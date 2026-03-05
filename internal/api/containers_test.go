package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListContainers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/workers/containers/namespaces" {
			t.Errorf("path = %q, want /accounts/test-account/workers/containers/namespaces", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": [
				{
					"id": "container-1",
					"name": "test-container-1",
					"cpu_ms": 1000,
					"memory_mb": 128
				},
				{
					"id": "container-2",
					"name": "test-container-2",
					"cpu_ms": 2000,
					"memory_mb": 256
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	containers, err := client.ListContainers("test-account")
	if err != nil {
		t.Fatalf("ListContainers() error = %v", err)
	}

	if len(containers) != 2 {
		t.Fatalf("len(containers) = %d, want 2", len(containers))
	}

	if containers[0].ID != "container-1" {
		t.Errorf("containers[0].ID = %q, want %q", containers[0].ID, "container-1")
	}

	if containers[0].Name != "test-container-1" {
		t.Errorf("containers[0].Name = %q, want %q", containers[0].Name, "test-container-1")
	}

	if containers[0].CPUMS != 1000 {
		t.Errorf("containers[0].CPUMS = %d, want 1000", containers[0].CPUMS)
	}
}

func TestGetContainer_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/workers/containers/namespaces/container-1" {
			t.Errorf("path = %q", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": {
				"id": "container-1",
				"name": "test-container",
				"cpu_ms": 1500,
				"memory_mb": 256
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	container, err := client.GetContainer("test-account", "container-1")
	if err != nil {
		t.Fatalf("GetContainer() error = %v", err)
	}

	if container.ID != "container-1" {
		t.Errorf("ID = %q, want %q", container.ID, "container-1")
	}

	if container.Name != "test-container" {
		t.Errorf("Name = %q, want %q", container.Name, "test-container")
	}
}

func TestListContainers_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "errors": [{"message": "Internal Server Error"}]}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	_, err := client.ListContainers("test-account")
	if err == nil {
		t.Fatal("ListContainers() with API error: error = nil, want error")
	}
}

func TestListContainers_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "result": []}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	containers, err := client.ListContainers("test-account")
	if err != nil {
		t.Fatalf("ListContainers() error = %v, want nil", err)
	}

	if len(containers) != 0 {
		t.Errorf("len(containers) = %d, want 0", len(containers))
	}
}

func TestGetContainer_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success": false, "errors": [{"message": "Container not found"}]}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	_, err := client.GetContainer("test-account", "nonexistent")
	if err == nil {
		t.Fatal("GetContainer() with API error: error = nil, want error")
	}
}
