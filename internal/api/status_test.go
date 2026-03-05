package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/user/tokens/verify" {
			t.Errorf("Path = %q, want /user/tokens/verify", r.URL.Path)
		}

		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": {
				"id": "test-token-id",
				"status": "active"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	status, err := client.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	if !status.Valid {
		t.Error("Valid = false, want true")
	}

	if status.Status != "active" {
		t.Errorf("Status = %q, want %q", status.Status, "active")
	}
}

func TestClient_GetStatus_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success": false, "errors": [{"message": "Invalid token"}]}`))
	}))
	defer server.Close()

	client := NewClient("invalid-token")
	client.baseURL = server.URL

	status, err := client.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus() error = %v, want nil", err)
	}

	if status.Valid {
		t.Error("Valid = true, want false")
	}
}

func TestClient_GetAccountInfo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/accounts" {
			t.Errorf("Path = %q, want /accounts", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": [
				{
					"id": "test-account-id",
					"name": "Test Account",
					"settings": {
						"enforce_twofactor": false
					},
					"legacy_flags": {
						"enterprise_zone_quota": {
							"maximum": 5000,
							"current": 100,
							"available": 4900
						}
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	accountInfo, err := client.GetAccountInfo()
	if err != nil {
		t.Fatalf("GetAccountInfo() error = %v", err)
	}

	if accountInfo.ID != "test-account-id" {
		t.Errorf("ID = %q, want %q", accountInfo.ID, "test-account-id")
	}

	if accountInfo.Name != "Test Account" {
		t.Errorf("Name = %q, want %q", accountInfo.Name, "Test Account")
	}

	if accountInfo.PlanType != "Enterprise" {
		t.Errorf("PlanType = %q, want %q", accountInfo.PlanType, "Enterprise")
	}
}

func TestClient_GetAccountInfo_FreePlan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": [
				{
					"id": "free-account-id",
					"name": "Free Account",
					"settings": {
						"enforce_twofactor": false
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	accountInfo, err := client.GetAccountInfo()
	if err != nil {
		t.Fatalf("GetAccountInfo() error = %v", err)
	}

	if accountInfo.PlanType != "Free" {
		t.Errorf("PlanType = %q, want %q", accountInfo.PlanType, "Free")
	}
}

func TestClient_GetAccountInfo_NoAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": []
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	_, err := client.GetAccountInfo()
	if err == nil {
		t.Fatal("GetAccountInfo() error = nil, want error")
	}
}
