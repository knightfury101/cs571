package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goombaio/dag"
)

var wgOne sync.WaitGroup

func schedulerTimed(dagMatrix [][]*dag.Vertex, dag *dag.DAG, ms Master, fileInput string) {

	//!!!!!!!!!!!!!!!!!!!!
	startTimeScheduler := time.Now()

	numSourceNodes := len(dag.SourceVertices())

	fmt.Println("Starting splitting")

	if numSourceNodes > 1 {

		//!!!!!!!!!!!!!!!!!!!!
		startInitialSplit := time.Now()
		listOfSourceFiles := splitChunksTimed(fileInput, "schedulerInput", "00", numSourceNodes)
		//!!!!!!!!!!!!!!!!!!!!
		duration := time.Since(startInitialSplit)
		logTimenotes("Inital Splitting finished in ..."+duration.String(), "logtime.txt")
		for _, file := range listOfSourceFiles {
			ms.files <- file
		}
	} else {
		ms.files <- fileInput
	}

	fmt.Println("Ending splitting")
	//////////////////////////////

	////////////////////////////
	for i := 0; i < len(dagMatrix); i++ {
		count := 0

		//!!!!!!!!!!!!!!!!!!!!
		rowstarttime := time.Now()

		for j := 0; j < len(dagMatrix[i]); j++ {

			//fmt.Println(dagCommand[i][j])

			if dagMatrix[i][j] == nil {
				continue
			}

			wgOne.Add(1)
			// fmt.Println(dagMatrix[i][j].OutDegree())
			// fmt.Println(dagMatrix[i][j].InDegree())

			//go performAction(dagMatrix[i][j], ms, dagMatrix[i][j].InDegree(), dagMatrix[i][j].OutDegree())
			go performAction_experimentTimed(dagMatrix[i][j], ms, dagMatrix[i][j].InDegree(), dagMatrix[i][j].OutDegree())
			count += 1
		}

		//waitOne would hold the number of elements in a row
		//that way we wait until all those elements in a row are done
		//fmt.Println("Done waiting in this section")
		wgOne.Wait()

		//!!!!!!!!!!!!!!!!!!!!
		duration := time.Since(rowstarttime)
		logTimenotes("ROW FINISHED RUNNING ITS THREADS IN..."+duration.String(), "logtime.txt")

	}

	fmt.Println("finished")
	close(ms.files)
	//!!!!!!!!!!!!!!!!!!!!
	duration := time.Since(startTimeScheduler)
	logTimenotes("Scheduler finished in ..."+duration.String(), "logtime.txt")
	//ms.merge()
	//////////////////////////////

}

func scheduler(dagMatrix [][]*dag.Vertex, dag *dag.DAG, ms Master, fileInput string) {

	numSourceNodes := len(dag.SourceVertices())

	fmt.Println("Starting splitting")

	if numSourceNodes > 1 {

		listOfSourceFiles := splitChunks(fileInput, "schedulerInput", "00", numSourceNodes)

		for _, file := range listOfSourceFiles {
			ms.files <- file
		}
	} else {
		ms.files <- fileInput
	}

	fmt.Println("Ending splitting")
	//////////////////////////////

	////////////////////////////
	for i := 0; i < len(dagMatrix); i++ {
		count := 0
		for j := 0; j < len(dagMatrix[i]); j++ {

			//fmt.Println(dagCommand[i][j])

			if dagMatrix[i][j] == nil {
				continue
			}

			wgOne.Add(1)
			// fmt.Println(dagMatrix[i][j].OutDegree())
			// fmt.Println(dagMatrix[i][j].InDegree())

			//go performAction(dagMatrix[i][j], ms, dagMatrix[i][j].InDegree(), dagMatrix[i][j].OutDegree())
			go performAction_experiment(dagMatrix[i][j], ms, dagMatrix[i][j].InDegree(), dagMatrix[i][j].OutDegree())
			count += 1
		}

		//waitOne would hold the number of elements in a row
		//that way we wait until all those elements in a row are done
		//fmt.Println("Done waiting in this section")
		wgOne.Wait()

	}

	fmt.Println("finished")
	close(ms.files)

	//ms.merge()
	//////////////////////////////

}

