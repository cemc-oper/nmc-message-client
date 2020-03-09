.PHONY: all

all:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nmc-message-client/cmd/nmc_message_consumer/app.Version=${VERSION}\" \
        -X \"github.com/nwpc-oper/nmc-message-client/cmd/nmc_message_consumer/app.BuildTime=${BUILD_TIME}\" \
        -X \"github.com/nwpc-oper/nmc-message-client/cmd/nmc_message_consumer/app.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nmc_monitor_consumer \
		main.go