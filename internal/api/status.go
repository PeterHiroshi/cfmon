package api

import (
	"fmt"
	"net/http"
)

// TokenStatus represents the status of a Cloudflare API token
type TokenStatus struct {
	Valid  bool
	Status string
}

// AccountInfo represents information about a Cloudflare account
type AccountInfo struct {
	ID       string
	Name     string
	PlanType string
}

// GetStatus verifies the API token and returns its status
func (c *Client) GetStatus() (*TokenStatus, error) {
	var response struct {
		Success bool `json:"success"`
		Result  struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"result"`
	}

	err := c.doRequest("GET", "/user/tokens/verify", &response)
	if err != nil {
		// Check if it's an unauthorized error
		if resp, ok := err.(interface{ StatusCode() int }); ok && resp.StatusCode() == http.StatusUnauthorized {
			return &TokenStatus{Valid: false}, nil
		}
		// For unauthorized errors from API, return invalid token
		return &TokenStatus{Valid: false}, nil
	}

	return &TokenStatus{
		Valid:  response.Success,
		Status: response.Result.Status,
	}, nil
}

// GetAccountInfo retrieves information about the Cloudflare account
func (c *Client) GetAccountInfo() (*AccountInfo, error) {
	var response struct {
		Success bool `json:"success"`
		Result  []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Settings struct {
				EnforceTwoFactor bool `json:"enforce_twofactor"`
			} `json:"settings"`
			LegacyFlags struct {
				EnterpriseZoneQuota struct {
					Maximum   int `json:"maximum"`
					Current   int `json:"current"`
					Available int `json:"available"`
				} `json:"enterprise_zone_quota"`
			} `json:"legacy_flags"`
		} `json:"result"`
	}

	err := c.doRequest("GET", "/accounts", &response)
	if err != nil {
		return nil, err
	}

	if len(response.Result) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}

	account := response.Result[0]
	planType := "Free"

	// Determine plan type based on enterprise zone quota
	if account.LegacyFlags.EnterpriseZoneQuota.Maximum > 0 {
		planType = "Enterprise"
	}

	return &AccountInfo{
		ID:       account.ID,
		Name:     account.Name,
		PlanType: planType,
	}, nil
}
