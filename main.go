package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"sort"

	"github.com/fatih/color"
)

const (
	defaultIP   = "0.0.0.0"
	defaultPort = 9001
	defaultDir  = "./"
)

var (
	ip   *string
	port *int
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

func handler(w http.ResponseWriter, r *http.Request) {

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "<b>Request Headers:</b></br>", r.URL.Path[1:])
	for _, k := range keys {
		fmt.Fprintln(w, k, ":", r.Header[k], "</br>", r.URL.Path[1:])
	}
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

	http.Handle("/", newFileHandler(http.FileServer(http.Dir(dir))))
	if err := http.ListenAndServe(host, nil); err != nil {
		color.HiRed(err.Error())
	}
}
