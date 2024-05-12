package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

type MD5 string

const HASH_LOG string = "hash_sum.log"

var path_log *string
var path_source *string

func init() {
	initArgs()
}

func main() {
	createHashSummFile(*path_source)

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func createHashSummFile(path string) {
	if err := os.Remove(*path_log); err != nil {
		println(err.Error())
	}
	appendFileData(*path_log, fmt.Sprintf("%-150s %-50s %s\n", "PATH", "NAME", "HASH"))

	var totalFilesCount int

	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			totalFilesCount++
		}

		return nil
	})

	bar := pb.StartNew(totalFilesCount)

	filepath.Walk(path, func(wPath string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if wPath != path {
			data := fmt.Sprintf("%-150s %-50s %s\n", wPath, info.Name(), hashSum(wPath))
			appendFileData(*path_log, data)
		}
		bar.Increment()
		return nil
	})
	bar.Finish()

}

func hashSum(path string) MD5 {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error open file or dir: ", err)
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

func initArgs() {
	path_log = flag.String("log", filepath.Join(".", HASH_LOG), "path for logs")
	path_source = flag.String("s", filepath.Join("."), "path to the directory to scan")
	flag.Parse()
}
