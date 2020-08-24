package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"io/ioutil"
	"testing"

	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/test"
)

const defaultTenant = "default"

func TestMappingsReader_Env(t *testing.T) {
	envSource := stubEnvSource{
		getAll: map[string]string{
			"RESTQL_MAPPING_HERO":     "http://hero.api/",
			"RESTQL_MAPPING_SIDEKICK": "http://sidekick.api/",
			"RESTQL_MAPPING_":         "http://failed.api/",
			"TEST_VAR":                "foo",
		},
	}
	db := stubDatabase{}

	reader := NewMappingReader(noOpLogger, envSource, map[string]string{}, db)

	heroMapping, err := restql.NewMapping("hero", "http://hero.api/")
	test.VerifyError(t, err)

	sidekickMapping, err := restql.NewMapping("sidekick", "http://sidekick.api/")
	test.VerifyError(t, err)

	expected := map[string]restql.Mapping{
		"hero":     heroMapping,
		"sidekick": sidekickMapping,
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

func TestMappingsReader_Local(t *testing.T) {
	envSource := stubEnvSource{getAll: map[string]string{}}
	local := map[string]string{
		"hero":     "http://hero.api/",
		"sidekick": "http://sidekick.api/",
	}
	db := stubDatabase{}

	reader := NewMappingReader(noOpLogger, envSource, local, db)

	heroMapping, err := restql.NewMapping("hero", "http://hero.api/")
	test.VerifyError(t, err)

	sidekickMapping, err := restql.NewMapping("sidekick", "http://sidekick.api/")
	test.VerifyError(t, err)

	expected := map[string]restql.Mapping{
		"hero":     heroMapping,
		"sidekick": sidekickMapping,
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

func TestMappingsReader_Database(t *testing.T) {
	envSource := stubEnvSource{getAll: map[string]string{}}
	local := map[string]string{}

	heroMapping, err := restql.NewMapping("hero", "http://hero.api/")
	test.VerifyError(t, err)

	sidekickMapping, err := restql.NewMapping("sidekick", "http://sidekick.api/")
	test.VerifyError(t, err)

	db := stubDatabase{findMappingsForTenant: []restql.Mapping{heroMapping, sidekickMapping}}

	reader := NewMappingReader(noOpLogger, envSource, local, db)

	expected := map[string]restql.Mapping{
		"hero":     heroMapping,
		"sidekick": sidekickMapping,
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

func TestMappingsReader_ShouldOverwriteMappings(t *testing.T) {
	heroMapping, err := restql.NewMapping("hero", "https://hero.com/api/")
	test.VerifyError(t, err)

	sidekickMapping, err := restql.NewMapping("sidekick", "https://sidekick.com/api")
	test.VerifyError(t, err)

	villainMapping, err := restql.NewMapping("villain", "http://villain.api/")
	test.VerifyError(t, err)

	local := map[string]string{
		"hero":     "http://hero.api/",
		"sidekick": "http://sidekick.api/",
		"villain":  "http://villain.api/",
	}
	db := stubDatabase{
		findMappingsForTenant: []restql.Mapping{sidekickMapping},
	}
	envSource := stubEnvSource{
		getAll: map[string]string{
			"RESTQL_MAPPING_HERO": "https://hero.com/api/",
		},
	}

	reader := NewMappingReader(noOpLogger, envSource, local, db)

	expected := map[string]restql.Mapping{
		"hero":     heroMapping,
		"sidekick": sidekickMapping,
		"villain":  villainMapping,
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

var noOpLogger = logger.New(ioutil.Discard, logger.LogOptions{})

type stubDatabase struct {
	findMappingsForTenant []restql.Mapping
	findQuery             restql.SavedQuery
}

func (s stubDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]restql.Mapping, error) {
	return s.findMappingsForTenant, nil
}

func (s stubDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (restql.SavedQuery, error) {
	return s.findQuery, nil
}

type stubEnvSource struct {
	getAll map[string]string
}

func (s stubEnvSource) GetString(key string) string {
	return ""
}

func (s stubEnvSource) GetAll() map[string]string {
	return s.getAll
}
