package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clienttest "github.com/tasansga/terraform-provider-grantory/api/client/testutil"
)

func TestDataHostSource(t *testing.T) {
	t.Parallel()

	handler := newHostDataSourceTestHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	client := clienttest.New(t, server, "", "", "")

	t.Run("lookup by host_id", func(t *testing.T) {
		resource := dataHost()
		data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
			"host_id": "host-123",
		})

		assert.False(t, resource.ReadContext(context.Background(), data, client).HasError(), "unexpected diagnostics from host data read")

		assert.Equal(t, "host-123", data.Get("host_id"), "host_id should be preserved")
		assert.Equal(t, "unique:host", data.Get("unique_key"), "unique_key should be preserved")
		labelsValue, ok := data.Get("labels").(map[string]any)
		assert.True(t, ok, "labels should be a map")
		assert.Equal(t, map[string]any{"env": "prod"}, labelsValue, "labels should match stored data")
	})

	t.Run("lookup by unique_key success", func(t *testing.T) {
		resource := dataHost()
		data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
			"unique_key": "unique:host",
		})

		assert.False(t, resource.ReadContext(context.Background(), data, client).HasError(), "unexpected diagnostics from host data read")
		assert.Equal(t, "host-123", data.Id())
		assert.Equal(t, "host-123", data.Get("host_id"))
		assert.Equal(t, "unique:host", data.Get("unique_key"))
	})

	t.Run("lookup by unique_key not found", func(t *testing.T) {
		resource := dataHost()
		data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
			"unique_key": "nonexistent",
		})

		assert.False(t, resource.ReadContext(context.Background(), data, client).HasError())
		assert.Equal(t, "", data.Id())
	})

	t.Run("neither host_id nor unique_key provided", func(t *testing.T) {
		resource := dataHost()
		data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{})

		diags := resource.ReadContext(context.Background(), data, client)
		assert.True(t, diags.HasError())
		assert.Equal(t, "either host_id or unique_key must be specified", diags[0].Summary)
	})

	t.Run("both host_id and unique_key provided", func(t *testing.T) {
		resource := dataHost()
		data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
			"host_id":    "host-123",
			"unique_key": "unique:host",
		})

		diags := resource.ReadContext(context.Background(), data, client)
		assert.True(t, diags.HasError())
		assert.Equal(t, "cannot specify both host_id and unique_key", diags[0].Summary)
	})
}

func newHostDataSourceTestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createdAt := time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)
		sampleHost := apiHost{
			ID:        "host-123",
			UniqueKey: "unique:host",
			Labels:    map[string]string{"env": "prod"},
			CreatedAt: createdAt,
		}

		if r.Method == http.MethodGet && r.URL.Path == "/hosts" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]apiHost{sampleHost})
			return
		}

		if r.Method != http.MethodGet || !strings.HasPrefix(r.URL.Path, "/hosts/") {
			http.NotFound(w, r)
			return
		}

		hostID := strings.TrimPrefix(r.URL.Path, "/hosts/")
		if hostID == "" {
			http.Error(w, "missing host", http.StatusBadRequest)
			return
		}

		if hostID != "host-123" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sampleHost)
	}
}
