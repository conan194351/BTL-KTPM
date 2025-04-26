dev:
	export $(shell cat .env) && \
		air -c .air.toml

local:
	export $(shell cat .env) && \
		go run cmd/server/server.go

worker:
	export $(shell cat .env) && \
		go run cmd/worker/worker.go