package main

import (
	"fmt"
	"net/http"
	"stock-notifier/config"
	"stock-notifier/indicator"
	"stock-notifier/notifier"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	client := &http.Client{Timeout: 10 * time.Second}

	idx, err := indicator.GetFearAndGreed(client)
	if err != nil {
		fmt.Printf("❌ Failed to get API: %v\n", err)
		return
	}
	fmt.Printf("📊 Fear and Greed Index: %.0f (%s)\n", idx.Value, idx.Status)

	if idx.Value <= indicator.ExtremeFearThreshold {
		fmt.Println("🚨 Extreme Fear Detected!")

		notion := notifier.NewNotionNotifier(client, cfg)
		if err := notion.SendAlarm(idx.Value, idx.Status); err != nil {
			fmt.Printf("❌ Failed to send notification: %v\n", err)
		} else {
			fmt.Println("✅ Notification sent successfully")
		}
	} else {
		fmt.Println("✅ Market is stable. Skip the alarm.")
	}
}
