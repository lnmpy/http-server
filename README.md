# http-server

A static http-server in golang
Accept POST http request and save the request in local-file, works like the firebase.com does.

## Usage

```
$ http-server -h <ip> -p <port> <dir>

default:
ip: 0.0.0.0
port: 9001
dir: .
```

If POST a request with postman with url `http://localhost:9001/folder/test.json`, the server will create a folder named `folder` in current path and create a file named `test.json`, and store all the post request in it.

While do a GET request, it will run similar to the ftp-server, a simple static folder.

## Installation

```
$ go get github.com/lnmpy/http-server
```

or

```
$ git clone git@github.com:lnmpy/http-server.git
$ cd http-server && make install
```

## License

MIT

## Author

Elvis Macak (elvis@lnmpy.com)

## Thanks

* indexzero: base idea for http-server

    https://github.com/indexzero/http-server
