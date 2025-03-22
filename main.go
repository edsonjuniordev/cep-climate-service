package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"unicode"

	"github.com/joho/godotenv"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Response represents the final response structure.
type Response struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var (
	cepRegex = regexp.MustCompile(`^\d{8}$`)
)

// getLocation fetches the location data from ViaCEP API.
func getLocation(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("cannot find zipcode")
	}

	var data struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Erro {
		return "", errors.New("cannot find zipcode")
	}

	return data.Localidade, nil
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: Mark, Nonspacing
	}), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// getWeather fetches the weather data for a given city.
func getWeather(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to fetch weather data")
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}

// convertTemperatures converts temperatures from Celsius to Fahrenheit and Kelvin.
func convertTemperatures(tempC float64) Response {
	return Response{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !cepRegex.MatchString(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	city, err := getLocation(cep)
	if err != nil {
		if err.Error() == "cannot find zipcode" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("error fetching location: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	city = removeAccents(city)

	escapedCity := url.QueryEscape(city)

	tempC, err := getWeather(escapedCity)
	if err != nil {
		log.Printf("error fetching weather: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := convertTemperatures(tempC)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/weather", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
