package utils

import "fmt"

const (
	blue    = "\033[1;34m"
	green   = "\033[1;32m"
	yellow  = "\033[1;33m"
	red     = "\033[1;31m"
	reset   = "\033[0m"
)

func LogInfo(msg string) {
	fmt.Println(blue + "[*] " + msg + reset)
}

func LogSuccess(msg string) {
	fmt.Println(green + "[+] " + msg + reset)
}

func LogWarn(msg string) {
	fmt.Println(yellow + "[!] " + msg + reset)
}

func LogError(msg string) {
	fmt.Println(red + "[✘] " + msg + reset)
}
