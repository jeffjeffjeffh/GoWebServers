module github.com/jeffjeffjeffh/GoWebServers

go 1.21.3

require github.com/go-chi/chi/v5 v5.0.10

require internal/testDatabase v1.0.0
replace internal/testDatabase => ./internal/testDatabase