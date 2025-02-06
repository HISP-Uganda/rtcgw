package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"rtcgw/models"
)

const (
	TypeCreateClient = "client:create"
)

func NewClientTask(client models.ECHISRequest) (*asynq.Task, error) {
	payload, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeCreateClient, payload), nil
}
