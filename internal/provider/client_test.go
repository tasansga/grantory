package provider

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrantoryClientSetsAuthorizationHeader(t *testing.T) {
	t.Parallel()

	const token = "secret-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer "+token, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &grantoryClient{
		baseURL:    mustParseURL(t, server.URL),
		httpClient: server.Client(),
		token:      token,
	}

	assert.NoError(t, client.doJSON(context.Background(), http.MethodGet, "/auth", nil, nil))
}

func TestGrantoryClientSetsBasicAuth(t *testing.T) {
	t.Parallel()

	const user = "alice"
	const password = "s3cr3t"
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expected, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &grantoryClient{
		baseURL:    mustParseURL(t, server.URL),
		httpClient: server.Client(),
		user:       user,
		password:   password,
	}

	assert.NoError(t, client.doJSON(context.Background(), http.MethodGet, "/auth", nil, nil))
}
