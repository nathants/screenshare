package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/valyala/fasthttp"
)

var (
	lock   sync.RWMutex
	buffer bytes.Buffer
	auth   string
	millis int
	images = make(chan *image.RGBA, 1)
)

func router(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	if auth != string(args.Peek("auth")) {
		_, err := ctx.WriteString("bad ?auth=")
		if err != nil {
			panic(err)
		}
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
	if string(ctx.Method()) != "GET" {
		_, err := ctx.WriteString("GET only")
		if err != nil {
			panic(err)
		}
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
	switch string(ctx.Path()) {
	case "/":
		indexHandler(ctx)
	case "/axios.min.js":
		jsHandler(ctx)
	case "/img.jpg":
		imageHandler(ctx)
	}
}

func jsHandler(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString(axios)
	if err != nil {
		panic(err)
	}
	ctx.SetContentType("text/javascript; charset=utf8")
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	val := index
	val = strings.Replace(val, "AUTH", auth, -1)
	val = strings.Replace(val, "MILLIS_PER_FRAME", fmt.Sprint(millis), -1)
	_, err := ctx.WriteString(val)
	if err != nil {
		panic(err)
	}
	ctx.SetContentType("text/html; charset=utf8")
}

func imageHandler(ctx *fasthttp.RequestCtx) {
	lock.RLock()
	_, err := ctx.Write(buffer.Bytes())
	lock.RUnlock()
	if err != nil {
		panic(err)
	}
	ctx.SetContentType("application/octet-stream")
}

func encoder() {
	last_print := time.Now()
	var buf bytes.Buffer
	var count int64
	for {

		// encode jpg
		buf.Reset()
		img := <-images
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
		if err != nil {
			panic(err)
		}

		// update image bytes
		lock.Lock()
		buffer.Reset()
		buffer.Write(buf.Bytes())
		lock.Unlock()

		// print stats
		count++
		if time.Since(last_print) > time.Second {
			fmt.Println("millis per frame:", time.Since(last_print).Milliseconds()/count)
			last_print = time.Now()
			count = 0
		}

	}
}

func capturer(display int) {
	bounds := screenshot.GetDisplayBounds(display)
	for {
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		images <- img
	}
}

func flags() (int, int) {
	display := flag.Int("d", 1, "display number")
	port := flag.Int("p", 8080, "port")
	_millis := flag.Int("m", 30, "millis per frame")
	_auth := flag.String("a", "", "auth: http://localhost:8080?auth=AUTH")
	flag.Parse()
	auth = *_auth
	millis = *_millis
	return *display, *port
}

func serve(port int) {
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func main() {
	display, port := flags()
	go capturer(display)
	go encoder()
	serve(port)
}
