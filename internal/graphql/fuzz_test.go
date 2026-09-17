package graphql

import (
	"context"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
)

func FuzzGraphQLQueryParser(f *testing.F) {
	// Seed corpus with valid and invalid GraphQL patterns
	seeds := []string{
		`{ projects { id name capexCAD } }`,
		`{ project(id: "proj-01") { id name } }`,
		`query Introspection { __schema { sdl } }`,
		`{ __typename }`,
		`query { invalidField { nested { deep { recursion } } } }`,
		`query ($var: String!) { projects(sector: $var) }`,
		`{ projects(limit: -1, offset: 999999999999999999) { id } }`,
		`mutation { createProject(name: "Attack' OR '1'='1") }`,
		`{ """unclosed block string`,
		`{ { { { { { { { { { } } } } } } } } } }`,
		"\x00\x01\x02\xff\xfe\xfd",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	store := database.NewMemoryStore()
	handler := &Handler{resolver: NewResolver(store)}
	ctx := context.Background()

	f.Fuzz(func(t *testing.T, query string) {
		// Target 1: parseQuery directly - must never panic
		_, _, _, _ = parseQuery(query)

		// Target 2: execute through handler - must never panic or crash
		_, _ = handler.execute(ctx, gqlRequest{Query: query})
	})
}

func TestFuzzSmoke(t *testing.T) {
	store := database.NewMemoryStore()
	handler := &Handler{resolver: NewResolver(store)}
	ctx := context.Background()

	adversarialInputs := []string{
		"",
		" ",
		"\t\n\r",
		"{{{{{",
		"}}}}}",
		`{ "unterminated string`,
		`{ field(arg: "escape \n \t \" \x00") }`,
		`query { ` + string(make([]byte, 10000)),
		`{ projects(limit: -9999999999999999999999999999999999999) { id } }`,
		`{ projects(sector: "Critical Minerals'; DROP TABLE projects; --") { id } }`,
		`__schema`,
		`__typename`,
	}

	for _, input := range adversarialInputs {
		// Must not panic
		_, _, _, _ = parseQuery(input)
		_, _ = handler.execute(ctx, gqlRequest{Query: input})
	}
}
