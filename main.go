package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type collection struct {
	Host      string `json:"host"`
	Name      string `json:"collectionName"`
	Requestes []rq   `json:"requestes"`
}
type rq struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Params string `json:"params"`
}

func main() {
	host := flag.String("host", "http://localhost", "请求地址的域名")
	file := flag.String("file", "./ioc_requestes.json", "存放请求的json文件")

	flag.Parse()
	fmt.Println("host:", *host)
	fmt.Println("file", *file)
	bs, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(string(bs))
	c := new(collection)
	err = json.Unmarshal(bs, c)
	if err != nil {
		log.Println(err)
	}
	if len(c.Host) != 0 {
		*host = c.Host
	}
	for _, q := range c.Requestes {
		req, err := http.NewRequest("POST", strings.ReplaceAll(q.Path, "#host", *host), bytes.NewBufferString(q.Params))
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(req)
		log.Println(response)
	}
}
