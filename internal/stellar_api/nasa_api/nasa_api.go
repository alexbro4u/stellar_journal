package nasa_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stellar_journal/internal/models/nasa_api_models"
	"time"
)

type NasaApi struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

func NewNasaApiConnect(host, token string) *NasaApi {
	return &NasaApi{
		Host:  host,
		Token: token,
	}
}

func (a *NasaApi) createClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (a *NasaApi) createRequest(method, url string, body []byte) (*http.Request, error) {
	const op = "internal/stellar_api/nasa_api.createRequest"

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (a *NasaApi) doRequest(req *http.Request, target interface{}) error {
	const op = "internal/stellar_api/nasa_api.doRequest"

	client := a.createClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: failed to send request: %w", op, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%s: failed to close response body: %w", op, err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: unexpected status code: %d", op, resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (a *NasaApi) GetAPOD() (*nasa_api_models.APODResp, error) {
	const op = "internal/stellar_api/nasa_api.GetAPOD"

	url := fmt.Sprintf("%s/planetary/apod?api_key=%s", a.Host, a.Token)
	req, err := a.createRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	var apodResp nasa_api_models.APODResp
	err = a.doRequest(req, &apodResp)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to do request: %w", op, err)
	}

	return &apodResp, nil
}
