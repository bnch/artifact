package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Conf is a configuration file of artifact.
type Conf struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBName     string
	DBType     string
	ListenTo   string
}

// NewConf gets the configuration file for artifact from conf.json, or creates it.
func NewConf() Conf {
	data, err := ioutil.ReadFile("conf.json")
	if os.IsNotExist(err) {
		defaultConfig := Conf{
			DBUsername: "root",
			DBName:     "artifact",
			DBType:     "mysql",
		}
		b, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("conf.json", b, 0644)
		if err != nil {
			panic(err)
		}
		data, _ = ioutil.ReadFile("conf.json")
	} else if err != nil {
		panic(err)
	}
	c := Conf{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	return c
}

func (c Conf) Save() {
	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("conf.json", b, 0644)
	if err != nil {
		panic(err)
	}
}
