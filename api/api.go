package api

import (
	"fmt"
	"github.com/avvero/kid_agent_direct/utils"
)

type Message struct {
	ConversationId string `json:"conversationId"`
	Text           string `json:"text"`
}

type ApiClient struct {
	sendMessageEndpoint string
}

func NewApiClient(host string) *ApiClient {
	sendMessageEndpoint := fmt.Sprintf("%s/api/message", host)
	return &ApiClient{sendMessageEndpoint}
}

func (this ApiClient) SendMessage(conversationId string, text string) (error) {
	message := &Message{conversationId, text}
	_, err := utils.HttpPost(this.sendMessageEndpoint, message)
	return err
}
