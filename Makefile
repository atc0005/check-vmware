# Copyright 2021 Adam Chalkley
#
# https://github.com/atc0005/check-vmware
#
# Licensed under the MIT License. See LICENSE file in the project root for
# full license information.

# References:
#
# https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies
# https://github.com/mapnik/sphinx-docs/blob/master/Makefile
# https://stackoverflow.com/questions/23843106/how-to-set-child-process-environment-variable-in-makefile
# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
# https://unix.stackexchange.com/questions/124386/using-a-make-rule-to-call-another
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
# https://www.gnu.org/software/make/manual/html_node/Recipe-Syntax.html#Recipe-Syntax
# https://www.gnu.org/software/make/manual/html_node/Special-Variables.html#Special-Variables
# https://danishpraka.sh/2019/12/07/using-makefiles-for-go.html
# https://gist.github.com/subfuzion/0bd969d08fe0d8b5cc4b23c795854a13
# https://stackoverflow.com/questions/10858261/abort-makefile-if-variable-not-set
# https://stackoverflow.com/questions/38801796/makefile-set-if-variable-is-empty

SHELL = /bin/bash

# Space-separated list of cmd/BINARY_NAME directories to build
WHAT 					= check_vmware_tools \
							check_vmware_vcpus \
							check_vmware_vhw \
							check_vmware_hs2ds2vms \
							check_vmware_datastore_space \
							check_vmware_datastore_performance \
							check_vmware_snapshots_age \
							check_vmware_snapshots_count \
							check_vmware_snapshots_size \
							check_vmware_rps_memory \
							check_vmware_host_memory \
							check_vmware_host_cpu \
							check_vmware_vm_power_uptime \
							check_vmware_disk_consolidation \
							check_vmware_question \
							check_vmware_alarms \
							check_vmware_vm_backup_via_ca \

PROJECT_NAME			= check-vmware

# What package holds the "version" variable used in branding/version output?
# VERSION_VAR_PKG			= $(shell go list .)
VERSION_VAR_PKG			= $(shell go list .)/internal/config
# VERSION_VAR_PKG			= main

OUTPUTDIR 				= release_assets

ROOT_PATH				:= $(CURDIR)/$(OUTPUTDIR)

# https://gist.github.com/TheHippo/7e4d9ec4b7ed4c0d7a39839e6800cc16
VERSION 				= $(shell git describe --always --long --dirty)

# Used when generating download URLs when building assets for public release
RELEASE_TAG 			= $(shell git describe --exact-match --tags)

BASE_URL				= https://github.com/atc0005/$(PROJECT_NAME)/releases/download

# The default `go build` process embeds debugging information. Building
# without that debugging information reduces the binary size by around 28%.
#
# We also include additional flags in an effort to generate static binaries
# that do not have external dependencies. As of Go 1.15 this still appears to
# be a mixed bag, so YMMV.
#
# See https://github.com/golang/go/issues/26492 for more information.
#
# -s
#	Omit the symbol table and debug information.
#
# -w
#	Omit the DWARF symbol table.
#
# -tags 'osusergo,netgo'
#	Use pure Go implementation of user and group id/name resolution.
#	Use pure Go implementation of DNS resolver.
#
# -extldflags '-static'
#	Pass 'static' flag to external linker.
#
# -linkmode=external
#	https://golang.org/src/cmd/cgo/doc.go
#
#   NOTE: Using external linker requires installation of `gcc-multilib`
#   package when building 32-bit binaries on a Debian/Ubuntu system. It also
#   seems to result in an unstable build that crashes on startup. This *might*
#   be specific to the WSL environment used for builds, but since this is a
#   new issue and and I do not yet know much about this option, I am leaving
#   it out.
#
# CGO_ENABLED=0
#	https://golang.org/cmd/cgo/
#	explicitly disable use of cgo
#	removes potential need for linkage against local c library (e.g., glibc)
BUILDCMD				=	CGO_ENABLED=0 go build -mod=vendor -trimpath -a -ldflags "-s -w -X $(VERSION_VAR_PKG).version=$(VERSION)"
QUICK_BUILDCMD			=	go build -mod=vendor
GOCLEANCMD				=	go clean -mod=vendor ./...
GITCLEANCMD				= 	git clean -xfd
CHECKSUMCMD				=	sha256sum -b

