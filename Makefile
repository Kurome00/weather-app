.PHONY: run build clean test run-console run-file run-json run-multi

run:
	go run ./cmd/linux/cli/main.go

run-console:
	go run ./cmd/linux/cli/main.go -logger=console -debug=true

run-file:
	go run ./cmd/linux/cli/main.go -logger=file -debug=true -logfile=app.log

run-json:
	go run ./cmd/linux/cli/main.go -logger=json -debug=true

run-multi:
	go run ./cmd/linux/cli/main.go -logger=multi -debug=true -logfile=app.log

build:
	go build -o build/weather-app ./cmd/linux/cli/main.go

clean:
	rm -rf build/
	rm -f *.log

test:
	go test -v ./...