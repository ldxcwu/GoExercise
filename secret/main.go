package main

import (
	"fmt"
	"secret/encrypt"
)

func main() {
	key := "ldxcwu"
	plain := "safari"
	cipher, _ := encrypt.Encrypt(key, plain)
	plaint, _ := encrypt.Decrypt(key, cipher)
	fmt.Printf("plaint: %v\n", plaint)
}
