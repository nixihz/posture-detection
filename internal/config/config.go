package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Detector struct {
		ScaleFactor              float64 `yaml:"scale_factor"`
		MinNeighbors             int     `yaml:"min_neighbors"`
		MinFaceSize              int     `yaml:"min_face_size"`
		MaxFaceSize              int     `yaml:"max_face_size"`
		EnableSideView           bool    `yaml:"enable_side_view"`
		EnableHunchbackDetection bool    `yaml:"enable_hunchback_detection"`
		HunchbackAngleThreshold  float64 `yaml:"hunchback_angle_threshold"`
		MinSitDistance           int     `yaml:"min_sit_distance"`
		MaxSitDistance           int     `yaml:"max_sit_distance"`
		MinSitHeight             int     `yaml:"min_sit_height"`
		MaxSitHeight             int     `yaml:"max_sit_height"`
	} `yaml:"detector"`
	Camera struct {
		Width        int  `yaml:"width"`
		Height       int  `yaml:"height"`
		FPS          int  `yaml:"fps"`
		Autofocus    bool `yaml:"autofocus"`
		Autoexposure bool `yaml:"autoexposure"`
		Brightness   int  `yaml:"brightness"`
		Contrast     int  `yaml:"contrast"`
	} `yaml:"camera"`
	Notification struct {
		Enable   bool `yaml:"enable"`
		Interval int  `yaml:"interval"`
	} `yaml:"notification"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{}
		data, err := ioutil.ReadFile("config/config.yaml")
		if err != nil {
			log.Fatalf("读取配置文件失败: %v", err)
		}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			log.Fatalf("解析配置文件失败: %v", err)
		}
	}
	return config
}
