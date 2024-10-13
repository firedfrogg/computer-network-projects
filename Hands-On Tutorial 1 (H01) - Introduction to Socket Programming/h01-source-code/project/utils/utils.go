package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

type LRTPIDSPacketFixed struct {
	TransactionId     uint16 // 16-bit ID transaksi
	IsAck             bool   // 1-bit tanda ACK
	IsNewTrain        bool   // 1-bit tanda kereta baru
	IsUpdateTrain     bool   // 1-bit tanda update kereta
	IsDeleteTrain     bool   // 1-bit tanda hapus kereta
	IsTrainArriving   bool   // 1-bit tanda kereta tiba
	IsTrainDeparting  bool   // 1-bit tanda kereta berangkat
	TrainNumber       uint16 // 16-bit nomor kereta
	DestinationLength uint8  // 8-bit panjang destinasi
}

type LRTPIDSPacket struct {
	LRTPIDSPacketFixed
	Destination string // String destinasi dengan panjang DestinationLength
}

// Method untuk LRTPIDSPacket ke byte array
func Encoder(packet LRTPIDSPacket) []byte {
	var buf bytes.Buffer

	// Encode bagian tetap dari paket
	err := binary.Write(&buf, binary.BigEndian, packet.LRTPIDSPacketFixed)
	if err != nil {
		log.Fatalf("Gagal meng-encode paket: %v", err)
	}

	// Encode string destinasi
	err = binary.Write(&buf, binary.BigEndian, []byte(packet.Destination))
	if err != nil {
		log.Fatalf("Gagal meng-encode destinasi paket: %v", err)
	}

	return buf.Bytes()
}

// Method untuk mengubah byte array ke LRTPIDSPacket
func Decoder(data []byte) LRTPIDSPacket {
	packet := LRTPIDSPacket{}
	buf := bytes.NewReader(data)

	// Decode bagian tetap dari paket
	err := binary.Read(buf, binary.BigEndian, &packet.LRTPIDSPacketFixed)
	if err != nil {
		log.Fatalf("Gagal meng-decode paket: %v", err)
	}

	// Decode string destinasi sesuai panjang yang diberikan
	destination := make([]byte, packet.DestinationLength)
	err = binary.Read(buf, binary.BigEndian, &destination)
	if err != nil {
		log.Fatalf("Gagal meng-decode destinasi paket: %v", err)
	}

	packet.Destination = string(destination)
	return packet
}
