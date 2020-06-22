build:fileserver

fileserver:*.go
	go build -ldflags "-s -w"
