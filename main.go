package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"os"
	"rtcgw/config"
	"rtcgw/models"
	"sync"
	"time"
)

func init() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = time.RFC3339
	formatter.FullTimestamp = true

	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}

var splash = `
┏━┓╺┳╸┏━╸┏━╸╻ ╻
┣┳┛ ┃ ┃  ┃╺┓┃╻┃
╹┗╸ ╹ ┗━╸┗━┛┗┻┛
`

var client *asynq.Client

func main() {
	fmt.Printf(splash)
	var wg sync.WaitGroup
	client = asynq.NewClient(asynq.RedisClientOpt{Addr: config.RTCGwConf.Server.RedisAddress})
	defer func(client *asynq.Client) {
		_ = client.Close()
	}(client)

	wg.Add(1)
	go startAPIServer(&wg)

	wg.Wait()
}

func startAPIServer(wg *sync.WaitGroup) {
	defer wg.Done()
	router := gin.Default()
	v2 := router.Group("/api", models.BasicAuth())
	{
		v2.GET("/test2", func(c *gin.Context) {
			c.String(200, "Authorized")
		})

	}
	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		c.String(404, "Page Not Found!")
	})

	if err := router.Run(":" + fmt.Sprintf("%s", config.RTCGwConf.Server.Port)); err != nil {
		log.Fatalf("Could not start GIN server: %v", err)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
