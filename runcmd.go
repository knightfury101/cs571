package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func PerformCommmand_experiment(commandString string, file string) []byte {

	fmt.Println(commandString)
	cmdStr := fmt.Sprintf(commandString, file)
	fmt.Println(cmdStr)
	//cmd := exec.Command(myCommandInputsarray[0], myCommandInputsarray[1:]...)
	////////////
	cmd := exec.Command("bash", "-c", cmdStr)
	fmt.Println("MY COMMAND IS:")
	fmt.Println(cmd)
	////////////
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "::: " + stderr.String())
		panic(err)
	}

	return out.Bytes()
}

func PerformCommmand(command string, input string, fileIn string) []byte {

	cmd := exec.Command(command, input, fileIn)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	return out.Bytes()
}

func writeToFile(data []byte, filename string) {

	f, err := os.Create(filename)
	check(err)

	defer f.Close()

	f.Write(data)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func splitChunks(filename string, commandString string, idPos string, chunks int) []string {

	split := chunks

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	texts := make([]string, 0)
	fmt.Println("Reading in text")
	for scanner.Scan() {
		text := scanner.Text()
		texts = append(texts, text)
	}
	fmt.Println("Finished reading text")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var listofFileMade []string

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	fmt.Println("Splitting files")
	lengthPerSplit := len(texts) / split
	for i := 0; i < split; i++ {
		if i+1 == split {
			chunkTexts := texts[i*lengthPerSplit:]
			filename := writefile(strings.Join(chunkTexts, "\n"), commandString, idPos, r1.Intn(1000000))
			listofFileMade = append(listofFileMade, filename)
		} else {
			chunkTexts := texts[i*lengthPerSplit : (i+1)*lengthPerSplit]
			filename := writefile(strings.Join(chunkTexts, "\n"), commandString, idPos, r1.Intn(1000000))
			listofFileMade = append(listofFileMade, filename)
		}
	}
	fmt.Println("finished splitting files")
	texts = nil

	return listofFileMade
}

func splitChunksTimed(filename string, commandString string, idPos string, chunks int) []string {

	split := chunks

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	texts := make([]string, 0)

	//!!!!!!!!!!!!!!!!!!!!!!!!!
	readInTextStart := time.Now()

	fmt.Println("Reading in text")
	for scanner.Scan() {
		text := scanner.Text()
		texts = append(texts, text)
	}
	fmt.Println("Finished reading text")

	//!!!!!!!!!!!!!!!!!!!!
	duration := time.Since(readInTextStart)
	logTimenotes("Finished Reading text (SPLITCHUNKS) in "+filename+"=="+duration.String(), "logtime.txt")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var listofFileMade []string

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//!!!!!!!!!!!!!!!!!!!!!!!!!
	splitTime := time.Now()

	fmt.Println("Splitting files")
	lengthPerSplit := len(texts) / split
	for i := 0; i < split; i++ {
		if i+1 == split {
			chunkTexts := texts[i*lengthPerSplit:]
			filename := writefile(strings.Join(chunkTexts, "\n"), commandString, idPos, r1.Intn(1000000))
			listofFileMade = append(listofFileMade, filename)
		} else {
			chunkTexts := texts[i*lengthPerSplit : (i+1)*lengthPerSplit]
			filename := writefile(strings.Join(chunkTexts, "\n"), commandString, idPos, r1.Intn(1000000))
			listofFileMade = append(listofFileMade, filename)
		}
	}

	//!!!!!!!!!!!!!!!!!!!!
	duration = time.Since(splitTime)
	logTimenotes("Finished assigning text to split files (SPLITCHUNKS) in "+filename+"=="+duration.String(), "logtime.txt")

	fmt.Println("finished splitting files")
	texts = nil

	return listofFileMade
}

func writefile(data string, commandString string, idPos string, uniqueid int) string {

	file, err := os.Create("./filestore/chunks-" + commandString + "_" + idPos + "_" + strconv.Itoa(uniqueid) + ".txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file.WriteString(data)

	return file.Name()
}

func randomInputFileName(uniqueid int, command string) string {

	nameofFile := "input_" + command + "_" + strconv.Itoa(uniqueid) + ".txt"

	return nameofFile
}

func mergeFiles(files []string, commandString string) string {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	uniqueid := r1.Intn(1000000)

	mergedFileName := "./filestore/" + "mergeFile" + "_for_" + commandString + "_" + strconv.Itoa(uniqueid) + ".txt"

	f, err := os.OpenFile(mergedFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		buf.Write(b)
		f.Write(buf.Bytes())
		f.Write([]byte("\n"))
		buf.Reset()
	}
	return mergedFileName
}
