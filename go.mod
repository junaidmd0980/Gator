module github.com/junaidmd0980/gator

go 1.27.2

require (
	config v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.12.3
)

replace config => ./internal/config
