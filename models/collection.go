package models

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ Model = (*Collection)(nil)

const (
	CollectionTypeBase = "base"
	CollectionTypeAuth = "auth"
	CollectionTypeView = "view"
)

type Collection struct {
	BaseModel

	Name       string                  `db:"name" json:"name"`
	Type       string                  `db:"type" json:"type"`
	System     bool                    `db:"system" json:"system"`
	Schema     schema.Schema           `db:"schema" json:"schema"`
	Indexes    types.JsonArray[string] `db:"indexes" json:"indexes"`
	ListRule   *string                 `db:"listRule" json:"listRule"`
	ViewRule   *string                 `db:"viewRule" json:"viewRule"`
	CreateRule *string                 `db:"createRule" json:"createRule"`
	UpdateRule *string                 `db:"updateRule" json:"updateRule"`
	DeleteRule *string                 `db:"deleteRule" json:"deleteRule"`
	Options    types.JsonMap           `db:"options" json:"options"`
}

func (m *Collection) TableName() string {
	return "_collections"
}

func (m *Collection) IsBase() bool {
	return m.Type == CollectionTypeBase
}

func (m *Collection) IsAuth() bool {
	return m.Type == CollectionTypeAuth
}

func (m *Collection) IsView() bool {
	return m.Type == CollectionTypeView
}

// AuthOptions returns the collection options parsed as AuthCollectionOptions.
func (m *Collection) AuthOptions() AuthCollectionOptions {
	options := AuthCollectionOptions{}
	if len(m.Options) > 0 {
		_ = json.Unmarshal(m.Options, &options)
	}
	return options
}

// ViewOptions returns the collection options parsed as ViewCollectionOptions.
func (m *Collection) ViewOptions() ViewCollectionOptions {
	options := ViewCollectionOptions{}
	if len(m.Options) > 0 {
		_ = json.Unmarshal(m.Options, &options)
	}
	return options
}

type AuthCollectionOptions struct {
	ManageRule *string `json:"manageRule"`

	// auth options
	AllowOAuth2Auth     bool     `json:"allowOAuth2Auth"`
	AllowUsernameAuth   bool     `json:"allowUsernameAuth"`
	AllowEmailAuth      bool     `json:"allowEmailAuth"`
	RequireEmail        bool     `json:"requireEmail"`
	ExceptEmailDomains  []string `json:"exceptEmailDomains"`
	OnlyEmailDomains    []string `json:"onlyEmailDomains"`
	MinPasswordLength   int      `json:"minPasswordLength"`

	// oauth2 providers
	Providers []schema.AuthProviderConfig `json:"providers"`
}

type ViewCollectionOptions struct {
	Query string `json:"query"`
}

// MarshalJSON ad-hoc marshals the collection to hide the oauth2 client secrets.
func (m *Collection) MarshalJSON() ([]byte, error) {
	type alias Collection

	if m.IsAuth() {
		options := m.AuthOptions()
		for i, provider := range options.Providers {
			provider.ClientSecret = ""
			options.Providers[i] = provider
		}

		clone := *m
		clone.Options, _ = types.ParseJsonMap(options)
		return json.Marshal((*alias)(&clone))
	}

	return json.Marshal((*alias)(m))
}