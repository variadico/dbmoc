.PHONY: test
test:
	go vet
	go test -cover
