rediscache.so: $(shell find . -name '*.go')
	go build -buildmode=plugin -o rediscache.so rediscache.go