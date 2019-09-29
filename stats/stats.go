package stats

import (
	"github.com/evandro-slv/arcentry-docker/parser"
	"os/exec"
	"strings"
)

func fetch() (string, error) {
	args := "stats -a --no-trunc --no-stream --format {\"container\":\"{{.Container}}\",\"memory\":{\"raw\":\"{{.MemUsage}}\",\"percent\":\"{{.MemPerc}}\"},\"cpu\":\"{{.CPUPerc}}\"}"
	split := strings.Split(args, " ")
	out, err := exec.Command("docker", split...).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func GetStats() ([]*parser.Stat, error) {
	s, err := fetch()

	if err != nil {
		return nil, err
	}

	split := strings.Split(s, "\n")

	stats := make([]*parser.Stat, len(split))

	for i, js := range split {
		if js != "" {
			s, err := parser.ReadStat(js)

			if err != nil {
				return nil, err
			}

			stats[i] = s
		}
	}

	return stats, nil
}
