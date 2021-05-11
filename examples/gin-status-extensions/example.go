package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	ix "github.com/theovassiliou/instanceidentification/examples/gin-status-extensions/instanceidextended"

	iid "github.com/theovassiliou/instanceidentification"
)

var thisServiceCIID iid.Ciid

const VERSION = "0.1-src"

//set this via ldflags (see https://stackoverflow.com/q/11354518)
// version is the current version number as tagged via git tag 1.0.0 -m 'A message'
var (
	version   = VERSION
	commit    = ""
	branch    = ""
	cmdName   = "ginstatus-vbc"
	startTime time.Time
)

func init() {
	thisServiceCIID = ix.NewExtCiid(cmdName, version, branch, commit)
	startTime = time.Now()
}

func main() {

	r := gin.Default()

	r.Use(GenerateInstanceId())

	// -- Example returning only default MIID as CIID
	r.GET("/status", func(c *gin.Context) {
		fmt.Println(c.Request.Header)
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
		log.Println("We called the following services:", ciid.(*iid.StdCiid).TreePrint())

		c.JSON(200, gin.H{
			"health": "degraded",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func GenerateInstanceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &CiidResponseWriter{c.Writer, ix.NewExtCiid(cmdName, version, branch, commit)}
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
		fmt.Println(w.Ciid)
		w.Ciid.SetEpoch(startTime)
		fmt.Println(w.Ciid)
		w.Header().Add(iid.XINSTANCEID, w.Ciid.String())
	}

	w.ResponseWriter.WriteHeader(code)
}
