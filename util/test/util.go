package test

import (
	"encoding/json"
	"io/ioutil"
)

// LoadDataFromJSON load data from a json file
func LoadDataFromJSON(path string, obj interface{}) error {
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonFile, obj)
}

// ConvertStructToMap to convert a struct to a map
func ConvertStructToMap(s interface{}) (map[string]interface{}, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
