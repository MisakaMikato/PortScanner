package parameter

import (
	"main/src/gip"
	"testing"
)

func TestGetTargetList(t *testing.T) {
	param := paramStruct{target: "192.168.0.0/24,127.0.0.1", inputFileName: "target.test"}

	excepted, _ := gip.GetIPSubnet("192.168.0.0/24", 32)
	excepted = append(excepted, "127.0.0.1")
	excepted = append(excepted, "192.168.118.234")
	tmp, _ := gip.GetIPSubnet("192.168.4.0/19", 32)
	excepted = append(excepted, tmp...)

	result, err := getTargetList(param)
	if err != nil {
		t.Fatalf("[ERROR] %s", err)
	}
	for i := 0; i < len(excepted); i++ {
		if result[i] != excepted[i] {
			t.Fatalf("getTargetList function failed, fail index %d, excepted %s, result %s", i, excepted[i], result[i])
		}
	}
}
