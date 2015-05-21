package service

import (
	"encoding/xml"
	"log"
	"os"
)

// Config ...
var Config AppConfig

// AppConfig ...
type AppConfig struct {
	Server struct {
		Host string `xml:"host"`
		Port string `xml:"port"`
	} `xml:"server"`
	Logger         LoggerConfig         `xml:"logger"`
	SessionStorage SessionStorageConfig `xml:"storage"`
}

// InitConfig ...
func InitConfig(path, environment string) {
	file, err := os.Open(path + "/" + environment + ".xml")
	if err != nil {
		log.Fatalf("Cannot open configuration file: %v\n", err)
	}

	defer file.Close()
	decoder := xml.NewDecoder(file)

	decoder.Decode(&Config)
}
