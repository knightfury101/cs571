package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/goombaio/dag"
)

func main() {

	start()

}

func split_file_experiment(filename string, numChunks int) {

	file, err := os.Open(filename)

	stats, err := file.Stat()

	fileSize := stats.Size()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	numlines := lineCounter(filename)
	sizeofChunks := int64(numlines) / int64(numChunks)

	fmt.Println("sizeofFile")
	fmt.Println(fileSize)
	fmt.Println("Size of chunks")
	fmt.Println(sizeofChunks)
	fmt.Println("-----")

	currentFile := 0
	numLines := 0
	listOfText := []string{}
	totalChunks := 0

	for scanner.Scan() {
		linetext := scanner.Text()
		numLines += 1
		listOfText = append(listOfText, linetext)

		if numLines >= int(sizeofChunks) && totalChunks != (int(sizeofChunks)-1) {
			fileName := "./filestore/somebigfile_" + strconv.FormatInt(int64(currentFile), 10)
			filePtr, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			for _, text := range listOfText {
				_, err := filePtr.WriteString(text + "\n")
				if err != nil {
					log.Fatal(err)
				}
			}
			filePtr.Close()
			numLines = 0
			currentFile += 1
			listOfText = nil
			totalChunks += 1
		}
	}

	fileName := "./filestore/somebigfile_" + strconv.FormatInt(int64(currentFile), 10)
	filePtr, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	for _, text := range listOfText {
		_, err := filePtr.WriteString(text + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	filePtr.Close()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file.Close()
}

func lineCounter(file string) int {
	fileptr, _ := os.Open(file)
	fileScanner := bufio.NewScanner(fileptr)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}

func printall(args ...string) {

	for _, word := range args {
		fmt.Println(word)
	}
}

func start() {

	var inputText string = os.Args[1]

	var dagText string = os.Args[2]

	//commands, numberofCommands := readFile("./myinput/dag_input3.txt")
	commands, numberofCommands := readFile(inputText)
	dag1 := dag.NewDAG()

	matrix := createDag2(dag1, commands, numberofCommands)
	// fmt.Println(dag1)

	// fmt.Println(matrix)
	master := Master{make(chan string, 100)}

	//can pick between two schedulers for time or not
	scheduler(matrix, dag1, master, dagText)
	//schedulerTimed(matrix, dag1, master, dagText)

}
