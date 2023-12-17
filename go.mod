module github.com/jeffjeffjeffh/GoWebServers

go 1.21.3

require github.com/go-chi/chi/v5 v5.0.10

require internal/testDatabase v1.0.0

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.16.0 // indirect
)

replace internal/testDatabase => ./internal/testDatabase
