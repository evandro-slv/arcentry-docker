package api

import (
	"bytes"
	"github.com/evandro-slv/arcentry-docker/parser"
	"io/ioutil"
	"log"
	"net/http"
)

func UpdateDoc(arcentry parser.Arcentry, b []byte) error {
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://arcentry.com/api/v1/doc/"+arcentry.DocId, bytes.NewReader(b))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+arcentry.ApiKey)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	log.Println(string(body))

	return nil
}
