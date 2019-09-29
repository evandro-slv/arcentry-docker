package parser

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type Item struct {
	Cpu    string
	Memory string
}

type Container struct {
	Chart Item
	Text  Item
}

type Arcentry struct {
	ApiKey string `yaml:"apiKey"`
	DocId  string `yaml:"docId"`
}

type Config struct {
	Config struct {
		Arcentry Arcentry
		Watch    struct {
			Interval string
		}
		Containers []map[string]Container
	}
}

type Stat struct {
	Container string
	Memory    struct {
		Raw     string
		Percent string
	}
	Cpu string
}

type ChartDocument struct {
	Data       [][]float64 `json:"data"`
	MaxY       float64     `json:"maxY"`
	MinY       float64     `json:"minY"`
	YAxisSpace string      `json:"yAxisSpace"`
}

type TextDocument struct {
	Text string `json:"text"`
}

type Request struct {
	Objects map[string]interface{} `json:"objects"`
}

func ReadConfig(bytes []byte) (config *Config, err error) {
	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadStat(content string) (stat *Stat, err error) {
	err = json.Unmarshal([]byte(content), &stat)

	if err != nil {
		return nil, err
	}

	return stat, nil
}

func ParseRequest(content interface{}) (req []byte, err error) {
	bytes, err := json.Marshal(content)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
