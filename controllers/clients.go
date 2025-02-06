package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"rtcgw/models"
	"rtcgw/tasks"
)

type ClientsController struct{}

func (b *ClientsController) Start(c *gin.Context) {
	var clientRequest models.ECHISRequest
	if err := c.ShouldBindJSON(&clientRequest); err != nil {
		RespondWithError(400, "Invalid JSON payload", c)
		return
	}
	c.JSON(200, gin.H{
		"message": "client queued for saving to DHIS2",
	})
	client := c.MustGet("asynqClient").(*asynq.Client)
	task, err := tasks.NewClientTask(clientRequest)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
