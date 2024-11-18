// Copyright 2024 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.
//
// Code in this file inspired by or generated with the help of ChatGPT, OpenAI
// and Google Gemini.

package nagios

import (
	"bytes"
	"compress/gzip"
	"encoding/ascii85"
	"fmt"
	"io"
	"regexp"
)

// getEncodedPayloadDelimiterLeft retrieves the custom left delimiter used
// when enclosing an encoded payload if set, otherwise returns the default
// value.
func (p Plugin) getEncodedPayloadDelimiterLeft() string {
	switch {
	case p.encodedPayloadDelimiterLeft != nil:
		return *p.encodedPayloadDelimiterLeft
	default:
		return defaultPayloadDelimiterLeft
	}
}

// getEncodedPayloadDelimiterRight retrieves the custom right delimiter used
// when enclosing an encoded payload if set, otherwise returns the default
// value.
func (p Plugin) getEncodedPayloadDelimiterRight() string {
	switch {
	case p.encodedPayloadDelimiterRight != nil:
		return *p.encodedPayloadDelimiterRight
	default:
		return defaultPayloadDelimiterRight
	}
}

// encodeASCII85 encodes the given input as Ascii85. If no input is provided,
// an empty string is returned. No compression is performed on given input.
//
// If specified, the given left and right delimiters are used to enclose the
// encoded payload. If not specified, no delimiters are used.
func encodeASCII85(data []byte, leftDelimiter string, rightDelimiter string) string {
	if len(data) == 0 {
		return ""
	}

	encoded := make([]byte, ascii85.MaxEncodedLen(len(data)))

	// Encode and trim the encoded slice to the exact number of encoded bytes.
	n := ascii85.Encode(encoded, data)
	encoded = encoded[:n]

	// Add optional delimiters.
	return leftDelimiter + string(encoded) + rightDelimiter
}

