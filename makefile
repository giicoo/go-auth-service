run:
	go build ./cmd/main.go
	./main
vet:
	go vet ./...