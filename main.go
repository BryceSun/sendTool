package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type collection struct {
	Host      string `json:"host"`
	Name      string `json:"collectionName"`
	Pass      int    `json:"pass"`
	Fail      int    `json:"fail"`
	Requestes []rq   `json:"requestes"`
}
type rq struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Params     string `json:"params"`
	Output     string `json:"output"`
	StatusCode int    `json:"status_code"`
	Cost       int64  `json:"cost"`
}

func main() {
	logFile, err := os.OpenFile("./sendtool_error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open error log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	defer logFile.Close()

	host := flag.String("host", "http://localhost", "请求地址的域名")
	jsonFile := flag.String("file", "./ioc_requestes.json", "存放请求的json文件")
	flag.Parse()
	fmt.Println("host:", *host)
	fmt.Println("file", *jsonFile)
	bs, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	//fmt.Println(string(bs))
	c := new(collection)
	err = json.Unmarshal(bs, c)

	if err != nil {
		log.Println(err)
		return
	}
	if len(c.Host) != 0 {
		*host = c.Host
	} else {
		c.Host = *host
	}
	_, err = http.DefaultClient.Get(*host)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("\n\n")
	log.Println("---------start send requests for " + c.Name + "------------")
	for i, q := range c.Requestes {
		log.Println(*host + q.Path)
		req, err := http.NewRequest(q.Method, *host+q.Path, bytes.NewBufferString(q.Params))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		now := time.Now()
		response, err := http.DefaultClient.Do(req)
		c.Requestes[i].Cost = time.Since(now).Milliseconds()
		c.Requestes[i].StatusCode = response.StatusCode
		if response.StatusCode != http.StatusOK {
			c.Fail++
		} else {
			c.Pass++
		}
		r, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		c.Requestes[i].Output = string(r)

	}
	r, _ := json.Marshal(c)

	var out bytes.Buffer
	err = json.Indent(&out, r, "", " ")
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile("./"+c.Name+"_sendtool_result.json", out.Bytes(), 0644)
}
