package main

import (
	"fmt"
	"net/http"
	"os"
	//"strconv"
	"sync"
	"time"

	//	"github.com/IdeaEvolver/cutter-pkg/service"
	"github.com/go-chi/chi"
)

const html string = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
<div id="root"></div>
<script>
function push (n) {
	document.getElementById("root").innerHTML = n
}
</script>
`

var mu = sync.Mutex{}
var shutdownCalled = false
var shutdownStarted time.Time

var termCalled = false
var termStarted time.Time

func hi(res http.ResponseWriter, req *http.Request) {
	time.Sleep(2 * time.Second)
	mu.Lock()
	wasTermCalled := termCalled
	termStartTime := termStarted
	mu.Unlock()
	extra := ""
	if wasTermCalled {
		ds := time.Now().Sub(termStartTime)
		extra = fmt.Sprintf(" (grace period started %v seconds ago)", ds.Seconds())
	}

	x := fmt.Sprintf("VERSION: %s%s\n", os.Getenv("VERSION"), extra)
	res.Write([]byte(x))
}

func stream(res http.ResponseWriter, req *http.Request) {

	res.Write([]byte(html))
	f, _ := res.(http.Flusher)
	start := time.Now()

	for {
		select {
		case <-time.After(time.Second):
			mu.Lock()
			wasShutdownCalled := shutdownCalled
			shutdownStartTime := shutdownStarted
			mu.Unlock()
			dt := time.Now().Sub(start)
			extra := ""
			if wasShutdownCalled {
				ds := time.Now().Sub(shutdownStartTime)
				extra = fmt.Sprintf(" (shutdown started %v seconds ago)", ds.Seconds())
			}
			res.Write([]byte(fmt.Sprintf(`<script>push("VERSION: %s: connection alive for %d ms %s");</script>`, os.Getenv("VERSION"), dt.Milliseconds(), extra)))
			f.Flush()
		case <-req.Context().Done():
			return
		}
	}

}

func main() {

	//maxShutdown, _ := strconv.ParseInt(os.Getenv("MAX_SHUTDOWN_TIME"), 10, 64)
	//gracePeriod, _ := strconv.ParseInt(os.Getenv("SHUTDOWN_GRACE_TIME"), 10, 64)

	router := chi.NewRouter()
	router.Method("GET", "/hi", http.HandlerFunc(hi))
	router.Method("GET", "/", http.HandlerFunc(stream))

	//svr := service.GracefulServer(&service.Config{
	//	Addr:                ":1234",
	//	MaxShutdownTime:     time.Second * time.Duration(maxShutdown),
	//	ShutdownGracePeriod: time.Second * time.Duration(gracePeriod),
	//}, router)

	//svr.RegisterOnShutdown(func() {
	//	mu.Lock()
	//	shutdownCalled = true
	//	shutdownStarted = time.Now()
	//	mu.Unlock()
	//})

	//svr.RegisterOnReady(func() {
	//	fmt.Println("GOT TERM")
	//	mu.Lock()
	//	termCalled = true
	//	termStarted = time.Now()
	//	mu.Unlock()
	//})

	http.ListenAndServe(":1234", router)

	//svr.ListenAndServe()

}
