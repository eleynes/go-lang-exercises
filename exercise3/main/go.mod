module example.com/exercise3

go 1.22.0

replace example.com/passwordgenerator => ../passwordgenerator

replace example.com/auth => ../auth

replace example.com/db => ../db

require example.com/passwordgenerator v0.0.0-00010101000000-000000000000

require (
	example.com/auth v0.0.0-00010101000000-000000000000
	example.com/db v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.22.0
)

require github.com/lib/pq v1.10.9 // indirect