.DEFAULT_GOAL := help

  ##########################################################################
  # Targets will not work properly if a file with the same name is ever
  # created in this directory. We explicitly declare our targets to be phony
  # by making them a prerequisite of the special target .PHONY
  ##########################################################################

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: lintinstall
## lintinstall: install common linting tools
# https://github.com/golang/go/issues/30515#issuecomment-582044819
lintinstall:
	@echo "Installing linting tools"

	@export PATH="${PATH}:$(go env GOPATH)/bin"

	@echo "Installing latest staticcheck version ..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck --version

	@echo Installing latest stable golangci-lint version per official installation script ...
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
	golangci-lint --version

	@echo "Finished updating linting tools"

.PHONY: linting
## linting: runs common linting checks
linting:
	@echo "Running linting tools ..."

	@echo "Running go vet ..."
	@go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

	@echo "Running golangci-lint ..."
	@golangci-lint --version
	@golangci-lint run

	@echo "Running staticcheck ..."
	@staticcheck --version
	@staticcheck $(shell go list -mod=vendor ./... | grep -v /vendor/)

	@echo "Finished running linting checks"

.PHONY: gotests
## gotests: runs go test recursively, verbosely
gotests:
	@echo "Running go tests ..."
	@go test -mod=vendor ./...
	@echo "Finished running go tests"

.PHONY: goclean
## goclean: removes local build artifacts, temporary files, etc
goclean:
	@echo "Removing object files and cached files ..."
	@$(GOCLEANCMD)
	@echo "Removing any existing release assets"
	@mkdir -p "$(OUTPUTDIR)"
	@rm -vf $(wildcard ${OUTPUTDIR}/*/*-linux-*)
	@rm -vf $(wildcard ${OUTPUTDIR}/*/*-windows-*)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-links.txt)
	@find ${OUTPUTDIR} -mindepth 1 -type d -empty -delete

.PHONY: clean
## clean: alias for goclean
clean: goclean

.PHONY: gitclean
## gitclean: WARNING - recursively cleans working tree by removing non-versioned files
gitclean:
	@echo "Removing non-versioned files ..."
	@$(GITCLEANCMD)

.PHONY: pristine
## pristine: run goclean and gitclean to remove local changes
pristine: goclean gitclean

.PHONY: quick
## quick: generates non-release binaries for current platform, arch
quick:
	@echo "Building non-release assets for current platform, arch ..."

	@for target in $(WHAT); do \
		echo "  Building $$target binary" && \
		$(QUICK_BUILDCMD) ${PWD}/cmd/$$target; \
	done

	@echo "Completed tasks for quick build"

.PHONY: all
# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
## all: generates assets for Linux distros and Windows
all: clean windows linux
	@echo "Completed all cross-platform builds ..."

.PHONY: windows-x86
## windows-x86: generates assets for Windows x86 systems
windows-x86:
	@echo "Building release assets for windows x86 ..."

	@for target in $(WHAT); do \
		mkdir -p $(ROOT_PATH)/$$target && \
		echo "  Building $$target 386 binary" && \
		env GOOS=windows GOARCH=386 $(BUILDCMD) -o $(ROOT_PATH)/$$target/$$target-$(VERSION)-windows-386.exe ${PWD}/cmd/$$target && \
		echo "  Generating $$target checksum file" && \
		cd $(ROOT_PATH)/$$target && \
		$(CHECKSUMCMD) $$target-$(VERSION)-windows-386.exe > $$target-$(VERSION)-windows-386.exe.sha256 && \
		cd $$OLDPWD; \
	done

	@echo "Completed build tasks for windows x86"

.PHONY: windows-x86-links
## windows-x86-links: generates download URLs for Windows x86 assets
windows-x86-links:
	@echo "Generating download links for windows x86 assets ..."

	@for target in $(WHAT); do \
		echo "  Generating $$target download links" && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-windows-386.exe" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-windows-386.exe.sha256" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt; \
	done

	@echo "Completed generating download links for windows x86 assets"

.PHONY: windows-x64
## windows-x64: generates assets for Windows x64 systems
windows-x64:
	@echo "Building release assets for windows x64 ..."

	@for target in $(WHAT); do \
		mkdir -p $(ROOT_PATH)/$$target && \
		echo "  Building $$target amd64 binary" && \
		env GOOS=windows GOARCH=amd64 $(BUILDCMD) -o $(ROOT_PATH)/$$target/$$target-$(VERSION)-windows-amd64.exe ${PWD}/cmd/$$target && \
		echo "  Generating $$target checksum file" && \
		cd $(ROOT_PATH)/$$target && \
		$(CHECKSUMCMD) $$target-$(VERSION)-windows-amd64.exe > $$target-$(VERSION)-windows-amd64.exe.sha256 && \
		cd $$OLDPWD; \
	done

	@echo "Completed build tasks for windows x64"

