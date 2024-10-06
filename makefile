build: 
	go build -o bin/botman internal/cmd/botman/*.go 
	go build -o bin/botman-config internal/cmd/botmanconfig/*.go 

install: 
	go install internal/cmd/botman/*.go 
	go install internal/cmd/botmanconfig/*.go 

run: 
	go run internal/cmd/botman/*.go 

conf: 
	go run internal/cmd/botmanconfig/*.go 

test:
	go test ./...