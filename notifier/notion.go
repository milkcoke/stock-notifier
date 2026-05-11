package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stock-notifier/config"
	"time"
)

type NotionNotifier struct {
	Client *http.Client
	Config *config.Config
}

func NewNotionNotifier(client *http.Client, cfg *config.Config) *NotionNotifier {
	return &NotionNotifier{Client: client, Config: cfg}
}

func (n *NotionNotifier) SendAlarm(value float64, status string) error {
	currentTime := time.Now().In(n.Config.Location).Format("2006-01-02 15:04")
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", n.Config.PageID)

	payload := map[string]interface{}{
		"children": []interface{}{
			map[string]interface{}{
				"object": "block",
				"type":   "paragraph",
				"paragraph": map[string]interface{}{
					"rich_text": []interface{}{
						map[string]interface{}{
							"type": "mention",
							"mention": map[string]interface{}{
								"type": "user",
								"user": map[string]interface{}{"id": n.Config.UserID},
							},
						},
						map[string]interface{}{
							"type": "text",
							"text": map[string]string{
								"content": fmt.Sprintf(" [%s] 🚨 지수: %.0f (%s) - 극도의 공포 단계입니다!", currentTime, value, status),
							},
						},
					},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+n.Config.NotionToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	resp, err := n.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
