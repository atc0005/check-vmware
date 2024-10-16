// Copyright 2024 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.
//
// Code in this file inspired by or generated with the help of ChatGPT, OpenAI.

package nagios

import (
	"bytes"
	"encoding/ascii85"
	"fmt"
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

// EncodeASCII85Payload encodes the given input as Ascii85. If no input is
// provided, an empty string is returned.
//
// If specified, the given left and right delimiters are used to enclose the
// encoded payload. If not specified, no delimiters are used.
func EncodeASCII85Payload(data []byte, leftDelimiter string, rightDelimiter string) string {
	if len(data) == 0 {
		return ""
	}

	encoded := make([]byte, ascii85.MaxEncodedLen(len(data)))
	ascii85.Encode(encoded, data)

	// Add optional delimiters.
	return leftDelimiter + string(encoded) + rightDelimiter
}

// decodeASCII85 decodes given Ascii85 encoded input or an error if one occurs
// during decoding.
//
// The caller is expected to remove any delimiters from the input before
// calling this function.
//
// This function is also not intended for extraction of an Ascii encoded
// payload from surrounding text.
func decodeASCII85(encodedInput []byte) ([]byte, error) {
	if len(encodedInput) == 0 {
		return nil, fmt.Errorf(
			"failed to decode empty payload: %w",
			ErrMissingValue,
		)
	}

	decoded := make([]byte, len(encodedInput))
	n, _, err := ascii85.Decode(decoded, encodedInput, true)
	if err != nil {
		return nil, err
	}

	decodedBytes := decoded[:n]

	// Remove any trailing null (\x00) bytes that may have been added as
	// padding during the encoding process.
	//
	// https://blog.manugarri.com/note-to-self-fixing-encoding-in-golang-ascii85/
	decodedBytes = bytes.Trim(decodedBytes, "\x00")

	return decodedBytes, nil
}

// DecodeASCII85Payload decodes given Ascii85 encoded input or an error if one
// occurs during decoding. If provided, the left and right delimiters are
// trimmed from the given input before decoding is performed.
//
// This function is not intended to extract an Ascii encoded payload from
// surrounding text.
func DecodeASCII85Payload(encodedInput []byte, leftDelimiter string, rightDelimiter string) ([]byte, error) {
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

	decodedPayload, err := decodeASCII85(encodedInput)
	if err != nil {
		// return nil, err
		return nil, fmt.Errorf(
			"failed to decode %d bytes input payload: %w",
			len(encodedInput),
			err,
		)
	}

	return decodedPayload, nil
}

// ExtractEncodedASCII85Payload extracts an Ascii85 encoded payload from given
// text input using specified delimiters.
//
// If not provided, a default regular expression for the Ascii85 encoding
// format is used to perform matching/extraction.
//
// If specified, delimiters are removed during the extraction process.
//
// NOTE: While technically optional, the use of delimiters for matching an
// encoded payload is *highly* recommended; reliability of payload matching is
// *greatly* reduced without using delimiters.
//
// The extracted payload is Ascii85 encoded and will need to be decoded before
// the original content is accessible.
func ExtractEncodedASCII85Payload(text string, customRegex string, leftDelimiter string, rightDelimiter string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf(
			"failed to extract encoded payload from empty input: %w",
			ErrMissingValue,
		)
	}

	// Regular expression to match Ascii85 block without delimiters.
	ascii85EncodingPattern := DefaultASCII85EncodingPatternRegex

	defaultMatchPattern := leftDelimiter + ascii85EncodingPattern + rightDelimiter

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
		return "", fmt.Errorf("no encoded Ascii85 data found: %w", ErrEncodedPayloadNotFound)
	}

	// Dynamically remove the delimiters based on input delimiter length.
	leftDelimiterLength := len(leftDelimiter)
	rightDelimiterLength := len(rightDelimiter)

	return matches[0][leftDelimiterLength : len(matches[0])-rightDelimiterLength], nil
}

// ExtractAndDecodeASCII85Payload extracts and decodes a Ascii85 encoded
// payload from given input text.
//
// If not provided, a default regular expression for the Ascii85 encoding
// format is used to perform matching/extraction.
//
// If specified, delimiters are removed during the extraction process.
//
// NOTE: While technically optional, the use of delimiters for matching an
// encoded payload is *highly* recommended; without delimiters, reliability of
// payload matching is *greatly* reduced (LOTS of false positives).
//
// The extracted content is the original unencoded payload before Ascii85
// encoding was performed. Depending on the type of the original data, the
// retrieved payload may require additional processing (e.g., JSON vs
// plaintext).
func ExtractAndDecodeASCII85Payload(text string, customRegex string, leftDelimiter string, rightDelimiter string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf(
			"failed to extract and decode payload from empty input: %w",
			ErrMissingValue,
		)
	}

	encodedPayload, err := ExtractEncodedASCII85Payload(text, customRegex, leftDelimiter, rightDelimiter)
	if err != nil {
		return "", err
	}

	decodedPayload, err := decodeASCII85([]byte(encodedPayload))
	if err != nil {
		return "", err
	}

	return string(decodedPayload), nil
}
