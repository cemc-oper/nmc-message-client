.PHONY: all

all:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nmc-message-client/cmd/nmc_monitor_client/app.Version=${VERSION}\" \
        -X \"github.com/nwpc-oper/nmc-message-client/cmd/nmc_monitor_client/app.BuildTime=${BUILD_TIME}\" \
        -X \"github.com/nwpc-oper/nmc-message-client/cmd/nmc_monitor_client/app.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nmc_monitor_client \
		main.go