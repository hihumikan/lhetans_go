package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	TravelMode      string `json:"travel_mode"`
	WebhookURL      string `json:"webhook_url"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	duration, err := getRouteInfo(requestBody.HomeLocation, requestBody.CurrentLocation, requestBody.TravelMode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go sendWebhookNotification(requestBody.WebhookURL, duration, requestBody.HomeLocation, requestBody.CurrentLocation)

	fmt.Println("Notification sent to webhook:", requestBody.WebhookURL)
	w.WriteHeader(http.StatusOK)
}

func getRouteInfo(homeLocation string, currentLocation string, travelMode string) (string, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not found")
	}

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}

	r := &maps.DirectionsRequest{
		Origin:        homeLocation,
		Destination:   currentLocation,
		DepartureTime: "now",
		Mode:          maps.Mode(travelMode),
	}

	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		return "", err
	}

	if len(route) == 0 || len(route[0].Legs) == 0 {
		return "", fmt.Errorf("no route found")
	}

	duration := route[0].Legs[0].Duration.String()
	return duration, nil
}

func sendWebhookNotification(webhookURL string, duration string, homeLocation string, currentLocation string) {
	var discord Discord
	discord.Username = "Google Maps API"
	discord.AvatarUrl = "https://asset.watch.impress.co.jp/img/ktw/docs/1238/736/icon_l.png"
	discord.Content = fmt.Sprintf("Estimated duration: %s\nHome location: %s\nCurrent location: %s", duration, homeLocation, currentLocation)

	discordJSON, _ := json.Marshal(discord)

	go func() {
		res, err := http.Post(
			webhookURL,
			"application/json",
			bytes.NewBuffer(discordJSON),
		)
		if err != nil {
			log.Printf("Error sending webhook notification: %v", err)
			return
		}
		defer res.Body.Close()
		fmt.Println("Notification sent to webhook:", webhookURL)
	}()
}
