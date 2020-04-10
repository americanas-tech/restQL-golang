package plugins

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
)

type Runner struct {
	log *logger.Logger
}

func NewRunner(log *logger.Logger) *Runner {
	return &Runner{log: log}
}

func (r *Runner) BeforeQuery(plugin restql.Plugin, query string, queryCtx domain.QueryContext) {
	r.safeExecute(plugin.Name(), "BeforeQuery", func() {
		plugin.BeforeQuery(query, queryCtx)
	})
}

func (r *Runner) AfterQuery(plugin restql.Plugin, query string, result domain.Resources) {
	r.safeExecute(plugin.Name(), "AfterQuery", func() {
		m := convertQueryResult(result)
		plugin.AfterQuery(query, m)
	})
}

func (r *Runner) BeforeRequest(plugin restql.Plugin, request domain.HttpRequest) {
	r.safeExecute(plugin.Name(), "BeforeRequest", func() {
		plugin.BeforeRequest(request)
	})
}

func (r *Runner) AfterRequest(plugin restql.Plugin, request domain.HttpRequest, response domain.HttpResponse, err error) {
	r.safeExecute(plugin.Name(), "AfterRequest", func() {
		plugin.AfterRequest(request, response, err)
	})
}

func (r *Runner) safeExecute(pluginName string, hook string, fn func()) {
	go func() {
		defer func() {
			if reason := recover(); reason != nil {
				err := errors.Errorf("reason : %v", reason)
				r.log.Error("plugin produced a panic", err, "name", pluginName, "hook", hook)
			}
		}()

		fn()
	}()
}

func convertQueryResult(resource interface{}) map[string]interface{} {
	switch resource := resource.(type) {
	case domain.Resources:
		m := make(map[string]interface{})
		for k, v := range resource {
			m[string(k)] = convertDoneResource(v)
		}
		return m
	case domain.Details:
		return map[string]interface{}{
			"status":       resource.Status,
			"success":      resource.Success,
			"ignoreErrors": resource.IgnoreErrors,
			"debugging":    convertQueryResult(resource.Debug),
		}
	case *domain.Debugging:
		return map[string]interface{}{
			"method":          resource.Method,
			"url":             resource.Url,
			"requestHeaders":  resource.RequestHeaders,
			"responseHeaders": resource.ResponseHeaders,
			"params":          resource.Params,
			"requestBody":     resource.RequestBody,
			"responseTime":    resource.ResponseTime,
		}
	default:
		return nil
	}
}

func convertDoneResource(doneResource interface{}) interface{} {
	switch resource := doneResource.(type) {
	case domain.DoneResource:
		return map[string]interface{}{
			"details": convertQueryResult(resource.Details),
			"result":  resource.Result,
		}
	case domain.DoneResources:
		l := make([]interface{}, len(resource))
		for i, r := range resource {
			l[i] = convertQueryResult(r)
		}
		return l
	default:
		return resource
	}
}
