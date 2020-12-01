package main

import (
	"bytes"
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

var lock sync.RWMutex
var buffer bytes.Buffer
var images = make(chan *image.RGBA)
var auth = ""
var millis = 30

func router(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != "GET" {
		panic("GET only")
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

func main() {

	go func() {
		bounds := screenshot.GetDisplayBounds(1)
		for {
			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				panic(err)
			}
			images <- img
		}
	}()

	go func() {
		start := time.Now()
		last_print := time.Now()
		var buf bytes.Buffer
		for {
			buf.Reset()
			img := <-images
			err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
			if err != nil {
				panic(err)
			}
			lock.Lock()
			buffer.Reset()
			buffer.Write(buf.Bytes())
			lock.Unlock()
			now := time.Now()
			if time.Since(last_print) > time.Second {
				fps := int(1 / now.Sub(start).Seconds())
				fmt.Println("fps:", fps)
				last_print = now
			}
			start = now
		}
	}()

	err := fasthttp.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}

}
