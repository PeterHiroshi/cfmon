package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWorkers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/workers/scripts" {
			t.Errorf("path = %q, want /accounts/test-account/workers/scripts", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": [
				{
					"id": "worker-1",
					"name": "test-worker-1",
					"cpu_ms": 500,
					"requests": 1000
				},
				{
					"id": "worker-2",
					"name": "test-worker-2",
					"cpu_ms": 750,
					"requests": 2000
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	workers, err := client.ListWorkers("test-account")
	if err != nil {
		t.Fatalf("ListWorkers() error = %v", err)
	}

	if len(workers) != 2 {
		t.Fatalf("len(workers) = %d, want 2", len(workers))
	}

	if workers[0].ID != "worker-1" {
		t.Errorf("workers[0].ID = %q, want %q", workers[0].ID, "worker-1")
	}

	if workers[0].Name != "test-worker-1" {
		t.Errorf("workers[0].Name = %q, want %q", workers[0].Name, "test-worker-1")
	}

	if workers[0].CPUMS != 500 {
		t.Errorf("workers[0].CPUMS = %d, want 500", workers[0].CPUMS)
	}

	if workers[0].Requests != 1000 {
		t.Errorf("workers[0].Requests = %d, want 1000", workers[0].Requests)
	}
}

func TestListWorkers_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "errors": [{"message": "Internal Server Error"}]}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	_, err := client.ListWorkers("test-account")
	if err == nil {
		t.Fatal("ListWorkers() with API error: error = nil, want error")
	}
}

func TestListWorkers_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "result": []}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	workers, err := client.ListWorkers("test-account")
	if err != nil {
		t.Fatalf("ListWorkers() error = %v, want nil", err)
	}

	if len(workers) != 0 {
		t.Errorf("len(workers) = %d, want 0", len(workers))
	}
}
