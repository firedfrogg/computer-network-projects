package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
	utilsSample "jarkom.cs.ui.ac.id/h01/samples/quic/utils"
)

// Generate pesan dari paket
func Handler(packet utils.LRTPIDSPacket) string {
	if packet.IsTrainArriving {
		// Kedatangan
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan tiba di Peron 1.", packet.Destination)
	} else if packet.IsTrainDeparting {
		// Keberangkatan
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan diberangkatkan dari Peron 1.", packet.Destination)
	}
	// Default case
	return ""
}

const (
	serverIP      = ""                         // IP address of the server dikosongkan untuk menerima paket dari sumber manapun
	serverPort    = "2696"                     // Port Unik
	serverType    = "udp4"                     // Tipe Protokol UDP4
	bufferSize    = 2048                       // Buffer size for reading
	appLayerProto = "lrt-jabodebek-2206082606" // Application layer protocol identifier
)

func main() {
	localUdpAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalf("Gagal untuk menyelesaikan alamat UDP: %v", err)
	}

	socket, err := net.ListenUDP(serverType, localUdpAddress)
	if err != nil {
		log.Fatalf("Tidak dapat membuka socket UDP: %v", err)
	}
	defer socket.Close()

	tlsConfig := &tls.Config{
		Certificates: utilsSample.GenerateTLSSelfSignedCertificates(),
		NextProtos:   []string{appLayerProto},
	}

	fmt.Printf("\nPIDS LRT Jabodebek Indonesia\n\n")
	fmt.Printf("[%s] Menyiapkan socket UDP di %s\n", serverType, socket.LocalAddr())

	listener, err := quic.Listen(socket, tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalf("Tidak dapat mendengarkan socket UDP: %v", err)
	}
	defer listener.Close()

	fmt.Printf("[quic] Menunggu koneksi QUIC di %s\n", listener.Addr())

	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Gagal menerima sesi: %v", err)
			continue
		}
		go handleSession(session)
	}
}

// Fungsi untuk menangani sesi QUIC
func handleSession(session quic.Connection) {
	for {
		// Menerima stream baru dalam sesi
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println(err)
			return
		}
		// Tangani setiap stream dalam goroutine terpisah
		go handleStream(session.RemoteAddr(), stream)
	}
}

func handleStream(clientAddress net.Addr, stream quic.Stream) {
	fmt.Printf("[quic] [Klien: %s] Membuka stream dengan ID %d\n", clientAddress, stream.StreamID())

	_, err := io.Copy(logicProcessorAndWriter{stream}, stream)
	if err != nil {
		fmt.Println(err)
	}
}

type logicProcessorAndWriter struct{ io.Writer }

func (lp logicProcessorAndWriter) Write(receivedMessageRaw []byte) (int, error) {
	// Decode pesan yang diterima
	receivedMessage := utils.Decoder(receivedMessageRaw)
	fmt.Printf("[quic] Pesan diterima: ")

	// Buat respons berdasarkan pesan
	response := Handler(receivedMessage)
	receivedMessage.IsAck = true

	// Encode dan kirim kembali respons ke klien
	writeLength, err := lp.Writer.Write(utils.Encoder(receivedMessage))
	fmt.Printf("\n%s\n", response)

	return writeLength, err
}
