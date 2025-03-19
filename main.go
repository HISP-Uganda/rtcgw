package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"rtcgw/config"
	"rtcgw/controllers"
	"rtcgw/models"
	"rtcgw/models/stats"
	_ "rtcgw/models/stats"
	"rtcgw/utils"
	"strings"
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

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

func startAPIServer(wg *sync.WaitGroup) {
	defer wg.Done()
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})
		v.RegisterValidation("ugandaNIN", utils.UgandaNINValidation)
		v.RegisterValidation("dhis2UID", utils.Dhis2UIDValidation)
		v.RegisterValidation("yesNo", utils.YesNoValidation)
		v.RegisterValidation("maleFemale", utils.MaleFemaleValidation)
	}

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
	router.GET("/stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "stats.html", gin.H{"title": "Stats"})
	})
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			xAxisCategories := utils.LastXDays(7)
			log.Infof("ECHIS creations %v --- Days: %v: Updates. %v",
				models.SyncLogCreationLastXDays(7), xAxisCategories, models.SyncLogUpdateLastXDays(7))
			seriesData := []stats.SeriesData{
				{Name: "Created eCBSS clients", Data: models.SyncLogCreationLastXDays(7)},
				{Name: "Updated results in eCBSS", Data: models.SyncLogUpdateLastXDays(7)},
			}
			chartConfig := stats.ChartConfig{
				XAxisCategories: xAxisCategories,
				Series:          seriesData,
			}

			barData := randomData(5)
			pieData := randomPieData()

			data := map[string]interface{}{
				// "categories": []string{"Jan", "Feb", "Mar", "Apr", "May"},
				"categories":    utils.LastXDays(7),
				"barValues":     barData,
				"pieValues":     pieData,
				"timelineChart": chartConfig,
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Println("Write error:", err)
				break
			}

			time.Sleep(3 * time.Second)
		}
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

func randomData(n int) []int {
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = rand.Intn(300) + 50
	}
	return values
}

func randomPieData() []map[string]interface{} {
	return []map[string]interface{}{
		{"value": rand.Intn(40) + 10, "name": "Category A"},
		{"value": rand.Intn(40) + 10, "name": "Category B"},
		{"value": rand.Intn(40) + 10, "name": "Category C"},
		{"value": rand.Intn(40) + 10, "name": "Category D"},
	}
}
