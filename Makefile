dev:
	export $(shell cat .env) && \
		air -c .air.toml

local:
	export $(shell cat .env) && \
		go run cmd/server.go