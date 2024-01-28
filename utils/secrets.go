package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

func GetPathToEnv() string {
	var BinaryDir string

	if runtime.GOOS == "darwin" {
		BinaryDir = "/usr/local/bin"
	} else if runtime.GOOS == "linux" {
		BinaryDir = "/usr/bin"
	} else {
		BinaryDir = "./"
	}

	return BinaryDir
}

func GetSK() string {
	path := GetPathToEnv()
	fullPath := fmt.Sprintf("%s/.env.iago", path)
	godotenv.Load(fullPath)
	pk := os.Getenv("OPEN_API_SK")
	return pk
}
