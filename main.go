package main

import (
	"cli_hash/gui"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type MD5 string

const HASH_LOG string = "hash_sum.log"

var path_log *string
var path_source *string
var user_gui *bool
var len_path_source int = 0

func init() {
	fmt.Println("Initialization...")
	initArgs()
}

func main() {

	if *user_gui {
		gui.StartGui(path_log, path_source)
		println(*path_log)
		println(*path_source)
	}

	println(*path_log)
	println(*path_source)
	startScan()

}

func createHashSummFile(path string) {
	if err := os.Remove(*path_log); err != nil {
		println(err.Error())
	}
	appendFileData(*path_log, fmt.Sprintf("%-150s %-50s %s\n", "PATH", "NAME", "HASH"))
	//appendFileData(*path_log, fmt.Sprintf("%-50s %-50s\n", "NAME", "HASH"))

	var totalFilesCount int

	filepath.Walk(path, func(wpath string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			re := regexp.MustCompile(`^\.+`)
			match := re.MatchString(info.Name())

			reLib := regexp.MustCompile(`^(L|l)ibrary`)
			matchLib := reLib.MatchString(info.Name())
			if match || matchLib {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			totalFilesCount++
		}

		return nil
	})

	bar := pb.StartNew(totalFilesCount)

	filepath.Walk(path, func(wPath string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			re := regexp.MustCompile(`^\.+`)
			match := re.MatchString(info.Name())

			reLib := regexp.MustCompile(`^(L|l)ibrary`)
			matchLib := reLib.MatchString(info.Name())
			if match || matchLib {
				return filepath.SkipDir
			}
			return nil
		}

		if wPath != path {
			appendFileData(*path_log, fmt.Sprintf("%-150s %-50s %s\n", "."+wPath[len_path_source:], info.Name(), hashSum(wPath)))
			//appendFileData(*path_log, fmt.Sprintf("%-50s %-50s\n", info.Name(), hashSum(wPath)))
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

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	} else {
		len_path_source = len(cwd)
	}

	path_log = flag.String("log", filepath.Join(cwd, HASH_LOG), "path for logs")
	path_source = flag.String("s", filepath.Join(cwd), "path to the directory to scan")
	user_gui = flag.Bool("g", false, "On or Off GUI")
	flag.Parse()
}

func startScan() {
	start := time.Now()
	createHashSummFile(*path_source)
	elapsed := time.Since(start)

	fmt.Printf("Time spent: %s \n", elapsed)
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
