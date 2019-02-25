// References: https://golang.org/

package main

import( 
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	fmt.Println(blockedStr)		//Prints content as string
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

func listenAndForward(r *http.Request) (resp *http.Response, headerBody string){
	//If host is on blocked list, don't query website
	if isBlocked(r.URL.Host){
		fmt.Println("This website has been blocked")
		os.Exit(1)
	}	
	fmt.Println("Request from client is : ", r)
	//Send the http request to the destination
	resp, err := client.Do(r)
	if err != nil{
		fmt.Println("\nError getting a response\n")
		log.Fatal(err)
	}
	//fmt.Println("\nResponse from host is: " , resp, "\n\nResponse body is: ", resp.Body)
	defer resp.Body.Close()
	bodybytes, er := ioutil.ReadAll(resp.Body)
	if er != nil{
		log.Fatal(er)
	}
	//bodyString contains the response body
	bodyString := string(bodybytes)
	return resp, bodyString
}

//Generates a http request from a byte slice of data heard on DEFAULT_PORT
func makeHeader(byteHeader []byte) *http.Request{
	req, err := http.ReadRequest(bufio.NewReader(io.MultiReader(bytes.NewReader(byteHeader))))
	if err != nil{
		log.Fatal(err)
	}
	//Client should not modify the RequestURI
	req.RequestURI = ""
	return req
}

func handleme(conn net.Conn){
	buff := make([]byte, 1024)
	msgLen, err := conn.Read(buff)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("\n\nRequest from browser: \n",string(buff[:msgLen]))
	req := makeHeader(buff)
	resp, respBody := listenAndForward(req)
	respStr := fmt.Sprintf("%v%s", resp, respBody)
	fmt.Println(respStr)
	//Write back on the connection
	_, err = conn.Write([]byte(respStr))
	if err != nil{
		fmt.Println("Error writing connection back to browser")
		log.Fatal(err)
	}
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
		go handleme(conn)
		defer conn.Close()
	}
}
