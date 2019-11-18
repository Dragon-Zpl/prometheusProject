package consul

import (
	"PrometheusProject/conf"
	"encoding/json"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"strings"
)
var ConsulClient *consulapi.Client
func init() {
	var err error
	conf.InitConf()
	config := consulapi.DefaultConfig()
	config.Address = conf.Consul.Url
	ConsulClient, err = consulapi.NewClient(config)
	if err != nil {
		panic(err)
	}
}

func ConsulRegisterServer(serverId, serverName, serverIp, serverTag, serverPath string, serverPort int) error {
	checkPort := serverPort
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = serverIp + "_" + serverId
	registration.Name = serverName
	registration.Port = serverPort
	if serverTag != "" {
		registration.Tags = []string{serverTag}
	}
	registration.Address = serverIp
	if serverPath == "" {
		serverPath = "metrics"
	}
	if !strings.Contains(serverPath, "/") {
		serverPath = "/" + serverPath
	}
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, serverPath),
		Timeout:                        "3s",
		Interval:                       "5s", // 健康检查间隔
		//DeregisterCriticalServiceAfter: "1h", //check失败后30秒删除本服务
	}
	err := ConsulClient.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
		return err
	}
	return nil
}



func ConsulDeRegisterServer(serverId, serverIp string) error {
	return ConsulClient.Agent().ServiceDeregister(serverIp + "_" + serverId)
}

func IsExistService(serverId, serverIp string) bool {
	services, err := ConsulClient.Agent().Services()
	if err != nil {
		return false
	}
	if _, found := services[serverIp + "_" + serverId]; !found {
		return false
	}
	return true
}

func GetConsulRegisterService() []string {
	services, err := ConsulClient.Agent().Services()
	if err != nil {
		return []string{}
	}
	servers := make([]string, 0)
	for v, _ := range services{
		servers = append(servers, v)
	}
	return servers
}

type ConsulKVData struct {
	Name string
	Labels []string
	Type string
}

func GetKeyData(key string) []byte {
	kv, _, err := ConsulClient.KV().Get(key, nil)
	if err != nil || kv == nil {
		return nil
	}
	return kv.Value
}

func SetKeyValue(key string, value interface{}) error {
	data, _ := json.Marshal(value)
	_, err := ConsulClient.KV().Put(&consulapi.KVPair{
		Key:         key,
		Value:       data,
	}, nil)
	return err
}

func UpdateValue(key string, value []ConsulKVData) error {
	data, _ := json.Marshal(value)
	_, err := ConsulClient.KV().Put(&consulapi.KVPair{
		Key:         key,
		Value:       data,
	}, nil)
	return err
}

func DeleteData(key string) error {
	_, err := ConsulClient.KV().Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}