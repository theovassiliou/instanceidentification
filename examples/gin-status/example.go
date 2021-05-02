package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	iid "github.com/theovassiliou/instanceidentification"
)

var startTime time.Time
var thisServiceCIID iid.Ciid

type Stack []iid.Ciid

const THISSERVICE = "ourService/1.1%-1s"

func init() {
	thisServiceCIID = iid.NewCiid(THISSERVICE)
	startTime = time.Now()
}

func main() {

	r := gin.Default()
	r.Use(DefaultInstanceId())

	r.GET("/health", func(c *gin.Context) {
		ciid := thisServiceCIID
		var callStack = ciid.Ciids

		callStack.Push(iid.NewCiid("database/1.2%33s(storageService/0.2%77s)"))
		callStack.Push(iid.NewCiid("monitoring/1.1%22242s"))
		ciid.SetStack(callStack).SetEpoch(startTime)
		c.Header("X-Instance-Id", ciid.String())
		log.Debugln("We called the following services:", iid.PrintCiid(ciid))
		c.JSON(200, gin.H{
			"health": "degraded",
		})
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
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

func DefaultInstanceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &CiidResponseWriter{c.Writer, iid.NewCiid(THISSERVICE)}
		c.Writer = writer
		c.Next()
	}
}
