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

// ParamStruct stroe user's input parameter
type ParamStruct struct {
	Target        string
	Thread        int
	OutputNormal  string
	Port          string
	InputFileName string
	ExcludeFile   string
	Exclude       string
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

func getHostFromString(hostStr string) []string {
	var hostList []string
	// Read from terminal parameter
	tmpList := strings.Split(hostStr, ",")
	for _, i := range tmpList {
		// handling ip address, for example 192.168.1.0/24
		tmp, _ := gip.GetIPSubnet(i, 32)
		hostList = append(hostList, tmp...)
	}
	return hostList
}

func getHostFromFile(filePath string) ([]string, error) {
	var hostList []string
	isExits, err := PathFileExists(filePath)
	if isExits {
		content, _ := ioutil.ReadFile(filePath)
		result := strings.Replace(string(content), "\r", "", -1)
		tmpList := strings.Split(result, "\n")
		for _, i := range tmpList {
			// handling ip address, for example 192.168.1.0/24
			tmp, _ := gip.GetIPSubnet(i, 32)
			hostList = append(hostList, tmp...)
		}
		return hostList, nil
	}
	return nil, err
}

// GetTargetList function handles Target and InputFilename parameters of ParamStruct.
func GetTargetList(param ParamStruct) ([]string, error) {
	var targetList []string

	// Read from terminal parameter
	if param.Target != "" {
		tmp := getHostFromString(param.Target)
		targetList = append(targetList, tmp...)
	}

	// Read from file
	if param.InputFileName != "" {
		tmp, err := getHostFromFile(param.InputFileName)
		if err != nil {
			return nil, err
		}
		targetList = append(targetList, tmp...)
	}

	if len(targetList) == 0 {
		err := errors.New("WARNING: No targets were specified, so 0 hosts scanned")
		return nil, err
	}
	targetList = filterExcludeHost(param, targetList)
	return targetList, nil
}

func filterExcludeHost(param ParamStruct, targetList []string) []string {
	excludeMap := make(map[string]bool)
	var newTargetList []string
	if param.Exclude != "" {
		tmp := getHostFromString(param.Exclude)
		for i := 0; i < len(tmp); i++ {
			excludeMap[tmp[i]] = true
		}
	}

	if param.ExcludeFile != "" {
		tmp, _ := getHostFromFile(param.ExcludeFile)
		for i := 0; i < len(tmp); i++ {
			excludeMap[tmp[i]] = true
		}
	}

	for i := 0; i < len(targetList); i++ {
		if _, ok := excludeMap[targetList[i]]; !ok {
			newTargetList = append(newTargetList, targetList[i])
		}
	}
	return newTargetList
}

// GetPortList function handles Port parameter of ParamStruct.
func GetPortList(param ParamStruct) ([]int, error) {
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

	if param.Port == "" {
		return defaultPortList, nil
	}

	baseList := strings.Split(param.Port, ",")
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

// ParamInit handles users' input.
func ParamInit() (ParamStruct, error) {
	var param ParamStruct

	flag.StringVar(&param.Target, "t", "", "-t <host1[,host2][,host3],...>: Scan specified targets")
	flag.StringVar(&param.Port, "p", "", "-p <port ranges>: Only scan specified ports")
	flag.StringVar(&param.Exclude, "exclude", "", "--exclude <host1[,host2][,host3],...>: Exclude hosts/networks")
	flag.StringVar(&param.ExcludeFile, "excludefile", "", "--excludefile <exclude_file>: Exclude list from file")
	flag.StringVar(&param.InputFileName, "iL", "", "-iL <input_file_name>: Input from list of hosts/networks")
	flag.StringVar(&param.OutputNormal, "oN", "", "-oN <file>: Output scan in normal")
	flag.IntVar(&param.Thread, "thread", 64, "--thread <size>: Parallel target scan group size")
	flag.Usage = usage
	flag.Parse()

	if param.Target == "" && param.InputFileName == "" {
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
