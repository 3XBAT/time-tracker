package api

import (
	"encoding/json"
	"fmt"
	"github.com/3XBAT/time-tracker/internal/config"
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"net/http"
)

type ApiClient struct {
	config *config.Config
}

func NewApiClient(config *config.Config) *ApiClient {
	return &ApiClient{config: config}
}

func (api *ApiClient) UserInfo(passportNumber, passportSeries string) (*models.User, error) {
	const op = "api.UserInfo"
	url := fmt.Sprintf("%s/info?passportSerie=%s&passportNumber=%s", api.config.API.ExternalURL, passportSeries, passportNumber)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("%s bad request: %s", op, resp.Status)
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %s: the server can not find the requested resource", op, resp.Status)
		}
		return nil, fmt.Errorf("%s failed: %s", op, resp.Status)

	}
	var user models.User

	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
