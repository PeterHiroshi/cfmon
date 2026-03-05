package api

import "fmt"

// Container represents a Cloudflare container
type Container struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CPUMS    int    `json:"cpu_ms"`
	MemoryMB int    `json:"memory_mb"`
}

type containersResponse struct {
	Success bool        `json:"success"`
	Result  []Container `json:"result"`
}

type containerResponse struct {
	Success bool      `json:"success"`
	Result  Container `json:"result"`
}

// ListContainers lists all containers for an account
func (c *Client) ListContainers(accountID string) ([]Container, error) {
	path := fmt.Sprintf("/accounts/%s/workers/containers/namespaces", accountID)

	var resp containersResponse
	if err := c.doRequest("GET", path, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetContainer gets a specific container by ID
func (c *Client) GetContainer(accountID, containerID string) (*Container, error) {
	path := fmt.Sprintf("/accounts/%s/workers/containers/namespaces/%s", accountID, containerID)

	var resp containerResponse
	if err := c.doRequest("GET", path, &resp); err != nil {
		return nil, err
	}

	return &resp.Result, nil
}
