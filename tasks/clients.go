package tasks

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
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

	return asynq.NewTask(TypeCreateClient, payload, asynq.MaxRetry(3)), nil
}

func HandleClientTask(ctx context.Context, task *asynq.Task) error {
	var client models.ECHISRequest
	if err := json.Unmarshal(task.Payload(), &client); err != nil {
		log.Infof("failed to unmarshal payload: %v", err)
		return err
	}

	// Save the client to DHIS2, if not already in DHIS2
	syncLog, err := models.GetSyncLogByECHISID(client.ECHISID)
	if err != nil {
		log.Infof("Error getting sync log for patient: %s: %v", client.ECHISID, err)
		return err
	}
	if syncLog == nil {
		// No match found in localDB hence in DHIS2
		client.SaveClient(clients.Dhis2Client)

		log.Infof("Client saved to DHIS2: %s", client.ECHISID)
	} else {
		log.Infof("Client already exists in DHIS2: %s", client.ECHISID)
	}

	return nil
}
