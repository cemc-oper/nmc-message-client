.PHONY: all

all:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nwpc-message-client/cmd/nmic_message_client/app.Version=${VERSION}\" \
        -X \"github.com/nwpc-oper/nwpc-message-client/cmd/nmic_message_client/app.BuildTime=${BUILD_TIME}\" \
        -X \"github.com/nwpc-oper/nwpc-message-client/cmd/nmic_message_client/app.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nmic_message_client \
		main.go