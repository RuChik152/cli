package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type MD5 string

func main() {
	createHashSummFile(".")
}

func createHashSummFile(path string) {
	if err := os.Remove("./hash_sum.log"); err != nil {
		println(err.Error())
	}
	appendFileData("./hash_sum.log", fmt.Sprintf("%-80s %-80s %s\n", "PATH", "NAME", "HASH"))
	filepath.Walk(path, func(wPath string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if wPath != path {
			data := fmt.Sprintf("%-80s %-80s %s\n", wPath, info.Name(), hashSum(wPath))
			appendFileData("./hash_sum.log", data)
		}
		return nil
	})
}

func hashSum(path string) MD5 {
	file, err := os.Open(path)
	if err != nil {
		log.Println("ошибка открытия файла: ", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	hashByteSum := hash.Sum(nil)
	md5String := fmt.Sprintf("%x", hashByteSum)
	return MD5(md5String)
}

func appendFileData(path string, data string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		panic(err)
	}
}
