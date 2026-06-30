#!/bin/bash

# ==========================================
# WARNA ANSI & STYLE (BIAR KEREN)
# ==========================================
NC='\033[0m'
BOLD='\033[1m'
BLINK='\033[5m'

# Warna Terang
RED='\033[1;31m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
MAGENTA='\033[1;35m'
CYAN='\033[1;36m'
WHITE='\033[1;37m'
GRAY='\033[0;90m'

# Cek Akses Root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[!] Error: Script ini harus dijalankan dengan sudo / root!${NC}"
  exit 1
fi

# Fungsi untuk Menggambar Header
draw_header() {
  clear
  echo -e "${CYAN}┌────────────────────────────────────────────────────────┐${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${MAGENTA}  ███╗   ██╗███████╗████████╗██╗  ██╗██╗██╗  ██╗     ${NC} ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${MAGENTA}  ████╗  ██║██╔════╝╚══██╔══╝██║  ██║██║██║ ██╔╝     ${NC} ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${MAGENTA}  ██╔██╗ ██║█████╗     ██║   ███████║██║█████╔╝      ${NC} ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${MAGENTA}  ██║╚██╗██║██╔══╝     ██║   ██╔══██║██║██╔═██╗      ${NC} ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${MAGENTA}  ██║ ╚████║███████╗   ██║   ██║  ██║██║██║  ██╗     ${NC} ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${MAGENTA}  ╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝     ${NC} ${CYAN}│${NC}"
  echo -e "${CYAN}├────────────────────────────────────────────────────────┤${NC}"
  echo -e "${CYAN}│${NC} ${BOLD}${WHITE}  SYSTEM MENU HUB v17.5   ${NC} | ${YELLOW}OS: $(uname -s) ${NC} | ${GREEN}User: root${NC}    ${CYAN}│${NC}"
  echo -e "${CYAN}└────────────────────────────────────────────────────────┘${NC}"
}

# Fungsi Footer / Penutup boks
draw_footer() {
  echo -e "${CYAN}└────────────────────────────────────────────────────────┘${NC}"
}

