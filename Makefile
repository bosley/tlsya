all:
	go build -o tlsya-cli cmd/cli/main.go

clean:
	rm tlsya-cli

test:
	go mod tidy && go test -v .
