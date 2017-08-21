package common

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/spf13/viper"
)

type ActionsList struct {
	Actions []Action
}

var defaultList *ActionsList

func DefaultActionsList() *ActionsList {
	if defaultList == nil {
		path := viper.GetString("config")
		defaultList, _ = NewActionsListFromPath(path)
	}

	return defaultList
}

func NewActionsListFromPath(path string) (*ActionsList, error) {
	list := &ActionsList{}

	data, err := ioutil.ReadFile(path)
	if err != nil { return nil, err}

	err = yaml.Unmarshal(data, list)
	if err != nil { return nil, err}

	return list, nil
}

func (a *ActionsList) Add(typ, action string) {
	a.Actions = append(a.Actions, Action{typ: typ, action: action})
}

func (a *ActionsList) WriteToPath(path string) error {
	data, err := yaml.Marshal(a)
	if err != nil { return err }

	err = ioutil.WriteFile(path, data, 0600)
	if err != nil { return err }

	return nil
}

func (a *ActionsList) MarshalYAML() (interface{}, error) {
	mapList := []map[string]string{}

	for _, action := range a.Actions {
		item := map[string]string{}
		item[action.typ] = action.action
		mapList = append(mapList, item)
	}

	return mapList, nil
}

func (a *ActionsList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	mapList := []map[string]string{}

	err := unmarshal(&mapList)
	if err != nil { return err }

	for _, listItem := range mapList {
		for key, val := range listItem {
			action := Action{typ: key, action: val}
			a.Actions = append(a.Actions, action)
		}
	}

	return nil
}