/*
1.Respond to HTTP & HTTPS requests,and should display eachrequest on a management console. It should forward the request to the Web server and relay the response to the browser.

2.Handle websocket connections.

3.Dynamically blockselected URLsvia the management console.

4.Efficiently cache requestslocally and thus save bandwidth. You must gather timing and bandwidth data to prove the efficiency of your proxy.

5.Handle multiple requests simultaneouslyby implementing a threaded server.
*/

// References: https://golang.org/

package main

import( 
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
)

const (
	defaultPort = ":8081"
)

var (
	client = &http.Client{
		Transport: nil,
		CheckRedirect: nil,
		Jar: nil,
		Timeout: 0,
	}
)

//Response from server we receive after issuing a GET request to a given URL
func genRequest(u url.URL) *http.Response{
	//Make a GET request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err!=nil{
		log.Fatal(err)
	}

	listenAndForward(req)
	
	resp, err := client.Do(req)
	if err != nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK{
		fmt.Print("Response from client.Do: \n\n", resp)
	}
	return resp
}

//TODO
//Given a http request, this function will forward the request to the host 
func listenAndForward(r *http.Request) {
	host := r.URL.Host	//Host to send to
	method := r.Method
	browserAddr := r.RemoteAddr
	fmt.Println("Host: ", host, "\nMethod: ", method, "\nRemoteAddr: ", browserAddr, "\n\n")
	ip, err := net.ResolveIPAddr("tcp", browserAddr)
	if err != nil{
		fmt.Println("ip: ", ip)
		log.Fatal(err)
	}
	//Send the http request to the destination
	resp, err := client.Do(r)
	if err != nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()
	responseToClient(resp, ip)
}

//Response returned from a request is given to the browser
func responseToClient(resp *http.Response, ip *net.IPAddr){
	
}

//Need to fix main so it listens to the raw bytes that are sent to the default port
//Using these raw bytes we will be able to infer the header (Browser will send the header after it is configured to send all data to specific port)
func main() {
	ln, err := net.Listen("tcp", defaultPort)
	if err != nil{
		log.Fatal(err)
	}
	defer ln.Close()

	tmp := make([]byte, 256)
	buf := make([]byte, 0, 4096)
	for{
		conn, err := ln.Accept()
		if err != nil{
			log.Fatal(err)
		}
		defer conn.Close()
		n, err := conn.Read(buf)
		if err != nil{
			log.Fatal(err)
		}
		buf = append(buf, tmp[0:n]...)
		fmt.Println("\nconn.RemoteAddr()\t", conn.RemoteAddr())
		fmt.Println("n: ", n, "\nbuf: ", buf, "\n&buf", &buf, "\ntmp", tmp, "\n&tmp: ", &tmp)

		er := http.Serve(ln, nil)
		fmt.Println("er: ", er)
	}
}
