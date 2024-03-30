package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GoogleMapsAPIResponse struct {
	Routes []struct {
		Legs []struct {
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
		} `json:"legs"`
	} `json:"routes"`
}

type RequestBody struct {
	HomeLocation    string `json:"home_location"`
	CurrentLocation string `json:"current_location"`
	TravelMode      int    `json:"travel_mode"`
	WebhookURL      string `json:"webhook_url"`
}

func main() {
	http.HandleFunc("/notification", handleNotification)
	http.ListenAndServe(":3000", nil)
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	duration, err := getRouteDuration(requestBody.HomeLocation, requestBody.CurrentLocation, requestBody.TravelMode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sendWebhookNotification(requestBody.WebhookURL, duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Notification sent to webhook:", requestBody.WebhookURL)
	w.WriteHeader(http.StatusOK)
}

func getRouteDuration(homeLocation string, currentLocation string, travelMode int) (string, error) {
	// Google Maps APIを呼び出して経路の所要時間を取得する

	// ここにGoogle Maps APIを呼び出すコードを書く

	// 仮のダミーデータを返す（実際にはGoogle Maps APIから取得する）
	return "1 hour", nil
}

func sendWebhookNotification(webhookURL string, duration string) error {
	// Webhookに通知を送信する

	// ここにWebhookに通知を送るコードを書く

	// 仮のダミーデータを返す（実際にはWebhookに通知を送る）
	fmt.Println("Notification sent to webhook:", webhookURL)
	fmt.Println("Estimated duration:", duration)
	return nil
}
