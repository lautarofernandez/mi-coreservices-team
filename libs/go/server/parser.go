package server

import (
	"fmt"
	"strings"
)

// ContextFromScopeString parses the fury scope to get the execution context (environment)
// Scope format must be: {environment}-{app role}[-{app name}]
// For example: test-search, develop-indexer, production-indexer-feature-new-context
func ContextFromScopeString(scope string) (ApplicationContext, error) {
	parts := strings.Split(strings.ToLower(scope), "-")

	// If we receive a scope with only 1 part, then we lack information for bootstraping the server.
	if len(parts) <= 1 {
		return ApplicationContext{}, fmt.Errorf("invalid scope received: %v", scope)
	}

	env, role := Environment(parts[0]), Role(parts[1])

	// Validate Role
	if role != RoleIndexer && role != RoleRead && role != RoleSearch && role != RoleWrite {
		return ApplicationContext{}, fmt.Errorf("invalid role inferred from scope: %v", role)
	}

	// Validate Scope
	if env != EnvProduction && env != EnvDevelop && env != EnvIntegration && env != EnvTest {
		return ApplicationContext{}, fmt.Errorf("invalid environment inferred from scope %v", role)
	}

	// If fury scope has a 3rd part, then we use that as some kind of tag for the application
	// Eg.:  We might use this tag for running a specific branch from the git repository.
	var tag string
	if len(parts) >= 3 {
		tag = strings.Join(parts[2:], "-")
	}

	return ApplicationContext{
		Environment: env,
		Role:        role,
		Tag:         tag,
	}, nil
}
