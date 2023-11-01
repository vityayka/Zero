package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	args := os.Args
	for i, arg := range args {
		fmt.Println(i, arg)
	}

	switch os.Args[1] {
	case "hash":
		fmt.Println(hash(os.Args[2]))
	case "compare":
		if compare(os.Args[2], os.Args[3]) {
			fmt.Println("Password correct")
		} else {
			fmt.Println("You fucked up")
		}
	}
}

func compare(s1 string, s2 string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s1), []byte(s2)) == nil
}

func hash(s string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)

}
