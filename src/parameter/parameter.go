// Package parameter implements handling user's input
package parameter

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"main/src/gip"
	"os"
	"strconv"
	"strings"
)

type paramStruct struct {
	target        string
	thread        int
	outputNormal  string
	port          string
	inputFileName string
	excludeFile   string
	exclude       string
}

// PathFileExists checks if path is existential
func PathFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

// func RemoveRepElement(slc []string) []string {
// 	vis := make(map[string]bool)
// 	for i := 0; i < len(slc); i++ {

// 	}
// }

func getTargetList(param paramStruct) ([]string, error) {
	var targetList []string

	// Read from terminal parameter
	if param.target != "" {
		tmpList := strings.Split(param.target, ",")
		for _, i := range tmpList {
			// handling ip address, for example 192.168.1.0/24
			tmp, _ := gip.GetIPSubnet(i, 32)
			targetList = append(targetList, tmp...)
		}
	}

	// Read from file
	if param.inputFileName != "" {
		isExits, err := PathFileExists(param.inputFileName)
		if isExits {
			content, _ := ioutil.ReadFile(param.inputFileName)
			result := strings.Replace(string(content), "\r", "", -1)
			tmpList := strings.Split(result, "\n")
			for _, i := range tmpList {
				// handling ip address, for example 192.168.1.0/24
				tmp, _ := gip.GetIPSubnet(i, 32)
				targetList = append(targetList, tmp...)
			}
		} else {
			return nil, err
		}
	}

	if len(targetList) == 0 {
		err := errors.New("WARNING: No targets were specified, so 0 hosts scanned")
		return nil, err
	}
	return targetList, nil
}

func getPortList(param paramStruct) ([]int, error) {
	var portList []int
	vis := make(map[int]bool)

	defaultPortList := []int{
		21, 22, 23, 25, 53, 53, 80, 81, 110, 111, 123, 123, 135, 137, 139, 161, 389, 443,
		445, 465, 500, 515, 520, 523, 548, 623, 636, 873, 902, 1080, 1099, 1433, 1521, 1604,
		1645, 1701, 1883, 1900, 2049, 2181, 2375, 2379, 2425, 3128, 3306, 3389, 4730, 5060,
		5222, 5351, 5353, 5432, 5555, 5601, 5672, 5683, 5900, 5938, 5984, 6000, 6379, 7001,
		7077, 8080, 8081, 8443, 8545, 8686, 9000, 9042, 9092, 9100, 9200, 9418, 9999, 11211,
		27017, 37777, 50000, 50070, 61616,
	}

	if param.port == "" {
		return defaultPortList, nil
	}

	baseList := strings.Split(param.port, ",")
	for i := 0; i < len(baseList); i++ {
		// if the format of port is like 1-20000
		seqIndex := strings.Index(baseList[i], "-")
		if seqIndex != -1 {
			start, err1 := strconv.Atoi(baseList[i][0:seqIndex])
			end, err2 := strconv.Atoi(baseList[i][seqIndex+1:])
			if err1 != nil {
				return nil, err1
			}
			if err2 != nil {
				return nil, err2
			}
			if start > end {
				errorMess := fmt.Sprintf(
					"Your port range %d-%d is backwards. Did you mean %d-%d?",
					start, end, end, start,
				)
				return nil, errors.New(errorMess)
			}
			for i := start; i <= end; i++ {
				if _, ok := vis[i]; !ok {
					vis[i] = true
					portList = append(portList, i)
				}
			}
		} else {
			port, err := strconv.Atoi(baseList[i])
			if err != nil {
				return nil, err
			}
			if _, ok := vis[port]; !ok {
				vis[port] = true
				portList = append(portList, port)
			}
		}
	}

	return portList, nil

}

func paramInit() (paramStruct, error) {
	var param paramStruct

	flag.StringVar(&param.target, "t", "", "-t <host1[,host2][,host3],...>: Scan specified targets")
	flag.StringVar(&param.port, "p", "", "-p <port ranges>: Only scan specified ports")
	flag.StringVar(&param.exclude, "exclude", "", "--exclude <host1[,host2][,host3],...>: Exclude hosts/networks")
	flag.StringVar(&param.excludeFile, "excludefile", "", "--excludefile <exclude_file>: Exclude list from file")
	flag.StringVar(&param.inputFileName, "iL", "", "-iL <input_file_name>: Input from list of hosts/networks")
	flag.StringVar(&param.outputNormal, "oN", "", "-oN <file>: Output scan in normal")
	flag.IntVar(&param.thread, "thread", 64, "--thread <size>: Parallel target scan group size")
	flag.Usage = usage
	flag.Parse()

	if param.target == "" && param.inputFileName == "" {
		err := errors.New("WARNING: No targets were specified, so 0 hosts scanned")
		return param, err
	}

	return param, nil
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-t target] [-p port_range] [Options]\n\nOptions:\n",
		os.Args[0],
	)
	flag.PrintDefaults()
}
