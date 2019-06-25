.PHONY: nmc_monitor_client

VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null || true)

nmc_monitor_client:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nmc-monitor-client-go/cmd.Version=${VERSION}\" \
        -X \"github.com/nwpc-oper/nmc-monitor-client-go/cmd.BuildTime=${BUILD_TIME}\" \
        -X \"github.com/nwpc-oper/nmc-monitor-client-go/cmd.GitCommit=${GIT_COMMIT}\" " \
		-o bin/nmc_monitor_client \
		main.go