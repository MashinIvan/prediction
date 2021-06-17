package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main"
	"net/http"

	app_types "gitlab.ptnl.moscow/polymatica/common/backend/app-types"
)

type AggregateScheme struct {
	Blocks      *[]app_types.DataOptionBlock `json:"blocks"`
	DataOptions []app_types.DataOption       `json:"data_options"`
	Limit       *int                         `json:"limit"`
	Offset      *int                         `json:"offset"`
}

func LoadData(data AggregateScheme) (interface{}, error) {
	url := fmt.Sprintf("%s://%s:%s/%s",
		main.Config.CoreScheme,
		main.Config.CoreHost,
		main.Config.CorePort,
		"data/dataset/")

	jsonData, err := json.Marshal(data)

	result, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	return result, nil
}
