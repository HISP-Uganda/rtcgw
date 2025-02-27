package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"rtcgw/config"
	"rtcgw/controllers"
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

	// Define template functions
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s) // Mark string as safe HTML
		},
	}
	// Load templates with custom functions
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob(
		config.RTCGwConf.Server.TemplatesDirectory + "/*"))
	router.SetHTMLTemplate(tmpl)

	// Serve Static Files
	router.Static("/static", config.RTCGwConf.Server.StaticDirectory)

	// Home Route
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "API Documentation",
		})
	})

	// Documentation Routes
	router.GET("/docs/:page", func(c *gin.Context) {
		docName := c.Param("page")

		// Construct Markdown file path
		mdFile := fmt.Sprintf("%s/%s.md", config.RTCGwConf.Server.DocsDirectory, docName)

		// Read Markdown file
		mdContent, err := os.ReadFile(mdFile)
		if err != nil {
			c.String(http.StatusNotFound, "Documentation not found")
			return
		}

		// Convert Markdown to HTML
		htmlContent := template.HTML(markdown.ToHTML(mdContent, nil, nil))

		// Render docs.html template
		c.HTML(http.StatusOK, "docs.html", gin.H{
			"title":   docName,
			"content": htmlContent,
		})
	})

	v2 := router.Group("/api", BasicAuth())
	{
		v2.GET("/test2", func(c *gin.Context) {
			c.String(200, "Authorized")
		})
		r := new(controllers.ResultsController)
		v2.POST("/results", r.Start)

		e := new(controllers.ClientsController)
		v2.POST("/clients", e.Start)

		userController := &controllers.UserController{}
		v2.GET("/users/:uid", userController.GetUserByUID)
		v2.PUT("/users/:uid", userController.UpdateUser)
		v2.POST("/users/getToken", userController.CreateUserToken)
		v2.POST("/users/refreshToken", userController.RefreshUserToken)
	}
	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		c.String(404, "Page Not Found!")
	})

	if err := router.Run(":" + fmt.Sprintf("%s", config.RTCGwConf.Server.Port)); err != nil {
		log.Fatalf("Could not start GIN server: %v", err)
	}
}
