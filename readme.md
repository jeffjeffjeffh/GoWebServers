# Chirpy

Chirpy is a web server written in Go that supports CRUD operations for a Twitter(X?)-clone like service.

It uses a simple json file to store a database locally.

## Features

I completed this project is part of the [boot.dev](https://www.boot.dev/) curriculum.

The server uses the [chi](https://github.com/go-chi/chi) library to simplify routing.

After a user is created, it can be updated if the correct password is provided. Passwords are hashed before storage using [bcrypt](https://pkg.go.dev/golang.org/x/crypto@v0.16.0/bcrypt). When a user logs in, they receive a refresh token and an access token. The refresh endpoint accepts requests with a valid and unexpired refresh token to receive a new access token. Refresh tokens can be revoked on the server side.

After a user is logged in, they can post and delete chirps after they are authenticated and authorized.

There are endpoints for retrieving specific chirps, as well as a list of chirps, which supports finding chirps by specific users and sorting by ascending and descending order.

There are also endpoints for confirming that the service is running, receiving information about the number of API hits, and resetting the API hit count.
