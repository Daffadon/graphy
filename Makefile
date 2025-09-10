gql-init:
	@printf '//go:build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go
	@go mod tidy
	@go run github.com/99designs/gqlgen init

gql-generate:
	@go run github.com/99designs/gqlgen generate