# Loop Utama Menu
while true; do
  draw_header

  # List Menu dengan Box Drawing (Sudah disesuaikan dan dirapikan letaknya)
  echo -e "${CYAN}┌─────────────────── [ NETWORK TOOLS ] ──────────────────┐${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[1]${NC} Go SYN Flooder (Murni)                         ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[2]${NC} Go SYN Flooder (Mode Debug)                    ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[3]${NC} Fast Port Scanner (Nmap)                       ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[4]${NC} Check Website HTTP Status (Curl)               ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[5]${NC} Network Interface Monitor (IP Link)            ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[6]${NC} Go SYN-ACK Flooder                             ${CYAN}│${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${YELLOW}[7]${NC} Go ACK Flooder (Upgraded Vermin)               ${CYAN}│${NC}"
  echo -e "${CYAN}├────────────────────────────────────────────────────────┤${NC}"
  echo -e "${CYAN}│${NC}  ${BOLD}${RED}[0]${NC} Keluar / Exit                                  ${CYAN}│${NC}"
  draw_footer

  echo -e -n "\n${BOLD}${WHITE}Nethik-Menu>${NC} Pilih opsi [0-7]: "
  read -r pilihan

  case $pilihan in
    1)
      echo -e "\n${BOLD}${GREEN}[+] Menjalankan Go SYN Flooder Standard${NC}"
      echo -n "Masukkan IP Target: " && read -r ip
      echo -n "Masukkan Port Target (default 80): " && read -r port
      port=${port:-80}
      echo -n "Masukkan Interface (default eth0): " && read -r iface
      iface=${iface:-eth0}

      if [ -f "./gosynflood" ]; then
        ./gosynflood -t "$ip" -p "$port" -i "$iface"
      else
        echo -e "${RED}[!] File binary './gosynflood' tidak ditemukan! Silakan compile dulu.${NC}"
        read -p "Tekan Enter untuk kembali..."
      fi
      ;;

    2)
      echo -e "\n${BOLD}${RED}[+] Menjalankan Go SYN Flooder MODE DEBUG${NC}"
      echo -n "Masukkan IP Target: " && read -r ip
      echo -n "Masukkan Port Target (default 80): " && read -r port
      port=${port:-80}
      echo -n "Masukkan Interface (default eth0): " && read -r iface
      iface=${iface:-eth0}

      if [ -f "./gosynflood" ]; then
        ./gosynflood -t "$ip" -p "$port" -i "$iface" -d
      else
        echo -e "${RED}[!] File binary './gosynflood' tidak ditemukan! Silakan compile dulu.${NC}"
        read -p "Tekan Enter untuk kembali..."
      fi
      ;;

    3)
      echo -e "\n${BOLD}${CYAN}[+] Fast Port Scanner (Nmap)${NC}"
      echo -n "Masukkan Target IP/Domain: " && read -r scan_target
      if command -v nmap &> /dev/null; then
        nmap -F "$scan_target"
      else
        echo -e "${YELLOW}[*] Nmap belum terinstall. Menggunakan scanning internal netcat...${NC}"
        nc -zv -w 2 "$scan_target" 20-90
      fi
      echo "" && read -r -p "Tekan Enter untuk kembali ke menu..."
      ;;

    4)
      echo -e "\n${BOLD}${CYAN}[+] Check Website HTTP Status${NC}"
      echo -n "Masukkan URL (contoh: google.com): " && read -r url
      echo -e "${GRAY}[*] Menghubungi target...${NC}"
      curl -Is "$url" | head -n 10
      echo "" && read -r -p "Tekan Enter untuk kembali ke menu..."
      ;;

    5)
      echo -e "\n${BOLD}${CYAN}[+] Daftar Interface Jaringan Aktif:${NC}"
      ip -br link
      echo "" && read -r -p "Tekan Enter untuk kembali ke menu..."
      ;;

    6)                                                                                                                                              
      echo -e "\n${BOLD}${GREEN}[+] Menjalankan Go SYN-ACK Flooder Standard${NC}"
      echo -n "Masukkan IP Target: " && read -r ip
      echo -n "Masukkan Port Target (default 80): " && read -r port
      port=${port:-80}
      echo -n "Masukkan Interface (default eth0): " && read -r iface                                                                                
      iface=${iface:-eth0}

      if [ -f "./gosynackflood" ]; then
        ./gosynackflood -t "$ip" -p "$port" -i "$iface"
      else
        echo -e "${RED}[!] File binary './gosynackflood' tidak ditemukan! Silakan compile dulu.${NC}"                                                    
        read -p "Tekan Enter untuk kembali..."
      fi
      ;;

    7)
      echo -e "\n${BOLD}${GREEN}[+] Menjalankan Go ACK Flooder (Upgraded Vermin)${NC}"
      echo -n "Masukkan IP Target: " && read -r ip
      echo -n "Masukkan Port Target (default 80): " && read -r port
      port=${port:-80}
      echo -n "Jumlah Paket yang Dikirim (contoh: 5000): " && read -r counter

      if [ -f "./goackflood" ]; then
        # Memanggil binary Go dengan argument yang sudah diinput user di Bash
        ./goackflood -t "$ip" -p "$port" -c "$counter"
        echo "" && read -r -p "Tekan Enter untuk kembali ke menu..."
      else
        echo -e "${RED}[!] File binary './goackflood' tidak ditemukan! Silakan compile dulu.${NC}"
        read -p "Tekan Enter untuk kembali..."
      fi
      ;;

    0)
      echo -e "\n${BOLD}${BLUE}[*] Terima kasih telah menggunakan Nethik Menu Hub. Sampai jumpa!${NC}\n"
      exit 0
      ;;

    *)
      echo -e "\n${RED}[!] Pilihan tidak valid!${NC}"
      sleep 1
      ;;
  esac
done

