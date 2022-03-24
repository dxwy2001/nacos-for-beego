package nacos

import (
	"fmt"
	"testing"
)

func TestIniParser_Parse(t *testing.T) {
	str := "gopub=4\ntest=21"
	parser := IniParser{}
	data, _ := parser.Parse(str)
	fmt.Println(data)
}
