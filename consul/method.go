package consul

import (
	"encoding/json"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

func SetJobNameValue(key string, value ConsulKVData) error{
	var newData []ConsulKVData
	if oldData := GetJobNameData(key); oldData != nil {
		newData = append(oldData, value)
	} else {
		newData = append(newData, value)
	}
	data, _ := json.Marshal(newData)
	_, err := ConsulClient.KV().Put(&consulapi.KVPair{
		Key:         key,
		Value:       data,
	}, nil)
	return err
}

func GetJobNameData(key string) []ConsulKVData {
	data := make([]ConsulKVData, 0)
	kv, _, err := ConsulClient.KV().Get(key, nil)
	if err != nil || kv == nil{
		return nil
	}
	err = json.Unmarshal(kv.Value, &data)
	if err != nil {
		return nil
	}
	return data
}

func DeleteJobNameData(key, name string) error {
	datas := GetJobNameData(key)
	if datas == nil {
		return fmt.Errorf("JobName 不存在")
	}
	newDatas := make([]ConsulKVData, 0)
	for _, data:= range datas {

		if data.Name == name {
			continue
		}
		newDatas = append(newDatas, data)
	}

	err := UpdateValue(key, newDatas)
	if err != nil {
		return err
	}
	return nil
}