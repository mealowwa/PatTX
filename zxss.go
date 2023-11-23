package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type paramCheck struct {
	url   string
	param string
}

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: time.Second,
		DualStack: true,
	}).DialContext,
}

var httpClient = &http.Client{
	Transport: transport,
}

func main() {

	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	sc := bufio.NewScanner(os.Stdin)

	initialChecks := make(chan paramCheck, 40)

	done := makePool(initialChecks, func(c paramCheck, output chan paramCheck) {
		reflected, err := checkReflected(c.url)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "error from checkReflected: %s\n", err)
			return
		}

		if len(reflected) == 0 {
			// TODO: wrap in verbose mode
			//fmt.Printf("no params were reflected in %s\n", c.url)
			return
		}

		fmt.Printf("URL: %s Param: %s \n", c.url , reflected)
	})

	for sc.Scan() {
		initialChecks <- paramCheck{url: sc.Text()}
	}

	close(initialChecks)
	<-done
}



func checkReflected(targetURL string) ([]string, error) {
	out := make([]string, 0)
	maxRetries := 50

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			return out, err
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36")

		resp, err := httpClient.Do(req)
		if err != nil {
			return out, err
		}

		if resp.StatusCode == http.StatusServiceUnavailable && attempt < maxRetries {

			// Sleep for a short duration before retrying
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if resp.Body == nil {
			return out, err
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return out, err
		}

		if strings.HasPrefix(resp.Status, "3") || (resp.Header.Get("Content-Type") != "" && !strings.Contains(resp.Header.Get("Content-Type"), "html")) {
			return out, nil
		}

		body := string(b)


		if strings.Contains(body, "Type the characters you see in this image") {
			continue
		}


		u, err := url.Parse(targetURL)
		if err != nil {
			return out, err
		}

		for key, vv := range u.Query() {
			for _, v := range vv {
				if !strings.Contains(body, v) {
					continue
				}

				out = append(out, key+"="+v)
			}
		}

		return out, nil
	}

	return out, nil
}


type workerFunc func(paramCheck, chan paramCheck)

func makePool(input chan paramCheck, fn workerFunc) chan paramCheck {
	var wg sync.WaitGroup

	output := make(chan paramCheck)
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			for c := range input {
				fn(c, output)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
