// References: https://golang.org/

package main

import( 
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

const (
	DEFAULT_PORT = ":8081"
	BLOCKED_WEBSITES = "blocked.ini"
)

var (
	client = &http.Client{
		Transport: nil,
		CheckRedirect: nil,
		Jar: nil,
		Timeout: 0,
	}
	blockedStr string
)

func init(){
	//Reads in the blocked websites
	blocked, err := ioutil.ReadFile(BLOCKED_WEBSITES)
	if err != nil{
		fmt.Println("Error reading blocked.ini")
		log.Fatal(err)
	}
	blockedStr = string(blocked)
	fmt.Println("Blocked websites:\n", blockedStr)		//Prints content as string
}

func listenAndForward(r *http.Request) (resp *http.Response, bodyString string){
	//If host is on blocked list, don't query website, instead return blocked message
	if isBlocked(r.URL.Host){
		return nil, "<!DOCTYPE html><html><body><h1>This website has been blocked</h1></body></html>"
	}

	//Send the request to the destination
	resp, err := client.Do(r)
	if err != nil{
		fmt.Println("\nError getting a response\n")
		log.Fatal(err)
	}
	
	defer resp.Body.Close()
	bodybytes, er := ioutil.ReadAll(resp.Body)
	if er != nil{
		log.Fatal(er)
	}
	//bodyString contains the response body
	bodyString = string(bodybytes)	
	return resp, bodyString
}

//Checks if the host is blocked
func isBlocked(u string) bool{
	match, err := regexp.MatchString(u, blockedStr)
	if err != nil{
		log.Fatal(err)
	}
	if match{
		return true
	}
	return false
}

//Generates a http request from a byte slice of data heard on DEFAULT_PORT
func makeHeader(byteHeader []byte) *http.Request{
	//The below line reads the request but fills in fields for server to parse
	//As a result of this, if req is sent as a request to the host it does not return the correct response
	req, err := http.ReadRequest(bufio.NewReader(io.MultiReader(bytes.NewReader(byteHeader))))
	if err != nil{
		log.Fatal(err)
	}
	//Send request
	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err!= nil{
		log.Fatal(err)
	}
	return request
}

func handler(conn net.Conn){
	buff := make([]byte, 1024)
	msgLen, err := conn.Read(buff)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("\n\nRequest from client:\n",string(buff[:msgLen]))
	req := makeHeader(buff)
	resp, respBody := listenAndForward(req)
	if resp == nil{
		_, err = conn.Write([]byte(respBody))
		if err != nil{
			log.Fatal(err)
    }
  }
	fmt.Println("Response header received: ", resp)

	//Write back on the connection
	_, err = conn.Write([]byte(respBody))
	if err != nil{
		fmt.Println("Error writing data back to browser")
		log.Fatal(err)
	}
	defer conn.Close()
}


func main() {
	ln, err := net.Listen("tcp", DEFAULT_PORT)
	if err != nil{
		log.Fatal(err)
	}
	defer ln.Close()
	for{
		conn, err := ln.Accept()
		if err != nil{
			log.Fatal(err)
		}
		go handler(conn)
	}
}
