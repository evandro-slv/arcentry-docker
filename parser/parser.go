package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
	"time"
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

func parseUnit(str string, dur time.Duration, unit string) (time.Duration, error) {
	i, err := strconv.Atoi(strings.Replace(str, unit, "", -1))

	if err != nil {
		return 5 * time.Second, err
	}

	return time.Duration(i) * dur, nil
}

func ParseDuration(duration string) (dur time.Duration, err error) {
	if strings.Contains(duration, "s") {
		return parseUnit(duration, time.Second, "s")
	} else if strings.Contains(duration, "m") {
		return parseUnit(duration, time.Minute, "m")
	} else if strings.Contains(duration, "h") {
		return parseUnit(duration, time.Hour, "h")
	}

	return 5 * time.Second, errors.New(fmt.Sprintf("invalid duration format: %s", duration))
}
