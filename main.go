package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

var (
	client = &fasthttp.Client{
		ReadTimeout:              time.Second,
		NoDefaultUserAgentHeader: true,
		MaxConnsPerHost:          10000,
	}
)

func validateHref(s string) (ok bool) {
	if s == "javascript:;" {
		return false
	}
	if s == "javascript:void(0);" {
		return false
	}

	if s == "#" {
		return false
	}
	if s == "javascript: void(0);" {
		return false
	}
	return true
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			if validateHref(a.Val) {
				href = a.Val
				ok = true
			}

		}
	}
	return
}

func validateAnchor(t html.Token, url string) (link string, err bool) {
	isAnchor := t.Data == "a"
	result := ""
	err = true
	if isAnchor {
		ok, ref := getHref(t)
		if ok {
			if strings.Index(ref, "/") == 0 {
				ref = ref[1:]

			}
			hasProtocol := strings.Index(ref, "http") == 0
			if hasProtocol {
				result = ref

			} else {
				if strings.Index(url, "http") == 0 {
					result = url + "/" + ref
				} else {
					result = "https://" + url + "/" + ref
				}

			}
			err = false
		}

	}
	return result, err
}

func doRequest(i interface{}) {
	url := i.(string)
	//fmt.Println("Send request to target:", url)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	client.Do(req, resp)
	statusCode := resp.StatusCode()
	if statusCode == 301 {
		fmt.Print("Redirect found")
	}
	fmt.Printf("[%d] %s \n", statusCode, url)
	body := resp.Body()
	z := html.NewTokenizer(bytes.NewReader(body))

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			newurl, err := validateAnchor(t, url)
			if !err {
				fmt.Println(newurl)
			}

		}
	}

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
	for i := 0; i < 55; i++ {
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
