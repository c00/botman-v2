build: 
	go build -o bin/bot internal/cmd/*.go 

install: 
	go install internal/cmd/*.go 

run: 
	go run internal/cmd/*.go 

test:
	go test ./...