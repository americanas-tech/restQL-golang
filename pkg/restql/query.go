package restql

import "github.com/b2wdigital/restQL-golang/v4/internal/domain"

type QueryContext struct {
	Mappings map[string]Mapping
	Options  domain.QueryOptions
	Input    domain.QueryInput
}

type SavedQuery struct {
	Text       string
	Deprecated bool
}