func performAction_experimentTimed(vertexCommand *dag.Vertex, ms Master, numParent int, numChild int) {
	//fmt.Println(vertexCommand)

	//!!!!!!!!!!!!!!
	startPerformAction := time.Now()
	commandList := strings.Split(vertexCommand.Value.(string), " ")

	//assume numInput can be greater than 1
	//so we will need to keep that in mind

	var filename string
	var listofFilesToMerge []string
	if numParent > 1 {

		for i := 0; i < numParent; i++ {
			file := <-ms.files
			listofFilesToMerge = append(listofFilesToMerge, file)
		}

		//here we merge the files and place it in filename
		//we do NOT want to place this file in the channel, it is to be used to merge
		//different files (look at the lines after the else statement) and run a different command only

		//!!!!!!!!!!!!!!
		startfileMerge := time.Now()
		fileMerge := mergeFiles(listofFilesToMerge, commandList[0])
		//!!!!!!!!!!!!!!
		duration := time.Since(startfileMerge)
		logTimenotes("Num parents : "+strconv.FormatInt(int64(numParent), 10)+" Num child: "+strconv.FormatInt(int64(numChild), 10)+" Time: "+sizeofFile(fileMerge)+" : for MERGING files: "+fileMerge+" == "+duration.String()+" to be finished...", "logtime.txt")

		//!!!!!!!!!!!!!!
		startRunCommand := time.Now()

		fmt.Println("Running Command multi" + vertexCommand.Value.(string))
		filename = ms.RunCommand_experiment(strings.Trim(vertexCommand.Value.(string), " "), fileMerge)
		fmt.Println("Finished RUNNING Command multi " + filename)

		//!!!!!!!!!!!!!!
		duration = time.Since(startRunCommand)
		logTimenotes("Num parents : "+strconv.FormatInt(int64(numParent), 10)+" Num child: "+strconv.FormatInt(int64(numChild), 10)+" Time: "+sizeofFile(filename)+" : RUNNING command on MERGED file: "+" mergefile: "+fileMerge+" outputfile: "+filename+" == "+duration.String()+" to be finished...", "logtime.txt")

	} else {
		file := <-ms.files
		//we do NOT want to place this file in the channel, it is to be used to merge
		//different files (look at the lines after the else statement) and run a different command only

		//!!!!!!!!!!!!!!
		startsingleParent := time.Now()

		fmt.Println("Running Command single" + vertexCommand.Value.(string))
		filename = ms.RunCommand_experiment(strings.Trim(vertexCommand.Value.(string), " "), file)
		fmt.Println("Finished RUNNING Command single" + vertexCommand.Value.(string))

		//!!!!!!!!!!!!!!
		duration := time.Since(startsingleParent)
		logTimenotes("Num parents : "+strconv.FormatInt(int64(numParent), 10)+" Num child: "+strconv.FormatInt(int64(numChild), 10)+" Time: "+sizeofFile(filename)+" :Time for RUNNING command on a SINGLE file: "+filename+" == "+duration.String()+" to be finished...", "logtime.txt")
	}

	// the list of files grabbed from if numParent > 1 should then be merged

	//!!!!!!!!!!!!!!
	finalsplitstarttime := time.Now()

	// ms.files <- filename
	var listofFiles []string
	if !(numChild == 0 || numChild == 1) {

		//listofFiles = splitChunks(filename, commandList[0], vertexCommand.ID, numChild)
		listofFiles = splitChunksTimed(filename, commandList[0], vertexCommand.ID, numChild)
		//listofFiles = splitChunks_experiment(filename, commandList[0], vertexCommand.ID, numChild)
		for i := 0; i < len(listofFiles); i++ {
			ms.files <- listofFiles[i]
		}
	} else {
		ms.files <- filename
	}

	//!!!!!!!!!!!!!!
	duration := time.Since(finalsplitstarttime)
	logTimenotes("Time: "+sizeofFile(filename)+" : for SPLITTING files into chunks for file: "+filename+" == "+duration.String()+" to be finished...", "logtime.txt")

	//fmt.Println(commandList)
	//ms.RunCommand()

	//!!!!!!!!!!!!!!
	duration = time.Since(startPerformAction)
	logTimenotes("Time for (PerformAction function finishing) FUNCTION AND FILE:"+filename+" == "+duration.String()+" to be finished...", "logtime.txt")

	wgOne.Done()
}

func performAction_experiment(vertexCommand *dag.Vertex, ms Master, numParent int, numChild int) {
	//fmt.Println(vertexCommand)

	commandList := strings.Split(vertexCommand.Value.(string), " ")

	//assume numInput can be greater than 1
	//so we will need to keep that in mind

	var filename string
	var listofFilesToMerge []string
	if numParent > 1 {

		for i := 0; i < numParent; i++ {
			file := <-ms.files
			listofFilesToMerge = append(listofFilesToMerge, file)
		}

		//here we merge the files and place it in filename
		//we do NOT want to place this file in the channel, it is to be used to merge
		//different files (look at the lines after the else statement) and run a different command only

		fileMerge := mergeFiles(listofFilesToMerge, commandList[0])

		fmt.Println("Running Command multi" + vertexCommand.Value.(string))
		filename = ms.RunCommand_experiment(strings.Trim(vertexCommand.Value.(string), " "), fileMerge)
		fmt.Println("Finished running Command multi" + vertexCommand.Value.(string))

	} else {
		file := <-ms.files
		//we do NOT want to place this file in the channel, it is to be used to merge
		//different files (look at the lines after the else statement) and run a different command only

		fmt.Println("Running Command single" + vertexCommand.Value.(string))
		filename = ms.RunCommand_experiment(strings.Trim(vertexCommand.Value.(string), " "), file)
		fmt.Println("Finished running Command single" + vertexCommand.Value.(string))

	}

	// the list of files grabbed from if numParent > 1 should then be merged

	// ms.files <- filename
	var listofFiles []string
	if !(numChild == 0 || numChild == 1) {

		//listofFiles = splitChunks(filename, commandList[0], vertexCommand.ID, numChild)
		listofFiles = splitChunks(filename, commandList[0], vertexCommand.ID, numChild)
		//listofFiles = splitChunks_experiment(filename, commandList[0], vertexCommand.ID, numChild)
		for i := 0; i < len(listofFiles); i++ {
			ms.files <- listofFiles[i]
		}
	} else {
		ms.files <- filename
	}
	//fmt.Println(commandList)
	//ms.RunCommand()

	wgOne.Done()
}

func logTimenotes(text string, file string) {

	f, err := os.OpenFile("./filestore/"+file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	newtext := text + "\n\n"
	f.WriteString(newtext)
	check(err)
	f.Close()
}

func sizeofFile(filename string) string {

	fi, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	// get the size
	size := fi.Size()
	s1 := strconv.FormatInt(int64(size), 10)

	return s1
}
