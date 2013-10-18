package main

import (
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/cactus/go-statsd-client/statsd"
	"os"
	"regexp"
)

var (
	path        string
	statsd_host string
	prefix      string
)

func init() {
	flag.StringVar(&path, "l", "", "Path to the MySQL slow query log")
	flag.StringVar(&statsd_host, "h", "localhost:8125", "Hostname:port for statsd")
	flag.StringVar(&prefix, "m", "mysql.queries.slow", "Path prefix for the metric to be recorded")
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "You need arguments")
		flag.Usage()
		os.Exit(10)
	}
	flag.Parse()

	// Seek to the end at the start
	seek := &tail.SeekInfo{
		Offset: 0,
		Whence: 2,
	}

	config := tail.Config{
		ReOpen:    true,
		MustExist: true,
		Follow:    true,
		Location:  seek,
	}

	t, err := tail.TailFile(path, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not tail: %v\n", err)
	}

	conn, err := statsd.New(statsd_host, prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to statsd: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close()

	r, _ := regexp.Compile(`^# Query_time:\s+(?P<query_time>\S+)\s+Lock_time:\s+(?P<lock_time>\S+).+`)

	for line := range t.Lines {
		if r.MatchString(line.Text) {
			_ = conn.Inc("count", 1, 1.0)
		}
	}
}
