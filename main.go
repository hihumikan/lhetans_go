package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
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

type Discord struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
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
	fmt.Println("Request body:", requestBody)
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
	c, err := maps.NewClient(maps.WithAPIKey("YOUR_API_KEY"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(route)
	return "1 hour", nil
}

func sendWebhookNotification(webhookURL string, duration string) error {

	var discord Discord
	discord.Username = "Mr. Hogehoge"
	discord.AvatarUrl = "https://github.com/qiita.png"
	discord.Content = "Hello World!"

	// encode json
	discord_json, _ := json.Marshal(discord)
	fmt.Println(string(discord_json))

	// discord webhook_url
	webhook_url := ""
	res, _ := http.Post(
		webhook_url,
		"application/json",
		bytes.NewBuffer(discord_json),
	)
	defer res.Body.Close()

	fmt.Println("Notification sent to webhook:", webhookURL)
	fmt.Println("Estimated duration:", duration)
	return nil
}
