package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"strings"
)

const (
	SERVER_HOST  = ""
	SERVER_PORT  = "2696"
	SERVER_TYPE  = "tcp"
	BUFFER_SIZE  = 2048
	STUDENT_NAME = "Julian"
	STUDENT_NPM  = "2206082606"
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
	// Memulai server dan mendengarkan koneksi yang masuk
	listener, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server started on %s:%s\n", SERVER_HOST, SERVER_PORT)

	for {
		// Menerima koneksi baru
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		go HandleConnection(connection)
	}
}

func HandleConnection(connection net.Conn) {
	defer connection.Close()

	// Membaca request dari client
	buffer := make([]byte, BUFFER_SIZE)
	n, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	// Decode request
	request := RequestDecoder(buffer[:n])

	// Tangani request dan kirimkan response ke client
	response := HandleRequest(request)
	connection.Write(ResponseEncoder(response))
}

func HandleRequest(req HttpRequest) HttpResponse {
    var response HttpResponse

    if req.Uri == "/" {
        // Handling '/' URI remains the same
        response = HttpResponse{
            Version:     "HTTP/1.1",
            StatusCode:  "200",
            ContentType: "text/html",
            Data:        fmt.Sprintf("<html><body><h1>Halo, dunia! Aku %s</h1></body></html>", STUDENT_NAME),
        }
    } else if strings.HasPrefix(req.Uri, "/greet/") {
        // Extract NPM from URI
        uriParts := strings.Split(req.Uri, "/")
        if len(uriParts) > 2 && uriParts[2] == STUDENT_NPM {
            // Valid NPM, handle greeting
            greeter := STUDENT_NAME
            // Check for name parameter
            if strings.Contains(req.Uri, "?name=") {
                nameParts := strings.Split(req.Uri, "?name=")
                if len(nameParts) > 1 {
                    greeter = nameParts[1]
                }
            }
            student := Student{Nama: STUDENT_NAME, Npm: STUDENT_NPM}

            // Check Accept header to determine content type
            if req.Accept == "application/xml" {
                // Respond in XML format
                xmlData, _ := xml.MarshalIndent(GreetResponse{Student: student, Greeter: greeter}, "", "  ")
                xmlString := string(xmlData)
                response = HttpResponse{
                    Version:     "HTTP/1.1",
                    StatusCode:  "200",
                    ContentType: "application/xml",
                    Data:        xmlString,
                }

                // Parse XML and add to HttpResponse
                var parsedResponse GreetResponse
                err := xml.Unmarshal([]byte(xmlString), &parsedResponse)
                if err != nil {
                    fmt.Println("Error parsing XML:", err)
                } else {
                    response.Data = fmt.Sprintf("Body: %s\nParsed: {{%s %s} %s}\n", xmlString, parsedResponse.Student.Nama, parsedResponse.Student.Npm, parsedResponse.Greeter)
                }

            } else if req.Accept == "application/json" || req.Accept == "" {
                // Respond in JSON format (or default if Accept is empty or unsupported)
                jsonData, _ := json.Marshal(GreetResponse{Student: student, Greeter: greeter})
                jsonString := string(jsonData)
                response = HttpResponse{
                    Version:     "HTTP/1.1",
                    StatusCode:  "200",
                    ContentType: "application/json",
                    Data:        jsonString,
                }

                // Parse JSON and add to HttpResponse
                var parsedResponse GreetResponse
                err := json.Unmarshal([]byte(jsonString), &parsedResponse)
                if err != nil {
                    fmt.Println("Error parsing JSON:", err)
                } else {
                    response.Data = fmt.Sprintf("Body: %s\nParsed: {{%s %s} %s}\n", jsonString, parsedResponse.Student.Nama, parsedResponse.Student.Npm, parsedResponse.Greeter)
                }

            } else {
                // Unsupported content type, default to JSON response
                jsonData, _ := json.Marshal(GreetResponse{Student: student, Greeter: greeter})
                jsonString := string(jsonData)
                response = HttpResponse{
                    Version:     "HTTP/1.1",
                    StatusCode:  "200",
                    ContentType: "application/json",
                    Data:        jsonString,
                }

                // Parse JSON and add to HttpResponse
                var parsedResponse GreetResponse
                err := json.Unmarshal([]byte(jsonString), &parsedResponse)
                if err != nil {
                    fmt.Println("Error parsing JSON:", err)
                } else {
                    response.Data = fmt.Sprintf("Body: %s\nParsed: {{%s %s} %s}\n", jsonString, parsedResponse.Student.Nama, parsedResponse.Student.Npm, parsedResponse.Greeter)
                }
            }
        } else {
            // NPM doesn't match, respond with 404
            response = HttpResponse{
                Version:    "HTTP/1.1",
                StatusCode: "404",
                Data:       "",
            }
        }
    } else {
        // Handle URI not found (404)
        response = HttpResponse{
            Version:    "HTTP/1.1",
            StatusCode: "404",
            Data:       "",
        }
    }

    return response
}


func RequestDecoder(bytestream []byte) HttpRequest {
	// Meng-decode byte stream menjadi HttpRequest
	requestStr := string(bytestream)
	lines := strings.Split(requestStr, "\r\n")
	requestLine := strings.Split(lines[0], " ")
	method := requestLine[0]
	uri := requestLine[1]
	version := requestLine[2]

	return HttpRequest{
		Method:  method,
		Uri:     uri,
		Version: version,
		Host:    strings.Split(lines[1], " ")[1],
		Accept:  strings.Split(lines[2], " ")[1],
	}
}

func ResponseEncoder(res HttpResponse) []byte {
	// Meng-encode HttpResponse menjadi byte array
	statusLine := fmt.Sprintf("%s %s OK\r\n", res.Version, res.StatusCode)
	headers := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n\r\n", res.ContentType, len(res.Data))
	return []byte(statusLine + headers + res.Data)
}