.PHONY: windows-x64-links
## windows-x64-links: generates download URLs for Windows x64 assets
windows-x64-links:
	@echo "Generating download links for windows x64 assets ..."

	@for target in $(WHAT); do \
		echo "  Generating $$target download links" && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-windows-amd64.exe" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-windows-amd64.exe.sha256" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt; \
	done

	@echo "Completed generating download links for windows x64 assets"

.PHONY: windows
## windows: generates assets for Windows x86 and x64 systems
windows: windows-x86 windows-x64
	@echo "Completed all build tasks for windows"

.PHONY: windows-links
## windows-links: generates download URLs for Windows x86 and x64 assets
windows-links: windows-x86-links windows-x64-links
	@echo "Completed generating download links for windows x86 and x64 assets"

.PHONY: linux-x86
## linux-x86: generates assets for Linux x86 distros
linux-x86:
	@echo "Building release assets for linux x86 ..."

	@for target in $(WHAT); do \
		mkdir -p $(ROOT_PATH)/$$target && \
		echo "  Building $$target 386 binary" && \
		env GOOS=linux GOARCH=386 $(BUILDCMD) -o $(ROOT_PATH)/$$target/$$target-$(VERSION)-linux-386 ${PWD}/cmd/$$target && \
		echo "  Generating $$target checksum file" && \
		cd $(ROOT_PATH)/$$target && \
		$(CHECKSUMCMD) $$target-$(VERSION)-linux-386 > $$target-$(VERSION)-linux-386.sha256; \
	done

	@echo "Completed build tasks for linux x86"

.PHONY: linux-x86-links
## linux-x86-links: generates download URLs for Linux x86 assets
linux-x86-links:
	@echo "Generating download links for linux x86 assets ..."

	@for target in $(WHAT); do \
		echo "  Generating $$target download links" && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-linux-386" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-linux-386.sha256" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt; \
	done

	@echo "Completed generating download links for linux x86 assets"

.PHONY: linux-x64
## linux-x64: generates assets for Linux x64 distros
linux-x64:
	@echo "Building release assets for linux x64 ..."

	@for target in $(WHAT); do \
		mkdir -p $(ROOT_PATH)/$$target && \
		echo "  Building $$target amd64 binary" && \
		env GOOS=linux GOARCH=amd64 $(BUILDCMD) -o $(ROOT_PATH)/$$target/$$target-$(VERSION)-linux-amd64 ${PWD}/cmd/$$target && \
		echo "  Generating $$target checksum file" && \
		cd $(ROOT_PATH)/$$target && \
		$(CHECKSUMCMD) $$target-$(VERSION)-linux-amd64 > $$target-$(VERSION)-linux-amd64.sha256; \
	done

	@echo "Completed build tasks for linux x64"

.PHONY: linux-x64-links
## linux-x64-links: generates download URLs for Linux x86 assets
linux-x64-links:
	@echo "Generating download links for linux x64 assets ..."

	@for target in $(WHAT); do \
		echo "  Generating $$target download links" && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-linux-amd64" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt && \
		echo "$(BASE_URL)/$(RELEASE_TAG)/$$target-$(VERSION)-linux-amd64.sha256" >> $(ROOT_PATH)/$(PROJECT_NAME)-$(VERSION)-links.txt; \
	done

	@echo "Completed generating download links for linux x64 assets"

.PHONY: linux
## linux: generates assets for Linux x86 and x64 distros
linux: linux-x86 linux-x64
	@echo "Completed all build tasks for linux"

.PHONY: linux-links
## linux-links: generates download URLs for Linux x86 and x64 assets
linux-links: linux-x86-links linux-x64-links
	@echo "Completed generating download links for linux x86 and x64 assets"

.PHONY: links
## links: generates download URLs for Windows and Linux x86 and x64 assets
links: windows-x86-links windows-x64-links linux-x86-links linux-x64-links
	@echo "Completed generating download links for all release assets"

.PHONY: release-build
## release-build: generates assets for public release
release-build: clean linux-x64 linux-x64-links

	@echo "Completed all tasks for release build"
