dragass
=======

utility to generate slow/long HTTP requests

By trickling bunk headers. Slowly.

install
=======

	go get github.com/vbatts/dragass

usage
=====

	dragass -h
	Usage of dragass:
	  -c=1: concurrency of the requests
	  -delay=10: maximum delay between sending additional headers
	  -host="127.0.0.1": host address to connect to
	  -n=1: number of requests to make (total)
	  -port="8080": host port to connect to
	  -secs=30: number seconds to stretch request to

Default is a single, 30 second request to localhost:8080

For a series of long requests:

	dragass -c 4 -n 100 -secs 60

This would dispatch a total of 100 requests, spanned over 4 concurrent workers.
Each request lasting 60 seconds. 

output
======

'.' - A header sent over the wire
'+' - Response status was 200 OK
'-' - Response status was _NOT_ 200 OK

