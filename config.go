package nacos

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	nacosServerConfig []constant.ServerConfig
	nacosClientConfig constant.ClientConfig
	param             vo.ConfigParam
)

type NacosConfiger struct {
	dataId  string
	groupId string
	client  config_client.IConfigClient
	data    map[string]string
	sync.RWMutex
	config.BaseConfiger
}

func newNacosConfiger(client config_client.IConfigClient, dataId string, groupId string) *NacosConfiger {
	res := &NacosConfiger{
		client:  client,
		dataId:  dataId,
		groupId: groupId,
	}
	res.BaseConfiger = config.NewBaseConfiger(res.reader)
	return res
}

func (n *NacosConfiger) reader(ctx context.Context, key string) (string, error) {
	resp, _ := get(n.client, n.dataId, n.groupId)
	if resp != "" {
		logs.Error("err")
	}
	return "", nil
}

func (n NacosConfiger) Set(key, val string) error {
	n.data[key] = val
	fmt.Println(n.data)
	return errors.New("Unsupported operation")
}

func (n NacosConfiger) String(key string) (string, error) {
	return n.data[key], nil
}

func (n NacosConfiger) Bool(key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) Float(key string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DefaultString(key string, defaultVal string) string {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DefaultStrings(key string, defaultVal []string) []string {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DefaultInt(key string, defaultVal int) int {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DefaultInt64(key string, defaultVal int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DefaultBool(key string, defaultVal bool) bool {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DefaultFloat(key string, defaultVal float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) DIY(key string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) GetSection(section string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) Sub(key string) (config.Configer, error) {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) OnChange(key string, fn func(value string)) {
	//TODO implement me
	panic("implement me")
}

func (n NacosConfiger) SaveConfigFile(filename string) error {
	//TODO implement me
	panic("implement me")
}

type NacosConfigerProvider struct{}

func (provider *NacosConfigerProvider) Parse(key string) (config.Configer, error) {
	return provider.ParseData([]byte(key))
}

func (provider *NacosConfigerProvider) ParseData(data []byte) (config.Configer, error) {
	//cfg :=

	cfg := &NacosConfiger{
		data: make(map[string]string),
	}
	cfg.BaseConfiger = config.NewBaseConfiger(cfg.reader)

	cfg.Lock()
	defer cfg.Unlock()

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": nacosServerConfig,
		"clientConfig":  nacosClientConfig,
	})

	if err != nil {
		logs.Error(err)
	}
	content, err := configClient.GetConfig(param)
	if err != nil {
		panic("start app error: " + err.Error())
	}

	////lerr := configClient.ListenConfig(param)
	//if lerr != nil {
	//	logs.Error(err)
	//}
	parser := IniParser{}
	cfg.data, _ = parser.Parse(content)

	return cfg, err
	//if err != nil {
	//	logs.Error(err)
	//}
	//fmt.Println(content)
	//return newNacosConfiger(configClient, "example", "nacosconfig"), errors.New("")
}

func get(client config_client.IConfigClient, dataId string, groupId string) (string, error) {
	content, err := client.GetConfig(vo.ConfigParam{DataId: dataId, Group: groupId})
	if err != nil {
		logs.Error(err)
		return "", err
	}
	return content, err

}

func init() {
	fmt.Println("load nacos adapter")
	config.Register("nacos", &NacosConfigerProvider{})
	getNaconConfig()
}

func getNaconConfig() {

	section, _ := web.AppConfig.GetSection("nacos")

	if section == nil {
		logs.Error("没有获取相关信息")
		os.Exit(-1)
	}

	_, aerr := section["server_addr"]
	if !aerr {
		logs.Error("不存在nacos参数：server_addr")
		os.Exit(-1)
	}

	_, ok := section["data_id"]
	if !ok {
		logs.Error("不存在nacos参数：data_id")
		os.Exit(-1)
	}

	_, ok = section["group_id"]
	if !ok {
		logs.Error("不存在nacos参数：group_id")
		os.Exit(-1)
	}

	nacosServerAddr := section["server_addr"]
	nacosServerNamespaceId := section["namespaces_id"]

	var config, err = url.Parse(nacosServerAddr)
	if err != nil {
		logs.Error(err)
	}
	port := uint64(8848)
	host := strings.SplitN(config.Host, ":", 2)
	if len(host) == 2 {
		p, err := strconv.ParseUint(host[1], 10, 64)
		if err != nil {
			logs.Error(err)
		}
		port = p
	}
	nacosServerConfig = []constant.ServerConfig{
		{
			IpAddr:      host[0],
			ContextPath: config.Path,
			Port:        port,
			Scheme:      config.Scheme,
		},
	}

	nacosClientConfig = constant.ClientConfig{
		NamespaceId: nacosServerNamespaceId,
	}

	param = vo.ConfigParam{
		DataId: section["data_id"],
		Group:  section["group_id"],
	}
	logs.Info("解析成功")
}
