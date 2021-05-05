package main

import (
	"encoding/base64"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	iid "github.com/theovassiliou/instanceidentification"
)

// Default MIID of this service
const THISSERVICE = "ourService/1.1%-1s"

var startTime time.Time
var thisServiceCIID iid.Ciid

func init() {
	thisServiceCIID = iid.NewCiid(THISSERVICE)
	startTime = time.Now()
}

func main() {

	r := gin.Default()

	r.Use(GenerateInstanceId())

	// -- Example returning only default MIID as CIID
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
		})
	})

	// -- Example returning a simple call-graph
	r.GET("/health", func(c *gin.Context) {
		ciid := thisServiceCIID
		var callStack = ciid.Ciids

		callStack.Push(iid.NewCiid("database/1.2%33s(storageService/0.2%77s)"))
		callStack.Push(iid.NewCiid("monitoring/1.1%22242s"))
		ciid.SetStack(callStack).SetEpoch(startTime)
		c.Header("X-Instance-Id", ciid.String())
		log.Println("We called the following services:", ciid.PrintCiid())

		c.JSON(200, gin.H{
			"health": "degraded",
		})
	})

	// -- Example returning a simple call-graph with external services
	r.GET("/monitor", func(c *gin.Context) {
		ciid := thisServiceCIID
		var callStack = ciid.Ciids

		callStack.Push(iid.NewCiid("google/x/" + base64.StdEncoding.EncodeToString([]byte("http://www.google.com")) + "%-1s"))
		callStack.Push(iid.NewCiid("stackoverflow/x/" + base64.StdEncoding.EncodeToString([]byte("https://stackoverflow.com")) + "%-1s"))
		ciid.SetStack(callStack).SetEpoch(startTime)
		c.Header("X-Instance-Id", ciid.String())

		ciid.WithDecoding(func(s string) string {
			b1, _ := base64.StdEncoding.DecodeString(s)
			return string(b1)
		})

		log.Println("We called the following services:", ciid.PrintExtendedCiid())

		c.JSON(200, gin.H{
			"monitor": ciid.PrintExtendedCiid(),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func GenerateInstanceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &CiidResponseWriter{c.Writer, iid.NewCiid(THISSERVICE)}
		c.Writer = writer
		c.Next()
	}
}

type CiidResponseWriter struct {
	gin.ResponseWriter
	Ciid iid.Ciid
}

func (w *CiidResponseWriter) WriteHeader(code int) {
	if w.Header().Get("X-Instance-Id") == "" {
		w.Header().Add("X-Instance-Id", w.Ciid.SetEpoch(startTime).String())
	}

	w.ResponseWriter.WriteHeader(code)
}
