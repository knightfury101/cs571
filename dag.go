package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/goombaio/dag"
)

//NOT NEEDED
func createDag(commands []string, numberofCommands []string, numberOfChildList []string) {

	dagCommand := dag.NewDAG()
	createCommandMatrix(dagCommand, commands, numberofCommands)
	// vertexchan := make(chan *dag.Vertex, 10)
	// currentCommand := 0

	// headVertex := dag.NewVertex("head", nil)
	// dagCommand.AddVertex(headVertex)
	// curVertex := headVertex
	// previousTimesToRun := 1

	// for {

	// 	numberOfTimesToRun, err := strconv.Atoi(numberof[currentCommand])

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// }
	// vertex := dag.NewVertex(commands[0], nil)
	// vertexchan <- vertex
	// value := <-vertexchan
	// fmt.Println(value.String())
	// vertex1 := dag.NewVertex(commands[1], nil)
	// dagCommand.AddVertex(vertex)
	// dagCommand.AddVertex(vertex1)
	// dagCommand.AddEdge(vertex, vertex1)

	// fmt.Println(dagCommand.String())

}

func readFile(path string) ([]string, []string) {

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	stringOfCommands := strings.Split(scanner.Text(), ",")
	scanner.Scan()
	numberOfCommands := strings.Split(scanner.Text(), ",")

	return stringOfCommands, numberOfCommands

}

func createDag2(dagCommand *dag.DAG, commands []string, numberofCommandsList []string) [][]*dag.Vertex {
	//dagCommand := dag.NewDAG()
	commMatrix := createCommandMatrix(dagCommand, commands, numberofCommandsList)

	// fmt.Println(commMatrix)

	// fmt.Println()
	// fmt.Println()
	// fmt.Println("+++++++++++++++++++++++++++++")

	//commandRow := 0

	for i := 0; i < len(commands)-1; i++ {

		var commList []*dag.Vertex
		var childList []*dag.Vertex
		//var currCommPos int = 0

		j := 0
		val := getNum(numberofCommandsList[i])

		for j < val {
			commList = append(commList, commMatrix[i][j])
			j++
		}

		j = 0
		val = getNum(numberofCommandsList[i+1])

		for j < val {
			childList = append(childList, commMatrix[i+1][j])
			j++
		}

		// fmt.Println(commList)

		// fmt.Println(childList)
		// fmt.Println()
		// fmt.Println()
		// fmt.Println("---------------------")
		// fmt.Println()
		// fmt.Println()

		if len(commList) >= len(childList) {

			curChild := 0
			//numChildUsed := 0
			curComm := 0
			numOfChild := len(childList)
			numOfComm := len(commList)

			modVal := numOfComm % numOfChild
			whenToUpdate := 0
			if modVal != 0 {
				whenToUpdate = modVal
			} else {
				whenToUpdate = numOfComm / numOfChild
			}

			for k := 0; k < len(commList); k++ {

				if curComm == whenToUpdate && curChild != (len(childList)-1) {
					curChild += 1
					curComm = 0
					numOfComm -= whenToUpdate
					numOfChild -= 1

					modVal = numOfComm % numOfChild

					if modVal != 0 {
						whenToUpdate = modVal
					} else {
						whenToUpdate = numOfComm / numOfChild
					}

				}

				dagCommand.AddEdge(commList[k], childList[curChild])

				// fmt.Println(commList[k])
				// fmt.Println("EDGE TO ---> ")
				// fmt.Println(childList[curChild])

				// fmt.Println()
				// fmt.Print("---------------------------------------\n")
				curComm += 1
			}

		} else {

			curComm := 0
			curChild := 0
			numOfChild := len(childList)
			numOfComm := len(commList)

			modVal := numOfChild % numOfComm

			whenToUpdate := 0
			if modVal != 0 {
				whenToUpdate = modVal
			} else {
				whenToUpdate = numOfChild / numOfComm
			}

			for k := 0; k < len(childList); k++ {

				if curChild == whenToUpdate && curComm != (len(commList)-1) {
					curComm += 1
					curChild = 0
					/////////
					numOfChild -= whenToUpdate
					numOfComm -= 1

					modVal = numOfChild % numOfComm

					if modVal != 0 {
						whenToUpdate = modVal
					} else {
						whenToUpdate = numOfChild / numOfComm
					}
				}

				dagCommand.AddEdge(commList[curComm], childList[k])
				curChild += 1
			}
		}

	}

	// fmt.Println("((((((((((((((-------------")
	// fmt.Println()
	// fmt.Println()
	//fmt.Println(dagCommand.String())

	return commMatrix

}

