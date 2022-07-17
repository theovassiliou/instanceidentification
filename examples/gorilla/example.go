package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"

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

type Status struct {
	Name   string
	Status string
}

// Writing simple X-Instance-Id header
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	stat := Status{"status", "running"}
	js, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add(iid.XINSTANCEID, iid.NewStdCiid(THISSERVICE).SetEpoch(startTime).String())
	w.Write(js)
}

// Writing complex X-Instance-Id header
func HealthHandler(w http.ResponseWriter, r *http.Request) {

	ciid := thisServiceCIID
	var callStack = ciid.Ciids()

	callStack.Push(iid.NewStdCiid("database/1.2%33s(storageService/0.2%77s)"))
	callStack.Push(iid.NewStdCiid("monitoring/1.1%22242s"))
	ciid.SetCiids(callStack).SetEpoch(startTime)
	w.Header().Set(iid.XINSTANCEID, ciid.String())
	log.Println("We called the following services:", ciid.(*iid.StdCiid).TreePrint())

	stat := Status{"health", "degraded"}
	js, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func Noop(w http.ResponseWriter, r *http.Request) {

	stat := Status{"ciid", "ok"}
	js, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
func main() {

	r := CiidRouter{
		mux.NewRouter(),
		thisServiceCIID,
	}
	r.Use(InstanceIdMiddleware(&r))
	r.HandleFunc("/status", StatusHandler)
	r.HandleFunc("/health", HealthHandler)
	r.HandleFunc("/noop", Noop)

	http.ListenAndServe(":8080", r)
}

type CiidRouter struct {
	*mux.Router
	Ciid iid.Ciid
}

func InstanceIdMiddleware(r *CiidRouter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// We only want to reply with the header if requested
			if w.Header().Get(iid.XINSTANCEID) == "" {
				w.Header().Add(iid.XINSTANCEID, r.Ciid.SetEpoch(startTime).String())
			}
			next.ServeHTTP(w, req)
		})
	}
}
