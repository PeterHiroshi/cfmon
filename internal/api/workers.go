package api

import "fmt"

// Worker represents a Cloudflare worker
type Worker struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	CPUMS       int     `json:"cpu_ms"`
	Requests    int     `json:"requests"`
	Errors      int     `json:"errors,omitempty"`
	Status      string  `json:"status,omitempty"`
	SuccessRate float64 `json:"success_rate,omitempty"`
}

type workersResponse struct {
	Success bool     `json:"success"`
	Result  []Worker `json:"result"`
}

// ListWorkers lists all workers for an account
func (c *Client) ListWorkers(accountID string) ([]Worker, error) {
	path := fmt.Sprintf("/accounts/%s/workers/scripts", accountID)

	var resp workersResponse
	if err := c.doRequest("GET", path, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}