// unescapeASCII85 unescapes an Ascii85 input payload by removing escape
// patterns added to the payload as it passes through a monitoring system
// (e.g., for inclusion in a JSON API response).
func unescapeASCII85(encodedInput []byte) ([]byte, error) {
	if len(encodedInput) == 0 {
		return nil, fmt.Errorf(
			"failed to unescape empty payload: %w",
			ErrMissingValue,
		)
	}

	// Based on initial testing this is sufficient to unescape an Ascii85
	// payload that passes through the Nagios XI API.
	encodedInput = bytes.ReplaceAll(encodedInput, []byte(`\\`), []byte(`\`))

	return encodedInput, nil
}

// EncodePayload compresses and encodes the given input for inclusion in
// plugin output. If no input is provided, an empty string is returned. If an
// error is encountered during compression the given input is encoded directly
// without compression.
//
// If specified, the given left and right delimiters are used to enclose the
// encoded payload. If not specified, no delimiters are used.
func EncodePayload(data []byte, leftDelimiter string, rightDelimiter string) string {
	if len(data) == 0 {
		return ""
	}

	var encoded string

	compressedData, compressErr := compressPayloadContent(data)
	switch {
	case compressErr != nil:
		// Fallback to skipping compression if an error occurs, use original
		// payload buffer contents as-is.
		encoded = encodeASCII85(data, leftDelimiter, rightDelimiter)

	default:
		encoded = encodeASCII85(compressedData, leftDelimiter, rightDelimiter)
	}

	return encoded
}

// decompressPayloadContent returns given input in decompressed form or an
// error if one occurs during decompression.
func decompressPayloadContent(compressedContent []byte) ([]byte, error) {
	if len(compressedContent) == 0 {
		return nil, fmt.Errorf(
			"failed to decompress payload from empty input: %w",
			ErrMissingValue,
		)
	}

	// FIXME: Should we silently return the original content if it is not
	// compressed?
	if !isGzipCompressed(compressedContent) {
		return nil, fmt.Errorf(
			"failed to decompress payload: %w",
			ErrCompressedInputInvalid,
		)
	}

	dataReader := bytes.NewReader(compressedContent)

	gzipReader, err := gzip.NewReader(dataReader)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer func() {
		_ = gzipReader.Close()
	}()

	var decompressedData bytes.Buffer

	for {
		written, err := io.CopyN(&decompressedData, gzipReader, 1024)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf(
				"error reading from gzip reader: %w (%d bytes read, %d bytes written)",
				err, decompressedData.Len(), written,
			)
		}
	}

	return decompressedData.Bytes(), nil
}

// decodeASCII85 returns given Ascii85 encoded input in decoded form or an
// error if one occurs during decoding. No decompression is performed on given
// input.
//
// The caller is expected to decompress and remove any delimiters from the
// input before calling this function.
//
// This function is also not intended for extraction of an Ascii85 encoded
// payload from surrounding text.
func decodeASCII85(encodedInput []byte) ([]byte, error) {
	if len(encodedInput) == 0 {
		return nil, fmt.Errorf(
			"failed to decode empty payload: %w",
			ErrMissingValue,
		)
	}

	unescapedInput, unescapeErr := unescapeASCII85(encodedInput)
	if unescapeErr != nil {
		return nil, unescapeErr
	}

	decoded := make([]byte, len(unescapedInput))
	n, _, decodeErr := ascii85.Decode(decoded, unescapedInput, true)
	if decodeErr != nil {
		return nil, decodeErr
	}

	// Trim the decoded slice to the exact number of reportd decoded bytes to
	// prevent any extraneous null bytes (or other content) from being
	// unintentionally included.
	decodedBytes := decoded[:n]

	return decodedBytes, nil
}

// DecodePayload decodes and (if applicable) decompresses given encoded input
// or an error if one occurs during decoding. If provided, the left and right
// delimiters are trimmed from the given input before decoding is performed.
//
// This function is not intended to extract an encoded payload from
// surrounding text.
func DecodePayload(encodedInput []byte, leftDelimiter string, rightDelimiter string) ([]byte, error) {
	if len(encodedInput) == 0 {
		return nil, fmt.Errorf(
			"failed to decode empty payload: %w",
			ErrMissingValue,
		)
	}

	if leftDelimiter != "" {
		encodedInput = bytes.TrimPrefix(encodedInput, []byte(leftDelimiter))
	}

	if rightDelimiter != "" {
		encodedInput = bytes.TrimSuffix(encodedInput, []byte(rightDelimiter))
	}

	// fmt.Println("encodedInput after trimming:", string(encodedInput))

	decodedPayload, err := decodeASCII85(encodedInput)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode %d bytes input payload: %w",
			len(encodedInput),
			err,
		)
	}

	// An earlier payload compression attempt may have failed, causing the
	// encoding logic to fallback to using the original unencoded payload
	// buffer content.
	if isGzipCompressed(decodedPayload) {
		decodedPayload, err = decompressPayloadContent(decodedPayload)
		if err != nil {
			return nil, err
		}
	}

	return decodedPayload, nil
}

// ExtractEncodedPayload extracts an encoded payload from given text input
// using specified delimiters.
//
// If not provided, a default regular expression for the encoding format is
// used to perform matching/extraction.
//
// If specified, delimiters are removed during the extraction process.
//
// NOTE: While technically optional, the use of delimiters for matching an
// encoded payload is *highly* recommended; reliability of payload matching is
// *greatly* reduced without using delimiters.
//
// The extracted payload is encoded and will need to be decoded and then
// decompressed before the original content is accessible.
func ExtractEncodedPayload(text string, customRegex string, leftDelimiter string, rightDelimiter string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf(
			"failed to extract encoded payload from empty input: %w",
			ErrMissingValue,
		)
	}

	defaultMatchPattern := leftDelimiter + defaultEncodingPatternRegex + rightDelimiter

	chosenRegex := defaultMatchPattern
	if customRegex != "" {
		chosenRegex = leftDelimiter + customRegex + rightDelimiter
	}

	// Assert that combined expression is valid.
	re, err := regexp.Compile(chosenRegex)
	if err != nil {
		return "", fmt.Errorf(
			"failed to use regex %q to match encoded payload "+
				"in given text: %w",
			chosenRegex,
			ErrEncodedPayloadRegexInvalid,
		)
	}

	matches := re.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", fmt.Errorf("no encoded payload data found: %w", ErrEncodedPayloadNotFound)
	}

	// Dynamically remove the delimiters based on input delimiter length.
	leftDelimiterLength := len(leftDelimiter)
	rightDelimiterLength := len(rightDelimiter)

	return matches[0][leftDelimiterLength : len(matches[0])-rightDelimiterLength], nil
}

// ExtractAndDecodePayload extracts, decodes and decompresses an encoded
// payload from given input text.
//
// If not provided, a default regular expression for the encoding format is
// used to perform matching/extraction.
//
// If specified, delimiters are removed during the extraction process.
//
// NOTE: While technically optional, the use of delimiters for matching an
// encoded payload is *highly* recommended; without delimiters, reliability of
// payload matching is *greatly* reduced (LOTS of false positives).
//
// The final result is the original unencoded payload before compression and
// encoding was performed. Depending on the type of the original data, the
// retrieved payload may require additional processing (e.g., JSON vs
// plaintext).
func ExtractAndDecodePayload(text string, customRegex string, leftDelimiter string, rightDelimiter string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf(
			"failed to extract and decode payload from empty input: %w",
			ErrMissingValue,
		)
	}

	encodedPayload, err := ExtractEncodedPayload(text, customRegex, leftDelimiter, rightDelimiter)
	if err != nil {
		return "", err
	}

	decodedPayload, err := decodeASCII85([]byte(encodedPayload))
	if err != nil {
		return "", err
	}

	// An earlier payload compression attempt may have failed, causing the
	// encoding logic to fallback to using the original unencoded payload
	// buffer content. Due to this, we opt to skip decompressing what may
	// already be an uncompressed payload.
	if isGzipCompressed(decodedPayload) {
		decodedPayload, err = decompressPayloadContent(decodedPayload)
		if err != nil {
			return "", err
		}
	}

	return string(decodedPayload), nil
}

// compressPayloadBufferOrFallback returns the compressed payload buffer
// contents or the uncompressed/original payload buffer contents if an error
// occurs during compression.
func (p Plugin) compressPayloadBufferOrFallback() []byte {
	compressedData, compressErr := compressPayloadContent(p.encodedPayloadBuffer.Bytes())
	switch {
	case compressErr != nil:
		// Skip compression if an error occurs, use original payload buffer
		// contents as-is.
		p.logAction("failed to compress unencoded payload content, skipping compression")

		return p.encodedPayloadBuffer.Bytes()

	default:
		p.logAction("successfully compressed unencoded payload content")
		p.logPluginOutputSize(fmt.Sprintf("%d bytes plugin unencoded payload content after compression", len(compressedData)))

		return compressedData
	}
}

// compressPayloadContent compresses given input data or returns an error if
// one occurs.
func compressPayloadContent(uncompressedContent []byte) ([]byte, error) {
	var compressedBuffer bytes.Buffer

	gzipWriter, gzipWriterLevelErr := gzip.NewWriterLevel(&compressedBuffer, gzip.BestCompression)
	if gzipWriterLevelErr != nil {
		// Documentation notes that err is nil unless we specify an invalid
		// level; since we use a stdlib package constant that's highly
		// unlikely to produce a complaint, but we guard against it anyway.
		panic("invalid compression level specified")
	}
	defer func() {
		// Fallback close attempt in case later errors are encountered.
		if err := gzipWriter.Close(); err != nil {
			panic("failed to close gzip writer")
		}
	}()

	_, gzipWriteErr := io.Copy(gzipWriter, bytes.NewReader(uncompressedContent))
	if gzipWriteErr != nil {
		return nil, gzipWriteErr
	}

	// Explicitly close gzip writer to complete compression.
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return compressedBuffer.Bytes(), nil
}

// isGzipCompressed checks if the data is gzip-compressed by examining the
// header and asserting that the input data is at least two bytes long and
// starts with the gzip magic number (0x1F 0x8B).
func isGzipCompressed(data []byte) bool {
	return len(data) >= 2 && data[0] == 0x1F && data[1] == 0x8B
}

// isValidGzip checks entire content structure to assert that given input is
// valid gzip-compressed data. This is more resource intensive than just
// checking the header for the required gzip magic number but is also more
// reliable.
// func isValidGzip(data []byte) bool {
// 	_, err := gzip.NewReader(bytes.NewReader(data))
// 	return err == nil
// }

// func IsValidGzipHeader(data []byte) bool {
// 	// Check if data is at least 10 bytes to cover basic gzip header
// 	if len(data) < 10 {
// 		return false
// 	}
//
// 	// Check magic number
// 	if data[0] != 0x1F || data[1] != 0x8B {
// 		return false
// 	}
//
// 	// Check compression method
// 	if data[2] != 0x08 {
// 		return false
// 	}
//
// 	// Read flags and timestamp
// 	flags := data[3]
// 	timestamp := binary.LittleEndian.Uint32(data[4:8])
//
// 	fmt.Printf("Flags: %08b\n", flags)
// 	fmt.Printf("Timestamp: %d\n", timestamp)
//
// 	// Optionally, check for extra header fields based on flags
// 	const (
// 		FTEXT    = 1 << 0 // Text
// 		FHCRC    = 1 << 1 // Header CRC
// 		FEXTRA   = 1 << 2 // Extra fields
// 		FNAME    = 1 << 3 // Original file name
// 		FCOMMENT = 1 << 4 // File comment
// 	)
//
// 	if flags&FEXTRA != 0 {
// 		fmt.Println("Extra fields are present")
// 	}
//
// 	if flags&FNAME != 0 {
// 		fmt.Println("Original file name is present")
// 	}
//
// 	if flags&FCOMMENT != 0 {
// 		fmt.Println("File comment is present")
// 	}
//
// 	return true
// }
