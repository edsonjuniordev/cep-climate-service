package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConvertTemperatures(t *testing.T) {
	tempC := 25.0
	expected := Response{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}
	result := convertTemperatures(tempC)

	if result != expected {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestInvalidCEP(t *testing.T) {
	if !cepRegex.MatchString("123") {
		return
	}
	t.Errorf("Expected invalid CEP, but passed")
}

func TestValidCEP(t *testing.T) {
	if !cepRegex.MatchString("12345678") {
		t.Errorf("Expected valid CEP, but failed")
	}
}

// TestRemoveAccents tests the removeAccents function.
func TestRemoveAccents(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"São Paulo", "Sao Paulo"},
		{"Curitiba", "Curitiba"},
		{"Belo Horizonte", "Belo Horizonte"},
		{"João Pessoa", "Joao Pessoa"},
	}

	for _, c := range cases {
		result := removeAccents(c.input)
		if result != c.expected {
			t.Errorf("removeAccents(%q) == %q, want %q", c.input, result, c.expected)
		}
	}
}

// TestGetLocation tests the getLocation function.
func TestGetLocation(t *testing.T) {
	// Simulate the response from ViaCEP
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cep") == "01001000" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"localidade": "São Paulo"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	viaCEPBaseURL := "https://viacep.com.br/ws"

	originalViaCEP := viaCEPBaseURL

	viaCEPBaseURL = mockServer.URL
	defer func() { viaCEPBaseURL = originalViaCEP }()

	cases := []struct {
		cep      string
		expected string
		wantErr  bool
	}{
		{"01001000", "São Paulo", false},
		{"00000000", "", true},
	}

	for _, c := range cases {
		result, err := getLocation(c.cep)
		if (err != nil) != c.wantErr {
			t.Errorf("getLocation(%q) error = %v, wantErr %v", c.cep, err, c.wantErr)
			continue
		}
		if result != c.expected {
			t.Errorf("getLocation(%q) == %q, want %q", c.cep, result, c.expected)
		}
	}
}
