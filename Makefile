OWNER = lnmpy
REPO = http-server
TAG = `git describe --tags`


install:
	dep ensure

build: clean
	go build -o http-server

release: clean
	GOOS=darwin GOARCH=amd64 go build -o http-server && tar -czf http-server_darwin.tar.gz http-server
	GOOS=linux GOARCH=amd64 go build -o http-server && tar -czf http-server_linux.tar.gz http-server
	GOOS=windows GOARCH=amd64 go build -o http-server.exe  && tar -czf http-server_windows.tar.gz http-server.exe
	github-release release -u ${OWNER} -r ${REPO} -t ${TAG} -n ${TAG}
	github-release upload  -u ${OWNER} -r ${REPO} -t ${TAG} -n "http-server_darwin.tar.gz" -f http-server_darwin.tar.gz
	github-release upload  -u ${OWNER} -r ${REPO} -t ${TAG} -n "http-server_linux.tar.gz" -f http-server_linux.tar.gz
	github-release upload  -u ${OWNER} -r ${REPO} -t ${TAG} -n "http-server_windows.tar.gz" -f http-server_windows.tar.gz

clean:
	rm -f http-server*
