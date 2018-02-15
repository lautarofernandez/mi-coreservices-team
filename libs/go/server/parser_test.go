package server_test

import (
	"testing"

	"github.com/mercadolibre/coreservices-team/libs/go/server"
)

func TestContextFromScopeString(t *testing.T) {
	tt := []struct {
		Scope         string
		ExpectedRole  server.Role
		ExpectedScope server.Environment
		ExpectedTag   string
		ExpectedErr   bool
	}{
		{"production-indexer", server.RoleIndexer, server.EnvProduction, "", false},
		{"sandbox-indexer", server.RoleIndexer, server.EnvSandbox, "", false},
		{"develop-indexer", server.RoleIndexer, server.EnvDevelop, "", false},
		{"integration-indexer", server.RoleIndexer, server.EnvIntegration, "", false},
		{"test-indexer", server.RoleIndexer, server.EnvTest, "", false},

		{"production-indexer-tag", server.RoleIndexer, server.EnvProduction, "tag", false},
		{"production-indexer-feature-new-search", server.RoleIndexer, server.EnvProduction, "feature-new-search", false},
		{"invalid-indexer", server.RoleIndexer, server.EnvProduction, "appname", true},
		{"test-invalid", server.RoleIndexer, server.EnvTest, "", true},
		{"invalid", server.RoleIndexer, server.EnvTest, "", true},
	}

	for _, tc := range tt {
		t.Run(tc.Scope, func(t *testing.T) {
			ctx, err := server.ContextFromScopeString(tc.Scope)

			if err != nil {
				if !tc.ExpectedErr {
					t.Fatalf("Unexpected error returned: %v", err)
				}

				return
			}

			if ctx.Role != tc.ExpectedRole {
				t.Fatalf("Expected role to be %v, got: %v", tc.ExpectedRole, ctx.Role)
			}

			if ctx.Environment != tc.ExpectedScope {
				t.Fatalf("Expected scope to be %v, got: %v", tc.ExpectedScope, ctx.Environment)
			}

			if ctx.Tag != tc.ExpectedTag {
				t.Fatalf(`Expected tag to be "%v", got: "%v"`, tc.ExpectedTag, ctx.Tag)
			}
		})
	}
}
