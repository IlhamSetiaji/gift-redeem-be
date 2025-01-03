package utils

import (
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/request"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/response"
)

var ResponseChannel = make(chan map[string]interface{}, 100)

var Rchans = make(map[string](chan response.RabbitMQResponse))

type RabbitMsgPublisher struct {
	QueueName string                  `json:"queueName"`
	Message   request.RabbitMQRequest `json:"message"`
}

type RabbitMsgConsumer struct {
	QueueName string                    `json:"queueName"`
	Reply     response.RabbitMQResponse `json:"reply"`
}

// channel to publish rabbit messages
var Pchan = make(chan RabbitMsgPublisher, 10)
var Rchan = make(chan RabbitMsgConsumer, 10)
