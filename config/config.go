package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

//Config ...
type Config struct {
	Database struct {
		CSVFile string `yaml:"csv_file"`
	} `yaml:"database"`
	WEBServer struct {
		Port      string `yaml:"port"`
		Host      string `yaml:"host"`
		CertsFile string `yaml:"certs_file"`
		KeyFile   string `yaml:"key_file"`
	} `yaml:"web_server"`

	APIServer struct {
		Port      string `yaml:"port"`
		Host      string `yaml:"host"`
		CertsFile string `yaml:"certs_file"`
		KeyFile   string `yaml:"key_file"`
	} `yaml:"api_server"`
	SearchResults struct {
		Number int     `yaml:"number"`
		Radius float64 `yaml:"radius"`
	} `yaml:"search_results"`
	MapsAPIKEY string `yaml:"maps_api_key"`
}

//NewConfig ...
func NewConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	cfg := &Config{}
	yd := yaml.NewDecoder(file)
	err = yd.Decode(cfg)

	if err != nil {
		return nil, err
	}
	return cfg, nil
}
