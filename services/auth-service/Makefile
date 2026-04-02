start:
	fuser -k 8000/tcp 2>/dev/null ; go run cmd/service/main.go

test:
	go test -count=1 -cover ./... -v
