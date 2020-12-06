package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestScanPort(t *testing.T) {
	input := make(chan string, 2)
	output := make(chan ScanResult, 10)

	defer close(input)

	portList := []string{"445", "80", "443", "135", "23333"}

	ipList := []string{"127.0.0.1", "192.168.118.128"}

	ScanPort(ipList, portList, output, 5)

	fmt.Println("done")
	for result := range output {
		val, _ := json.Marshal(result)
		fmt.Println(string(val))
	}

}

func TestScanSinglePort(t *testing.T) {

	output := make(chan ScanResult, 10)

	ip := "127.0.0.1"
	port := "9100"
	ScanSinglePort(ip, port, output, nil)
	close(output)
	ch := <-output
	fmt.Printf("%s:%s\n", ch.IP, ch.Port)
}

func TestScanner(t *testing.T) {
	outPath := []string{"out.test", ""}
	portList := []string{"445", "80", "443", "135", "23333"}
	ipList := []string{"127.0.0.1", "192.168.118.128"}

	for i := 0; i < len(outPath); i++ {
		Scanner(ipList, portList, 5, outPath[i])
	}

}

func TestChannel(t *testing.T) {
	var ch = make(chan int, 10)
	ch <- 1
	ch <- 2
	l := len(ch)
	fmt.Printf("%d", l)
}
