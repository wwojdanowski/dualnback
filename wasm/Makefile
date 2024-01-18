build:
	GOARCH=wasm GOOS=js go build -o main.wasm wasm.go
	mv main.wasm static/

.PHONY: build
