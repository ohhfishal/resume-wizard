

build: 
	templ generate
	go fmt ./...
	go build

run:
	templ generate
	go fmt ./...
	go run . build examples/v2-john-doe.yaml > test.html
	cat test.html
		
