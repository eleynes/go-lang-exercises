module example.com/auth

go 1.22.0

replace example.com/db => ../db

require example.com/db v0.0.0-00010101000000-000000000000

require (
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.22.0
)
