package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var api = slack.New(os.Getenv("SLACK_AUTH_TOKEN"), slack.OptionDebug(true),
	slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))

func main() {
	http.HandleFunc("/hello/world/", helloWorld)

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8000", nil)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func query(city string) (weatherData, error) {
	weatherURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s", os.Getenv("API_WEATHER"), city)
	resp, err := http.Get(weatherURL)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}
