package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/valyala/fasthttp"
)

var (
	client = &fasthttp.Client{
		ReadTimeout:              time.Second,
		NoDefaultUserAgentHeader: true,
		MaxConnsPerHost:          10000,
	}
)

func doRequest(i interface{}) {
	// TODO
	// Add code for make request to target
	url := i.(string)
	fmt.Println("Send request to target:", url)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	client.Do(req, resp)
	bodyByte := resp.Body()
	fmt.Println(string(bodyByte))

}

func getAllLink(respondBody string) {
	return
}

func runScan(target string, concurency int) {
	defer ants.Release()
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(concurency, func(i interface{}) {
		doRequest(i)
		wg.Done()
	})
	defer p.Release()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		_ = p.Invoke(target)
	}
}

func main() {
	var target string
	flag.StringVar(&target, "t", "http://demo.testfire.net", "Target to scan")

	var concurency int
	flag.IntVar(&concurency, "c", 50, "Concurency number")

	flag.Parse()
	runScan(target, concurency)
}
