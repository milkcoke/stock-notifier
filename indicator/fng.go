package indicator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	// ExtremeFearThreshold 는 알람을 보내는 기준 지수입니다.
	ExtremeFearThreshold = 24
)

type IndexResult struct {
	Value  float64
	Status string
}

type fngResponse struct {
	Data []struct {
		Value               string `json:"value"`
		ValueClassification string `json:"value_classification"`
	} `json:"data"`
}

func GetFearAndGreed(client *http.Client) (*IndexResult, error) {
	resp, err := client.Get("https://api.alternative.me/fng/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res fngResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	val, _ := strconv.ParseFloat(res.Data[0].Value, 64)
	return &IndexResult{
		Value:  val,
		Status: res.Data[0].ValueClassification,
	}, nil
}