func printDag(curdag *dag.DAG, curverts []*dag.Vertex) {

	for _, item := range curverts {
		fmt.Println("++++++++++++++++++++++++++++++++++++")
		printLine(curdag, item)

	}
}

func printLine(curdag *dag.DAG, curvert *dag.Vertex) {

	fmt.Println(curvert)
	fmt.Println()
	fmt.Println("  ----->")
	successors, err := curdag.Successors(curvert)

	if err != nil {
		return
	}

	for i := 0; i < len(successors); i++ {
		printLine(curdag, successors[i])
	}
}

func getNum(value string) int {
	val, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return val
}

func createCommandMatrix(curdag *dag.DAG, commands []string, numberOfCommandsList []string) [][]*dag.Vertex {

	matrixX, matrixY := createMatrixValues(commands, numberOfCommandsList)

	commandMatrix := createMatrix(matrixX, matrixY)

	for i := 0; i < matrixX; i++ {
		for j := 0; j < matrixY; j++ {

			val, err := strconv.Atoi(numberOfCommandsList[i])

			if err != nil {
				panic(err)
			}

			if j >= val {
				break
			}
			id := strconv.Itoa(i) + strconv.Itoa(j)
			vertex := dag.NewVertex(id, commands[i])
			curdag.AddVertex(vertex)
			commandMatrix[i][j] = vertex
		}

	}

	//fmt.Println(curdag)

	return commandMatrix
}

func createMatrix(x int, y int) [][]*dag.Vertex {
	matrix := make([][]*dag.Vertex, x)
	for i := 0; i < x; i++ {
		matrix[i] = make([]*dag.Vertex, y)
	}

	return matrix

}

func createMatrixValues(commands []string, numberOfCommandsList []string) (int, int) {

	var maxY int = 0

	for i := 0; i < len(commands); i++ {
		val, err := strconv.Atoi(numberOfCommandsList[i])
		if err != nil {
			panic(err)
		}
		if val > maxY {
			maxY = val
		}
	}

	return len(commands), maxY

}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	file.Close()

	var commands []string
	var nums []string
	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if i%2 == 0 {
			commands = append(commands, scanner.Text())
			i++
		} else {
			nums = append(nums, scanner.Text())
			i++
		}
	}
	return commands, nums, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(commands []string, nums []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, command := range commands {
		fmt.Fprintln(w, command)
	}

	for _, num := range nums {
		fmt.Fprintln(w, num)
	}
	return w.Flush()
}

func prac() {
	commands, nums, err := readLines("test.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	if err := writeLines(commands, nums, "test.txt"); err != nil {
		log.Fatalf("writeLines: %s", err)
	}
	//numberOfNode := 0

	/////////////////////////////
	//fmt.Println(commands)
	//fmt.Println(nums)
	////////////////////////////
	//fmt.Println("Creating DAG...\n")

	/*
		// Create the dag
		dag1 := dag.NewDAG()

		// Create the vertices. Value is nil to simplify.
		vertex1 := dag.NewVertex("grep", 1)
		vertex2 := dag.NewVertex("sed1", 1)
		vertex3 := dag.NewVertex("sed2", 1)
		vertex4 := dag.NewVertex("awk1", 1)
		vertex5 := dag.NewVertex("awk2", 1)
		vertex6 := dag.NewVertex("awk3", 1)
		vertex7 := dag.NewVertex("awk4", 1)
		vertex8 := dag.NewVertex("cat", 1)

		// Add the vertices to the dag.
		dag1.AddVertex(vertex1)
		dag1.AddVertex(vertex2)
		dag1.AddVertex(vertex3)
		dag1.AddVertex(vertex4)
		dag1.AddVertex(vertex5)
		dag1.AddVertex(vertex6)
		dag1.AddVertex(vertex7)
		dag1.AddVertex(vertex8)

		// Add the edges (Note that given vertices must exist before adding an
		// edge between them).
		dag1.AddEdge(vertex1, vertex2)
		dag1.AddEdge(vertex1, vertex3)
		dag1.AddEdge(vertex2, vertex4)
		dag1.AddEdge(vertex2, vertex5)
		dag1.AddEdge(vertex3, vertex6) */
}
