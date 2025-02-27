package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rtcgw/models"
	"rtcgw/tasks"
)

type ResultsController struct{}

func (r *ResultsController) Start(c *gin.Context) {
	var result models.LabXpertResult
	if err := c.ShouldBindJSON(&result); err != nil {
		RespondWithError(http.StatusBadRequest, err.Error(), c)
		return
	}
	c.JSON(200, gin.H{
		"message": "results queued for saving to DHIS2",
	})
	client := c.MustGet("asynqClient").(*asynq.Client)
	task, err := tasks.NewResultsTask(result)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue results task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return
}

// RespondWithError returns an error if request has an error
func RespondWithError(i int, s string, c *gin.Context) {
	c.JSON(i, gin.H{"error": s})
	c.Abort()
}
