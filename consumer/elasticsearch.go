package consumer

import (
	"context"
	nmc_message_client "github.com/nwpc-oper/nmc-message-client"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type ElasticSearchTarget struct {
	Server string
}

type messageWithIndex struct {
	Index   string
	Id      string
	Message nmc_message_client.GribProduction
}

func pushMessages(client *elastic.Client, messages []messageWithIndex, ctx context.Context) {
	bulkRequest := client.Bulk()
	for _, indexMessage := range messages {
		request := elastic.NewBulkIndexRequest().
			Index(indexMessage.Index).
			Id(indexMessage.Id).
			Doc(indexMessage.Message)
		bulkRequest.Add(request)
	}
	_, err := bulkRequest.Do(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "elastic",
			"event":     "push",
		}).Errorf("%v", err)
		return
	}
}
