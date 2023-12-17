module github.com/jeffjeffjeffh/GoWebServers

go 1.21.3

require github.com/go-chi/chi/v5 v5.0.10

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1
	internal/testDatabase v1.0.0
)

require golang.org/x/crypto v0.16.0 // indirect

replace internal/testDatabase => ./internal/testDatabase
