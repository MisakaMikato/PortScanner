package parameter

import (
	"main/src/gip"
	"testing"
)

func TestGetTargetList(t *testing.T) {
	param := ParamStruct{Target: "192.168.0.0/24,127.0.0.1", InputFileName: "target.test"}

	excepted, _ := gip.GetIPSubnet("192.168.0.0/24", 32)
	excepted = append(excepted, "127.0.0.1")
	excepted = append(excepted, "192.168.118.234")
	tmp, _ := gip.GetIPSubnet("192.168.4.0/19", 32)
	excepted = append(excepted, tmp...)

	result, err := GetTargetList(param)
	if err != nil {
		t.Fatalf("[ERROR] %s", err)
	}
	for i := 0; i < len(excepted); i++ {
		if result[i] != excepted[i] {
			t.Fatalf("getTargetList function failed, fail index %d, excepted %s, result %s", i, excepted[i], result[i])
		}
	}
}

func TestGetPortList(t *testing.T) {
	var excepted1 []int
	defaultPortList := []int{
		21, 22, 23, 25, 53, 53, 80, 81, 110, 111, 123, 123, 135, 137, 139, 161, 389, 443,
		445, 465, 500, 515, 520, 523, 548, 623, 636, 873, 902, 1080, 1099, 1433, 1521, 1604,
		1645, 1701, 1883, 1900, 2049, 2181, 2375, 2379, 2425, 3128, 3306, 3389, 4730, 5060,
		5222, 5351, 5353, 5432, 5555, 5601, 5672, 5683, 5900, 5938, 5984, 6000, 6379, 7001,
		7077, 8080, 8081, 8443, 8545, 8686, 9000, 9042, 9092, 9100, 9200, 9418, 9999, 11211,
		27017, 37777, 50000, 50070, 61616,
	}
	for i := 1; i < 20001; i++ {
		excepted1 = append(excepted1, i)
	}
	excepted1 = append(excepted1, 30000, 30001)
	cases := []struct {
		param    ParamStruct
		excepted []int
	}{
		{ParamStruct{Port: "80"}, []int{80}},
		{ParamStruct{Port: "80,443,1443"}, []int{80, 443, 1443}},
		{ParamStruct{Port: "1-20000,30000,30001"}, excepted1},
		{ParamStruct{Port: "200-100,3,10"}, nil},
		{ParamStruct{Port: "1-20000,30000,30001,13000-16000,33"}, excepted1},
		{ParamStruct{Port: ""}, defaultPortList},
	}
	for _, c := range cases {
		result, _ := GetPortList(c.param)
		for i := 0; i < len(c.excepted); i++ {
			if result[i] != c.excepted[i] {
				t.Fatalf("getPortList function failed, fail index: %d, result: %d, excepted: %d", i, result[i], c.excepted[i])
			}
		}
	}
}

func TestFilterExcludeHost(t *testing.T) {
	param := ParamStruct{Target: "192.168.0.0/24,127.0.0.1", InputFileName: "target.test", Exclude: "127.0.0.1,192.168.0.0/25", ExcludeFile: "exclude.test"}
	targetList, _ := gip.GetIPSubnet("192.168.0.0/24", 32)
	expected, _ := gip.GetIPSubnet("192.168.0.192/26", 32)
	result := filterExcludeHost(param, targetList)

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Fatalf("filterExcludeHost function failed, fail index: %d, result: %s, expected: %s", i, result[i], expected[i])
		}
	}
}
