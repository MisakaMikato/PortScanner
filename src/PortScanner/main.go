package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/msterzhang/gpool"
)

// ScanResult represents the ip has opened a specific port.
type ScanResult struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
}

// Scanner launchs scan and output missions.
// if you want to print the result to stdout, just pass outPath empty("")
func Scanner(ipList []string, portList []string, threads int, outPath string) {
	output := make(chan ScanResult, threads)
	go ScanPort(ipList, portList, output, threads)
	Output(outPath, output)
}

// Output writes scan result to stdout if outPath is empty("")
// else writes to a specific file.
func Output(outPath string, outputChan chan ScanResult) {
	var f *os.File
	var err error
	if outPath == "" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(outPath, os.O_CREATE|os.O_RDWR, 0666)
		checkError(err)
		f.WriteString("[\n")
	}

	for result, ok := <-outputChan; ok || len(outputChan) != 0; result, ok = <-outputChan {
		content, err := json.MarshalIndent(result, "", "    ")
		checkError(err)
		f.WriteString(string(content) + ",\n")
	}

	if f != os.Stdout {
		f.WriteString("{}]")
	}
}

// ScanPort implements a simple port scanning.
func ScanPort(ipList []string, portList []string, output chan ScanResult, threads int) {
	defer close(output)

	pool := gpool.New(threads)
	for _, ip := range ipList {
		fmt.Printf("[INFO] Scanning %s\n", ip)
		for _, port := range portList {
			pool.Add(1)
			go ScanSinglePort(ip, port, output, pool)
		}
	}
	pool.Wait()
}

// ScanSinglePort implements a simple single port scanning.
func ScanSinglePort(ip string, port string, output chan ScanResult, pool *gpool.Pool) {
	defer func() {
		if pool != nil {
			pool.Done()
		}
	}()

	_, err := net.DialTimeout("tcp", ip+":"+port, time.Second*3)
	if err != nil {
		return
	}
	result := ScanResult{IP: ip, Port: port}
	output <- result
}

func main() {
	_, err := net.DialTimeout("tcp", "127.0.0.1:446", time.Second*3)
	checkError(err)
}
