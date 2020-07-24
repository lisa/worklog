SHELL := /usr/bin/env bash

GO_SOURCES := $(find $(CURDIR) -type f -name "*.go" -print)
GOPATH := $(shell go env GOPATH)

all: worklog oncall

.PHONY: worklog
worklog: $(GO_SOURCES)
	go install pkg/worklog/worklog.go

.PHONY: oncall
oncall: $(GO_SOURCES)
	go install pkg/oncall/oncall.go