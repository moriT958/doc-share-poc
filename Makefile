.PHONY: dev
dev: main.go
	go run main.go

.PHONY: serve
serve: dist/
	caddy file-server --root ./dist --listen :8080
