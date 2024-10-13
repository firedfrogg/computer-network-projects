package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP          = "35.247.92.49"
	serverPort        = "2696"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "lrt-jabodebek-2206082606"
	sslKeyLogFileName = "ssl-key.log"
)

func main() {
	// Membuat file log untuk SSL Key Log
	sslKeyLogFile, err := os.Create(sslKeyLogFileName)
	if err != nil {
		log.Fatalln("Gagal membuat file log SSL Key:", err)
	}
	defer sslKeyLogFile.Close()
	fmt.Printf("\nKendali Stasiun\n\n")

	// Konfigurasi TLS untuk QUIC dengan mengizinkan verifikasi sertifikat dilewati
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}

	// Menghubungkan ke server QUIC menggunakan alamat dan konfigurasi TLS yang telah ditentukan
	connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln("Gagal menghubungkan ke server QUIC:", err)
	}
	defer connection.CloseWithError(0x0, "Tidak ada kesalahan") // Menutup koneksi dengan pesan kesalahan "Tidak ada kesalahan"

	fmt.Printf("[quic] Menghubungkan dari %s ke %s\n", connection.LocalAddr(), connection.RemoteAddr())

	// Membuat buffer untuk menerima data
	fmt.Printf("[quic] Membuat buffer penerima dengan ukuran %d\n", bufferSize)
	receiveBuffer := make([]byte, bufferSize)

	// Membuka stream bidirectional pada koneksi QUIC
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln("Gagal membuka stream bidirectional:", err)
	}
	fmt.Printf("[quic] Membuka stream bidirectional %d ke %s\n", stream.StreamID(), connection.RemoteAddr())

	// Slice untuk menyimpan informasi paket
	packets := []utils.LRTPIDSPacket{
		{
			LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
				TransactionId:     0x01,                     // ID transaksi untuk paket pertama
				IsAck:             false,                    // Bukan paket acknowledge
				IsNewTrain:        false,                    // Bukan informasi kereta baru
				IsUpdateTrain:     false,                    // Bukan pembaruan informasi kereta
				IsDeleteTrain:     false,                    // Bukan penghapusan informasi kereta
				IsTrainArriving:   true,                     // Menandakan bahwa kereta tiba
				IsTrainDeparting:  false,                    // Tidak menandakan bahwa kereta berangkat
				TrainNumber:       42,                       // Nomor kereta
				DestinationLength: uint8(len("Harjamukti")), // Panjang nama tujuan
			},
			Destination: "Harjamukti", // Tujuan kereta
		},
		{
			LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
				TransactionId:     0x02,                     // ID transaksi untuk paket kedua
				IsAck:             false,                    // Bukan paket acknowledge
				IsNewTrain:        false,                    // Bukan informasi kereta baru
				IsUpdateTrain:     false,                    // Bukan pembaruan informasi kereta
				IsDeleteTrain:     false,                    // Bukan penghapusan informasi kereta
				IsTrainArriving:   false,                    // Tidak menandakan bahwa kereta tiba
				IsTrainDeparting:  true,                     // Menandakan bahwa kereta berangkat
				TrainNumber:       42,                       // Nomor kereta
				DestinationLength: uint8(len("Harjamukti")), // Panjang nama tujuan
			},
			Destination: "Harjamukti", // Tujuan kereta
		},
	}

	// Mengirim paket menggunakan loop
	for i, packet := range packets {
		// Mengkodekan paket ke format yang diinginkan
		packetEncoded := utils.Encoder(packet)

		// Mengirim paket melalui stream
		fmt.Printf("[quic] [Stream ID: %d] Mengirim paket %d\n", stream.StreamID(), i+1)
		_, err = stream.Write(packetEncoded)
		if err != nil {
			log.Fatalln("Gagal mengirim paket:", err)
		}
		fmt.Printf("[quic] [Stream ID: %d] Paket %d terkirim\n", stream.StreamID(), i+1)

		// Membaca respons dari server
		receiveLength, err := stream.Read(receiveBuffer)
		if err != nil {
			log.Fatalln("Gagal membaca respons dari server:", err)
		}
		fmt.Printf("[quic] [Stream ID: %d] Menerima %d byte pesan dari server untuk paket %d\n", stream.StreamID(), receiveLength, i+1)

		// Menguraikan pesan yang diterima dari server
		response := utils.Decoder(receiveBuffer[:receiveLength])
		fmt.Printf("[quic] [Stream ID: %d] Pesan diterima untuk paket %d:\n", stream.StreamID(), i+1)
		fmt.Print(response)
	}
}
