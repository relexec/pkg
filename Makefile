# Top-level Makefile containing make targets that automate build, release,
# deployment and local testing tasks.

GIT_VERSION ?= $(shell git describe --tags --always --dirty || echo "unknown")
GIT_COMMIT ?= $(shell git rev-parse HEAD)

include .mk/build.mk
include .mk/go.mk
include .mk/test.mk
include .mk/help.mk
