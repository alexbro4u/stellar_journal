package nasa_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stellar_journal/internal/models/nasa_api"
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
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (a *NasaApi) handleResponse(resp *http.Response, target interface{}) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (a *NasaApi) GetAPOD() (*nasa_api.APODResp, error) {
	op := "internal/stellar_api/nasa_api/nasa_api.GetAPOD"

	client := a.createClient()

	requestUrl := a.Host + "/planetary/apod" + "?api_key=" + a.Token
	request, err := a.createRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var apodResp nasa_api.APODResp
	err = a.handleResponse(resp, &apodResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &apodResp, nil
}
