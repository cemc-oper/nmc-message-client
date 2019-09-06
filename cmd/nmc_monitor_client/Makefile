.PHONY: all

all:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nmc-monitor-client-go/cmd.Version=${VERSION}\" \
        -X \"github.com/nwpc-oper/nmc-monitor-client-go/cmd.BuildTime=${BUILD_TIME}\" \
        -X \"github.com/nwpc-oper/nmc-monitor-client-go/cmd.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nmc_monitor_client \
		main.go