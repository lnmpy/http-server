install:
	dep ensure

release:
	GOOS=darwin GOARCH=amd64 go build -o http-server-v${TAG}_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o http-server-v${TAG}_linux_amd64
	GOOS=windows GOARCH=amd64 go build -o http-server-v${TAG}_windows_amd64.exe

clean:
	rm -f http-server*
