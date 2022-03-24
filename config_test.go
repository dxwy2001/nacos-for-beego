package nacos

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func Test_test(t *testing.T) {
	urls := "http://172.16.10.125/nacos"
	config, _ := url.Parse(urls)
	fmt.Println(config.Port())

	str := "172.16.10.12"
	host := strings.SplitN(str, ":", 2)
	fmt.Println(len(host))
	fmt.Println(host[0])
}
