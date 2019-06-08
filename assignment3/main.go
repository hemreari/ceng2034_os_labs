package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	createFolders()

	log.Println("Starting without concurrency...")
	startPlain := time.Now()

	splitFiles()

	elapsedPlain := time.Since(startPlain)
	log.Println("Operation finished without concurrency in: ", elapsedPlain)

	err := os.RemoveAll("jpegFolder")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("mp3Folder")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("textFolder")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("pngFolder")
	if err != nil {
		log.Fatal(err)
	}

	createFolders()

	log.Println("Starting with concurrency...")
	start := time.Now()
	done := make(chan string)
	splitFilesConcurrency(done)
	<-done
	elapsed := time.Since(start)
	log.Println("Operation finished with concurrency in: ", elapsed)
}

func splitFiles() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "file") {

			content, err := ioutil.ReadFile(f.Name())
			if err != nil {
				log.Fatal(err)
			}

			magicJpeg := [4]int{255, 216, 255, 224}
			magicMp3 := [4]int{73, 68, 51, 3}
			magicText := [4]int{120, 121, 122, 0}
			magicPng := [4]int{137, 80, 78, 71}

			var sumJpeg int
			var sumMp3 int
			var sumText int
			var sumPng int
			var sumMagic int

			for _, value := range magicJpeg {
				sumJpeg = sumJpeg + value
			}

			for _, value := range magicMp3 {
				sumMp3 = sumMp3 + value
			}

			for _, value := range magicText {
				sumText = sumText + value
			}

			for _, value := range magicPng {
				sumPng = sumPng + value
			}

			magicBytes := content[:4]

			for _, value := range magicBytes {
				sumMagic = sumMagic + int(value)
			}

			switch sumMagic {
			case sumJpeg:
				err = copy(f.Name(), "jpegFolder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			case sumMp3:
				err = copy(f.Name(), "mp3Folder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			case sumText:
				err = copy(f.Name(), "textFolder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			case sumPng:
				err = copy(f.Name(), "pngFolder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func splitFilesConcurrency(c chan string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "file") {

			content, err := ioutil.ReadFile(f.Name())
			if err != nil {
				log.Fatal(err)
			}

			magicJpeg := [4]int{255, 216, 255, 224}
			magicMp3 := [4]int{73, 68, 51, 3}
			magicText := [4]int{120, 121, 122, 0}
			magicPng := [4]int{137, 80, 78, 71}

			var sumJpeg int
			var sumMp3 int
			var sumText int
			var sumPng int
			var sumMagic int

			go func() {
				for _, value := range magicJpeg {
					sumJpeg = sumJpeg + value
				}
				c <- "done"
			}()

			go func() {
				for _, value := range magicMp3 {
					sumMp3 = sumMp3 + value
				}
				c <- "done"
			}()

			go func() {
				for _, value := range magicText {
					sumText = sumText + value
				}
				c <- "done"
			}()

			go func() {
				for _, value := range magicPng {
					sumPng = sumPng + value
				}
				c <- "done"
			}()

			magicBytes := content[:4]

			for _, value := range magicBytes {
				sumMagic = sumMagic + int(value)
			}

			switch sumMagic {
			case sumJpeg:
				err = copy(f.Name(), "jpegFolder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			case sumMp3:
				err = copy(f.Name(), "mp3Folder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			case sumText:
				err = copy(f.Name(), "textFolder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			case sumPng:
				err = copy(f.Name(), "pngFolder/"+f.Name())
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func createFolders() {
	if _, err := os.Stat("jpegFolder"); os.IsNotExist(err) {
		os.Mkdir("jpegFolder", 0777)
	}

	if _, err := os.Stat("mp3Folder"); os.IsNotExist(err) {
		os.Mkdir("mp3Folder", 0777)
	}

	if _, err := os.Stat("pngFolder"); os.IsNotExist(err) {
		os.Mkdir("pngFolder", 0777)
	}

	if _, err := os.Stat("textFolder"); os.IsNotExist(err) {
		os.Mkdir("textFolder", 0777)
	}
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
