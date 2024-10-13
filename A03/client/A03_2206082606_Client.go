package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 2048
)

type Student struct {
	Nama string
	Npm  string
}

type GreetResponse struct {
	Student Student
	Greeter string
}

type HttpRequest struct {
	Method  string
	Uri     string
	Version string
	Host    string
	Accept  string
}

type HttpResponse struct {
	Version       string
	StatusCode    string
	ContentType   string
	ContentLength int
	Data          string
}

func main() {
	var url, mimeType string

	// Input URL dari user
	fmt.Print("Input the url: ")
	fmt.Scan(&url)

	// Input MIME type dari user
	fmt.Print("Input the data type: ")
	fmt.Scan(&mimeType)

	// Parsing URL untuk mendapatkan host dan port
	parsedUrl := strings.Split(url, "/")
	hostPort := strings.Split(parsedUrl[2], ":")
	host := hostPort[0]
	port := hostPort[1]

	// Membuat HTTPRequest
	req := HttpRequest{
		Method:  "GET",
		Uri:     "/" + strings.Join(parsedUrl[3:], "/"),
		Version: "HTTP/1.1",
		Host:    host,
		Accept:  mimeType,
	}

	// Koneksi ke server
	connection, err := net.Dial(SERVER_TYPE, host+":"+port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer connection.Close()

	// Encode request dan kirim ke server
	requestMessage := RequestEncoder(req)
	_, err = connection.Write(requestMessage)
	if err != nil {
		fmt.Println("Error writing to server:", err)
		return
	}

	// Membaca response dari server
	buffer := make([]byte, BUFFER_SIZE)
	n, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	// Decode response dari server
	response := ResponseDecoder(buffer[:n])
	fmt.Printf("Status Code: %s\nBody: %s\n", response.StatusCode, response.Data)
}

func Fetch(req HttpRequest, connection net.Conn) HttpResponse {
	// Mengirim request ke server
	requestMessage := RequestEncoder(req)
	connection.Write(requestMessage)

	// Menerima response dari server
	buffer := make([]byte, BUFFER_SIZE)
	n, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
	}

	// Meng-decode response dan return hasilnya
	return ResponseDecoder(buffer[:n])
}

func RequestEncoder(req HttpRequest) []byte {
	// Encoding HTTP Request ke byte array
	requestLine := fmt.Sprintf("%s %s %s\r\n", req.Method, req.Uri, req.Version)
	headers := fmt.Sprintf("Host: %s\r\nAccept: %s\r\n\r\n", req.Host, req.Accept)
	return []byte(requestLine + headers)
}

func ResponseDecoder(bytestream []byte) HttpResponse {
	// Mengubah byte stream ke bentuk HttpResponse
	responseStr := string(bytestream)
	lines := strings.Split(responseStr, "\r\n")
	statusLine := strings.Split(lines[0], " ")
	statusCode := statusLine[1]

	body := strings.Join(lines[3:], "\n")

	return HttpResponse{
		Version:    statusLine[0],
		StatusCode: statusCode,
		Data:       body,
	}
}
