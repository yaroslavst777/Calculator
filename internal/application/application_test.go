package application

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalcHandlerSuccessCase(t *testing.T) {
	expected := "{\"result\":4}\n"
	body := strings.NewReader("{\"expression\": \"2+2\"}")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", body)
	w := httptest.NewRecorder()
	CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	sd := string(data)
	if sd != expected {
		t.Errorf("Expected %v but got %v", expected, string(data))
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("wrong status code")
	}
}
func TestCalcHandlerBadRequestCase(t *testing.T) {
	expected := "{\"error\":\"Expression is not valid\"}\n"
	body := strings.NewReader("{\"expression\": \"2+f2\"}")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", body)
	w := httptest.NewRecorder()
	CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	sd := string(data)
	if sd != expected {
		t.Errorf("Expected %v but got %v", expected, string(data))
	}

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("wrong status code")
	}
}
