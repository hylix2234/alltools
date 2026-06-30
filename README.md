# AllTools - Alat Pengujian Jaringan Berbasis Go

Kumpulan alat pengujian jaringan dan *stress-testing* berskala ringan namun berkinerja tinggi yang ditulis dalam bahasa Go. Repositori ini berisi berbagai skrip pengujian protokol jaringan untuk keperluan audit keamanan (*security auditing*), yang dikelola melalui satu menu interaktif terpusat.

## 📁 Struktur Repositori

* `main.go` / `menu.sh` — Antarmuka utama dan skrip bash interaktif untuk menjalankan setiap alat dengan mudah.
* `goackflood.go` (`goackflood`) — Implementasi Go untuk pengujian ACK.
* `gosynflood` (`gosynflood`) — Implementasi Go untuk pengujian SYN.
* `gosynackflood.go` (`gosynackflood`) — Implementasi Go untuk pengujian SYN-ACK.
* `go.mod` / `go.sum` — File manajemen dependensi modul Go.

---
![Screenshot](https://raw.githubusercontent.com/hylix2234/alltools/refs/heads/main/Screenshot_20260626_182247_com.termux.jpg)
## 🚀 Fitur Utama

* **Berbasis Go:** Menawarkan kemampuan konkurensi tinggi, penggunaan memori yang efisien, dan fitur *multi-threading* bawaan.
* **Menu Terpusat:** Dilengkapi dengan skrip interaktif (`menu.sh`) untuk menjalankan skrip dengan mudah tanpa perlu menghafal parameter perintah yang rumit.
* **Aplikasi Siap Pakai:** Menyertakan versi biner yang sudah dikompilasi untuk eksekusi instan.

---

## 🛠️ Pemasangan & Penggunaan

1. **Unduh Repositori** (Jika ingin digunakan di perangkat lain):
   ```bash
   git clone [https://github.com/hylix2234/alltools.git](https://github.com/hylix2234/alltools.git)
   cd alltools
2. **Masuk Ke Menu**
   ```bash
   ./menu.sh
   


