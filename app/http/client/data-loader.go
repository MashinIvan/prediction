package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/app/config"
	"net/http"
)

type AggregateScheme struct {
	Blocks      *[]interface{} `json:"blocks"`
	DataOptions []interface{}  `json:"data_options"`
	Limit       *int           `json:"limit"`
	Offset      *int           `json:"offset"`
}

func LoadData(data AggregateScheme) (interface{}, error) {
	url := fmt.Sprintf("%s://%s:%s/%s",
		config.Config.CoreScheme,
		config.Config.CoreHost,
		config.Config.CorePort,
		"data/dataset/")

	jsonData, err := json.Marshal(data)

	result, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	return result, nil
}
