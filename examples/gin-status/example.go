package main

import (
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
	thisServiceCIID = iid.NewStdCiid(THISSERVICE)
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
		var callStack = ciid.Ciids()

		callStack.Push(iid.NewStdCiid("database/1.2%33s(storageService/0.2%77s)"))
		callStack.Push(iid.NewStdCiid("monitoring/1.1%22242s"))
		ciid.SetCiids(callStack).SetEpoch(startTime)
		c.Header(iid.XINSTANCEID, ciid.String())
		log.Println("We called the following services:", ciid.(iid.StdCiid).TreePrint())

		c.JSON(200, gin.H{
			"health": "degraded",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func GenerateInstanceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &CiidResponseWriter{c.Writer, iid.NewStdCiid(THISSERVICE)}
		c.Writer = writer
		c.Next()
	}
}

type CiidResponseWriter struct {
	gin.ResponseWriter
	Ciid iid.Ciid
}

func (w *CiidResponseWriter) WriteHeader(code int) {
	if w.Header().Get(iid.XINSTANCEID) == "" {
		w.Header().Add(iid.XINSTANCEID, w.Ciid.SetEpoch(startTime).String())
	}

	w.ResponseWriter.WriteHeader(code)
}
