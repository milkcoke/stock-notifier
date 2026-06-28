package indicator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ExtremeFearThreshold = 25
)

type IndexResult struct {
	Value  float64
	Status string
}

type cnnFngResponse struct {
	FearAndGreed struct {
		Score  float64 `json:"score"`
		Rating string  `json:"rating"`
	} `json:"fear_and_greed"`
}

func GetFearAndGreed(client *http.Client) (*IndexResult, error) {
	req, err := http.NewRequest("GET", "https://production.dataviz.cnn.io/index/fearandgreed/graphdata", nil)
	if err != nil {
		return nil, err
	}

	// 브라우저 요청처럼 위장
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://edition.cnn.com/markets/fear-and-greed")
	req.Header.Set("Origin", "https://edition.cnn.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CNN API returned status %d", resp.StatusCode)
	}

	var res cnnFngResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.FearAndGreed.Score == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return &IndexResult{
		Value:  res.FearAndGreed.Score,
		Status: res.FearAndGreed.Rating,
	}, nil
}
