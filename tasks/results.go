package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"rtcgw/models"
)

const (
	TypeSendResults = "results:send"
)

func NewResultsTask(request models.LabXpertResult) (*asynq.Task, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendResults, payload), nil
}

func HandleResultsTask(cxt context.Context, task *asynq.Task) error {
	var result models.LabXpertResult
	if err := json.Unmarshal(task.Payload(), &result); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// Send the results to DHIS2
	log.Printf("Sending result to DHIS2: %v", result.PatientID)
	result.CheckDhis2Presence(clients.Dhis2Client)
	log.Printf("Done sending result to DHIS2: %v", result.PatientID)
	return nil
}
