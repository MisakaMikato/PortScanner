package main

import (
	"encoding/json"
	"fmt"
	"main/src/parameter"
	"net"
	"os"
	"time"

	"github.com/msterzhang/gpool"
)

// ScanResult represents the ip has opened a specific port.
type ScanResult struct {
	IP   string `json:"ip"`
	Port int    `json:"port,string"`
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
}

// Scanner launchs scan and output missions.
// if you want to print the result to stdout, just pass outPath empty("")
// TODO: use different 'outPath' format to implement variety output, eg: json://data.test, xml://data.test
func Scanner(ipList []string, portList []int, threads int, outPath string) {
	outputChan := make(chan ScanResult, threads)
	go ScanPort(ipList, portList, outputChan, threads)
	Output(outPath, outputChan)
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
func ScanPort(ipList []string, portList []int, output chan ScanResult, threads int) {
	defer close(output)

	pool := gpool.New(threads)
	for _, ip := range ipList {
		for _, port := range portList {
			pool.Add(1)
			go ScanSinglePort(ip, port, output, pool)
		}
	}
	pool.Wait()
}

// ScanSinglePort implements a simple single port scanning.
func ScanSinglePort(ip string, port int, output chan ScanResult, pool *gpool.Pool) {
	defer func() {
		if pool != nil {
			pool.Done()
		}
	}()
	tcpAddr := fmt.Sprintf("%s:%d", ip, port)
	_, err := net.DialTimeout("tcp", tcpAddr, time.Second*3)
	if err != nil {
		return
	}
	result := ScanResult{IP: ip, Port: port}
	output <- result
}

func main() {
	param, initErr := parameter.ParamInit()
	checkError(initErr)
	targetList, targetErr := parameter.GetTargetList(param)
	checkError(targetErr)
	portList, portErr := parameter.GetPortList(param)
	checkError(portErr)
	Scanner(targetList, portList, param.Thread, param.OutputNormal)
}
