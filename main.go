package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

const (
	defaultIP   = "0.0.0.0"
	defaultPort = 9001
	defaultDir  = "./"
)

var (
	ip    *string
	port  *int
	pubIP string
)

func init() {
	ip = flag.String("ip", defaultIP, "bind ip")
	port = flag.Int("port", defaultPort, "bind port")
}

func getIPAddrs() (ipAddrs []string) {
	ift, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, ifi := range ift {
		if ifi.Flags&net.FlagUp == 0 {
			continue
		} else if ifi.Flags&net.FlagPointToPoint != 0 {
			continue
		}

		addrs, err := ifi.Addrs()
		if err != nil {
			println(err.Error())
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); !ok {
				continue
			} else if ipv4 := ipnet.IP.To4(); ipv4 == nil {
				continue
			} else {
				if !ipv4.IsLoopback() {
					pubIP = ipv4.String()
				}
				ipAddrs = append(ipAddrs, ipv4.String())
			}
		}
	}
	sort.Strings(ipAddrs)
	return
}

func printStatus(addrs []string, ip string, port int, dir string) {
	if ip != defaultIP {
		addrs = []string{ip}
	}

	y := color.New(color.FgHiYellow)
	g := color.New(color.FgHiGreen, color.Bold)
	w := color.New(color.FgHiWhite)

	y.Print("Starting up http-server, serving ")
	g.Println(dir)
	y.Println("Available on:")
	for _, addr := range addrs {
		w.Printf("\thttp://%s:", addr)
		g.Printf("%d\n", port)
	}
	w.Println("Hit CTRL-C to stop the server")
}

func saveReqFile(c *gin.Context) {
	r := c.Request

	filePath := r.URL.Path[1:]
	dirPath := path.Dir(filePath)
	ret := struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
		URI string `json:"uri,omitempty"`
	}{}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			ret.Msg = fmt.Sprintf("failed to create folder '%s'", dirPath)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ret)
			return
		}
	}

	filePath = fmt.Sprintf("%s.%d", filePath, time.Now().UnixNano()/1000)
	if c.Request.Header.Get("Content-Type") == "application/json" {
		filePath += ".json"
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		ret.Msg = fmt.Sprintf("failed to create file '%s'", filePath)
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ret)
		return
	}

	_, err = io.Copy(f, r.Body)
	if err != nil {
		ret.Msg = fmt.Sprintf("failed to write file '%s'", filePath)
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ret)
		return
	}
	f.Close()

	ret.Msg = "success"
	ret.URI = fmt.Sprintf("http://%s:%d/%s", pubIP, *port, filePath)
	c.JSON(http.StatusOK, ret)
	return
}

func main() {
	flag.Parse()
	host, dir, addrs := fmt.Sprintf("%s:%d", *ip, *port), defaultDir, getIPAddrs()

	{
		args := flag.Args()
		if len(args) != 0 {
			dir = args[0]
		}
	}
	printStatus(addrs, *ip, *port, dir)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		r := color.New(color.FgHiRed, color.Bold)
		w := color.New(color.FgHiWhite)
		t := time.Now().Format("02/Jan/2006 03:04:05")
		return w.Sprintf("[%s] - - [%s] %s \"%s\"\n",
			param.ClientIP, t, param.Request.UserAgent(),
			r.Sprintf("%s %s %s", param.Method, param.Path, param.Request.Proto),
		)

	}))
	r.Use(gin.Recovery())
	r.StaticFS("/", gin.Dir(dir, true))
	r.POST("/*path", saveReqFile)
	r.PUT("/*path", saveReqFile)
	r.Run(host)
}
