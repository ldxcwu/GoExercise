package main

import (
	"fmt"
	"secret"
)

func main() {
	vault := secret.NewVault("my-fake-key", ".secrets")
	vault.Set("gophercises", "<account:ldxcwu@163.com>, <password:test>")
	vault.Set("google", "<account:ldxcwu@163.com>, <password:lxw's google>")
	vault.Set("gophernotes", "<account:ldxcwu@163.com>, <password:lxw's gophernotes>")
	value, err := vault.Get("gophercises")
	if err != nil {
		return
	}
	fmt.Printf("value: %v\n", value)
	value, err = vault.Get("google")
	if err != nil {
		return
	}
	fmt.Printf("value: %v\n", value)
	value, err = vault.Get("gophernotes")
	if err != nil {
		return
	}
	fmt.Printf("value: %v\n", value)
}
