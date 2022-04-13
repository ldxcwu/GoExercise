package main

import (
	"fmt"
	"secret"
)

func main() {
	vault := secret.NewVault("my-fake-key")
	vault.Set("gophercises", "<account:ldxcwu@163.com>, <password:test>")
	value, err := vault.Get("gophercises")
	if err != nil {
		return
	}
	fmt.Printf("value: %v\n", value)
}
