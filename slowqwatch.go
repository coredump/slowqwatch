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
	metric      string
	regex       string
)

func init() {
	flag.StringVar(&statsd_host, "h", "localhost:8125", "Hostname:port for statsd")
	flag.StringVar(&path, "l", "", "Path to log to be watched")
	flag.StringVar(&metric, "m", "mysql.queries.slow", "Metric to be increased")
	flag.StringVar(&regex, "r", "^# Query_time:\\s+(?P<query_time>\\S+)\\s+Lock_time:\\s+(?P<lock_time>\\S+).+", "Regex to be matched on the log")
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

	conn, err := statsd.New(statsd_host, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to statsd: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close()

	r := regexp.MustCompile(regex)

	for line := range t.Lines {
		if r.MatchString(line.Text) {
			_ = conn.Inc(metric, 1, 1.0)
		}
	}
}
