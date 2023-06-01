package middleware

import (
	"testing"

	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/valyala/fasthttp"
)

func Test_tenantByHost_setTenant(t *testing.T) {
	type fields struct {
		log            restql.Logger
		tenantsByHosts map[string]string
		defaultTenant  string
	}
	tests := []struct {
		name       string
		fields     fields
		host       string
		tenant     string
		wantTenant string
	}{
		{
			name: "should set tenant",
			fields: fields{
				log: noOpLogger{},
				tenantsByHosts: map[string]string{
					"americanas.test": "acom-npf",
				},
				defaultTenant: "default-tenant",
			},
			host:       "americanas.test",
			wantTenant: "acom-npf",
		},
		{
			name: "should not set tenant if already set",
			fields: fields{
				log: noOpLogger{},
				tenantsByHosts: map[string]string{
					"americanas.test": "acom-npf",
				},
				defaultTenant: "default-tenant",
			},
			tenant:     "previous-set-tenant",
			host:       "americanas.test",
			wantTenant: "previous-set-tenant",
		},
		{
			name: "should set default tenant",
			fields: fields{
				log: noOpLogger{},
				tenantsByHosts: map[string]string{
					"americanas.test": "acom-npf",
				},
				defaultTenant: "default-tenant",
			},
			host:       "unknown.test",
			wantTenant: "default-tenant",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tenantByHost{
				log:           tt.fields.log,
				tenantsByHost: tt.fields.tenantsByHosts,
				defaultTenant: tt.fields.defaultTenant,
			}
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetHost(tt.host)
			ctx.QueryArgs().Add("tenant", tt.tenant)
			u.setTenant(ctx)
			if got := string(ctx.QueryArgs().Peek("tenant")); got != tt.wantTenant {
				t.Errorf("tenantByHost.setTenant() = %v, want %v", got, tt.wantTenant)
			}
		})
	}
}

type noOpLogger struct{}

func (n noOpLogger) Panic(msg string, fields ...interface{})            {}
func (n noOpLogger) Fatal(msg string, fields ...interface{})            {}
func (n noOpLogger) Error(msg string, err error, fields ...interface{}) {}
func (n noOpLogger) Warn(msg string, fields ...interface{})             {}
func (n noOpLogger) Info(msg string, fields ...interface{})             {}
func (n noOpLogger) Debug(msg string, fields ...interface{})            {}
func (n noOpLogger) With(key string, value interface{}) restql.Logger   { return n }
