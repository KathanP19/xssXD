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
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go workers(inputs, &wg)
	}
	wg.Wait()
}

func buildurl(s string) {
	ur, err := url.Parse(s)
	if err != nil {
		return
	}
	x := ur.Query()
	if len(x) == 0 {
		return
	}
	baseurl := ur.Scheme + "://" + ur.Host + ur.Path + "?"
	params := url.Values{}
	for i := range x {
		params.Add(i, "<'\">")
	}
	finalurl := baseurl + params.Encode()
	//fmt.Println("Testing", s)
	if checkxss(finalurl) {
		fmt.Println(s, "might be vulnerable to xss")
	}
}

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func checkxss(s string) bool {
	resp, err := http.Get(s)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return strings.Contains(string(body), "<'\">")
}

func workers(cha chan string, wg *sync.WaitGroup) {
	for i := range cha {
		buildurl(i)
	}
	wg.Done()
}
