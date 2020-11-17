package global

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Config ServerConfig
)


type ServerConfig struct {
	Capacity int `json:"capacity"`
	Evict string `json:"evict"`
}

func init() {
	jsonFile, err := os.Open(".gocache.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Config)
	log.Println(Config)
}
