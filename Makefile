main: lib_bbchallenge/*
	go build .
tests: lib_bbchallenge/*
	go test ./... -v