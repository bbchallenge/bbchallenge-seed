main: lib_bbchallenge/*
	go build .
tests: lib_bbchallenge/*
	go test ./... -v
write_header: write_header.c
	gcc write_header.c -o write_header
