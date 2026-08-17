package models_test

import (
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestCollectionTableName(t *testing.T) {
	m := models.Collection{}
	if m.TableName() != "_collections" {
		t.Errorf("Expected _collections, got %s", m.TableName())
	}
}

func TestCollectionIsBase(t *testing.T) {
	m := models.Collection{Type: models.CollectionTypeBase}
	if !m.IsBase() {
		t.Errorf("Expected true, got false")
	}
}

func TestCollectionIsAuth(t *testing.T) {
	m := models.Collection{Type: models.CollectionTypeAuth}
	if !m.IsAuth() {
		t.Errorf("Expected true, got false")
	}
}

func TestCollectionIsView(t *testing.T) {
	m := models.Collection{Type: models.CollectionTypeView}
	if !m.IsView() {
		t.Errorf("Expected true, got false")
	}
}

func TestCollectionAuthOptions(t *testing.T) {
	m := models.Collection{}
	m.Options = types.JsonMap{"allowOAuth2Auth": true}
	options := m.AuthOptions()
	if !options.AllowOAuth2Auth {
		t.Errorf("Expected true, got false")
	}
}

func TestCollectionViewOptions(t *testing.T) {
	m := models.Collection{}
	m.Options = types.JsonMap{"query": "SELECT 1"}
	options := m.ViewOptions()
	if options.Query != "SELECT 1" {
		t.Errorf("Expected SELECT 1, got %s", options.Query)
	}
}

func TestCollectionMarshalJSON(t *testing.T) {
	m := models.Collection{
		Type: models.CollectionTypeAuth,
	}
	m.Options = types.JsonMap{
		"providers": []any{
			map[string]any{
				"name":         "google",
				"clientSecret": "secret123",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the secret is redacted in the JSON output
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	options, ok := result["options"].(map[string]any)
	if !ok {
		t.Fatalf("Expected options map, got %v", result["options"])
	}
	providers, ok := options["providers"].([]any)
	if !ok {
		t.Fatalf("Expected providers slice, got %v", options["providers"])
	}
	if len(providers) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(providers))
	}
	provider, ok := providers[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected provider map, got %v", providers[0])
	}
	if provider["clientSecret"] != "" {
		t.Errorf("Expected clientSecret to be empty, got %v", provider["clientSecret"])
	}

	// Check that the original collection's options are NOT mutated
	origOptions := m.AuthOptions()
	if len(origOptions.Providers) != 1 {
		t.Fatalf("Expected 1 original provider, got %d", len(origOptions.Providers))
	}
	if origOptions.Providers[0].ClientSecret != "secret123" {
		t.Errorf("Expected original clientSecret to remain 'secret123', got %q", origOptions.Providers[0].ClientSecret)
	}
}