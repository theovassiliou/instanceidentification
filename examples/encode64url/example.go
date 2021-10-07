package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/theovassiliou/base64url"
	iid "github.com/theovassiliou/instanceidentification"
)

// Default MIID of this service
const THISSERVICE = "encode64url/0.1%-1s"

var startTime time.Time
var thisServiceCIID iid.Ciid

func init() {
	thisServiceCIID = iid.NewStdCiid(THISSERVICE)
	startTime = time.Now()
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Use(GenerateInstanceId())

	r.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title":       "Home Page",
				"xinstanceid": thisServiceCIID.SetEpoch(startTime).String(),
			},
		)

	})

	r.POST("/encodestring", func(c *gin.Context) {
		fmt.Println()
		s, exists := c.GetPostForm("stringtoencode")
		if exists {
			c.JSON(200, gin.H{
				"encodedstring": base64url.Encode([]byte(s)),
				"x-instance-id": thisServiceCIID.SetEpoch(startTime).String(),
			})
		}
	})

	r.POST("/decodestring", func(c *gin.Context) {
		b, exists := c.GetPostForm("stringtodecode")
		if exists {
			s, _ := base64url.Decode(b)
			xiid := thisServiceCIID.SetEpoch(startTime)
			c.JSON(200, gin.H{
				"decodedstring": string(s),
				"x-instance-id": xiid.String(),
				"deciid":        xiid,
			})
		}
	})

	// -- Example returning a simple call-graph
	r.POST("/verify", func(c *gin.Context) {
		s, exists := c.GetPostForm("stringtoencode")
		if exists {
			xiid := thisServiceCIID.SetEpoch(startTime)
			fmt.Printf("%#v", xiid.Miid())
			encoded := base64url.Encode([]byte(s))
			encodedDecoded, _ := base64url.Decode(encoded)
			c.HTML(http.StatusOK, "verify.tmpl", gin.H{
				"stringtoencode":       s,
				"stringencoded":        encoded,
				"stringencodeddecoded": string(encodedDecoded),
				"equaling":             s == string(encodedDecoded),
				"xinstanceid":          xiid.String(),
				"sn":                   xiid.Miid().Sn(),
				"vn":                   xiid.Miid().Vn(),
				"va":                   xiid.Miid().Va(),
				"epoch":                xiid.Miid().T(),
			})
		}
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
