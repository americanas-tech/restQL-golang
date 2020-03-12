package runner

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strings"
)

type Runner struct {
	config domain.Configuration
	log    domain.Logger
	client domain.HttpClient
}

func NewRunner(config domain.Configuration, httpClient domain.HttpClient, log domain.Logger) Runner {
	return Runner{
		config: config,
		log:    log,
		client: httpClient,
	}
}

func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, mappings map[string]domain.Mapping) interface{} {
	responses := make([]interface{}, len(query.Statements))

	for i, statement := range query.Statements {
		request := r.makeRequest(mappings, statement)

		response, err := r.client.Do(ctx, request)
		if err != nil {
			r.log.Debug("request failed", "error", err)
			return nil
		}

		responses[i] = response
	}

	return responses
}

func (r Runner) makeRequest(mappings map[string]domain.Mapping, statement domain.Statement) domain.Request {
	mapping := mappings[statement.Resource]
	url := makeUrl(mapping, statement)

	queryArgs := make(map[string]string)
	for key, value := range statement.With {
		if contains(mapping.PathParams, key) {
			continue
		}

		str, ok := value.(string)

		if !ok {
			r.log.Debug("skipping query param on request build for failing string casting", "param-key", key, "param-value", value)
			continue
		}
		queryArgs[key] = str
	}

	headers := make(map[string]string)
	for key, value := range statement.Headers {
		str, ok := value.(string)
		if !ok {
			r.log.Debug("skipping header on request build for failing string casting", "header-name", key, "header-value", value)
			continue
		}
		headers[key] = str
	}

	return domain.Request{
		Url:     url,
		Query:   queryArgs,
		Body:    nil,
		Headers: headers,
	}
}

func makeUrl(mapping domain.Mapping, statement domain.Statement) string {
	resource := mapping.Url
	for _, pathParam := range mapping.PathParams {
		resource = strings.Replace(resource, fmt.Sprintf(":%v", pathParam), fmt.Sprintf("%v", statement.With[pathParam]), 1)
	}

	return resource
}

func contains(list []string, item string) bool {
	for _, el := range list {
		if el == item {
			return true
		}
	}

	return false
}
