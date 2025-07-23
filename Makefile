

build: 
	templ generate
	go fmt ./...
	go build

run:
	templ generate
	go fmt ./...
	go run . build examples/john-doe.yaml > john.html
	go run . build examples/jane-smith.yaml > jane.html
