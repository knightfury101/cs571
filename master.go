package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Master struct {
	files chan string
}

var wg sync.WaitGroup

//wg will need to have had called the add() method before calling the run command
func (mr *Master) RunCommand_experiment(commandString string, fileIn string) string {

	data := PerformCommmand_experiment(commandString, fileIn)
	commandArray := strings.SplitN(commandString, " ", 2)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	filename := randomInputFileName(r1.Intn(1000000), commandArray[0])

	// can get rid of this line if I want to
	filename = "./filestore/" + filename

	writeToFile(data, filename)

	return filename

	//The 2 value will need to be removed and replaced with the number of children a particular dag node has
	//filenames := splitChunks(filename, 3)

	//add the names of the file chunks
	// for _, s := range filenames {
	// 	mr.addFile(s)
	// }

	// wg.Done()
}

//wg will need to have had called the add() method before calling the run command
func (mr *Master) RunCommand(command string, input string, fileIn string) string {

	data := PerformCommmand(command, input, fileIn)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	filename := randomInputFileName(r1.Intn(1000000), command)

	// can get rid of this line if I want to
	filename = "./filestore/" + filename
	//

	writeToFile(data, filename)

	return filename
}

//Need to rewrite this as of april 7 2022
func (mr *Master) merge() {

	var buf bytes.Buffer
	for file := range mr.files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			// handle error
		}

		buf.Write(b)
	}

	err := ioutil.WriteFile("./filestore/OUPUT_file.txt", buf.Bytes(), 0644)
	if err != nil {
		// handle error
	}
}
