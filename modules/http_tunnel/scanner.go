package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

type HTTPTunnelConfig struct {
	Host string
	Port uint16
}

var channel chan *HTTPTunnelConfig

func CheckHTTPTunnel(host string, port uint16) (bool, error) {
	request := gorequest.New().Timeout(2 * time.Second)
	resp, body, errs := request.Proxy(fmt.Sprintf("http://%s:%d/", host, port)).Get("http://ifconfig.me/").End()
	fmt.Println(resp)
	fmt.Println(body)
	fmt.Println(errs)
	if body == host {
		return true, nil
	} else {
		return false, fmt.Errorf(body)
	}
}

func LoadHTTPTunnelConfigs(filepath string) {
	f, _ := os.OpenFile(filepath, os.O_RDONLY, 0644)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		host := strings.TrimSpace(scanner.Text())
		fmt.Println(host)
		channel <- &HTTPTunnelConfig{
			Host: host,
			Port: 8888,
		}
	}
}

func main() {
	// create channel
	channel = make(chan *HTTPTunnelConfig, 32)
	// load configs
	go LoadHTTPTunnelConfigs("input.txt")

	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// check proxies
	for httpTunnelConfig := range channel {
		if httpTunnelConfig == nil {
			break
		}
		status, err := CheckHTTPTunnel(httpTunnelConfig.Host, httpTunnelConfig.Port)
		f.WriteString(fmt.Sprintf("%s:%d,%t,%s\n", httpTunnelConfig.Host, httpTunnelConfig.Port, status, err))
		f.Sync()
	}
}
