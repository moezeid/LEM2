package LEM2

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func DieOnError(err error) {

	if err != nil {
		panic(err)
	}
}

type Env struct {
	AttributeMap map[string][]*AttributeObject

	DecisionMap map[string][]int

	AttributeValueBlock map[Tuple][]int

	AttributeList []string

	CasesCovered int

	NumericMap map[string][]float64

	FilePath   string
}

type AttributeObject struct {
	attribute string
	caseNum   int
	value     string
}

type Tuple struct {
	Attribute string
	Value     string
}

func (a *AttributeObject) String() {
	fmt.Printf("Attribute: %s, Value: %s case number: %d\n", a.attribute, a.value, a.caseNum)
}
func (t *Tuple) String() {
	fmt.Printf("(%s, %s) ", t.Attribute, t.Value)
}

func isSpecialCharacter(s string) bool {
	list := []string{" ", "[", "]", "!", "/", ""}
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}

func (e *Env) Parse() {

	listOfCases := make([]int,0)
	dir, err := os.Getwd()
	path := filepath.Join(dir, "/dataset/wine.txt")

	f, err := os.Open(path)
	for err != nil{
		var input string
		fmt.Print("Please enter a valid filepath: ")
		_, err = fmt.Scanln(&input)
		path := filepath.Join(dir,input)
		f, err = os.Open(path)

	}

	DieOnError(err)
	scanner := bufio.NewScanner(f)
	numAttributes := 0

	//b := make([]byte,1)
	num := 0
	e.NumericMap = make(map[string][]float64)

	for scanner.Scan() {
		//r is the entire line read in
		r := bytes.NewBuffer(scanner.Bytes())

		//this loop reads the first line that looks like " < a, a, ..., d >
		for num == 0 {

			//now one character (byte) is read from that line
			b, err := r.ReadByte()
			DieOnError(err)
			s := string(b)
			if s == "a" {
				numAttributes++
			} else if s == "d" {
				break
			}

		}

		//this loop reads the second line which has the names of the attributes along with decisions
		if num == 1 {

			line := strings.SplitAfter(r.String(),
				" ")
			for _, v := range line {
				if t := strings.TrimSpace(v); t != "d" && !isSpecialCharacter(t) {
					e.AttributeList = append(e.AttributeList, t)
				}
			}

		}

		//this loop retrieves the rest of the data by parsing attributes
		//from a case horizontally. It is also in charge of creating
		//mappings for attributes
		if num > 1 {
			caseNum := num - 1
			listOfCases = append(listOfCases,caseNum)
			line := strings.SplitAfter(r.String(), " ")
			//l := len(line)
			attributesCollected := 0
			for _, v := range line {

				t := strings.TrimSpace(v)
				if !isSpecialCharacter(t) && attributesCollected == numAttributes {
					if _, ok := e.DecisionMap[t]; ok {
						e.DecisionMap[t] = append(e.DecisionMap[t], caseNum)
					} else {
						e.DecisionMap[t] = []int{caseNum}
					}
				} else if !isSpecialCharacter(t) {
					a := new(AttributeObject)
					if attributesCollected < numAttributes {
						attributeName := e.AttributeList[attributesCollected]
						a.attribute = attributeName
						a.caseNum = caseNum
						a.value = t
						tuple := Tuple{
							a.attribute,
							a.value,
						}
						if isNumeric(tuple) {
							//append value to map of list of numeric values

							//convert string value to float value
							num, _ := strconv.ParseFloat(tuple.Value, 64)
							if _, ok := e.NumericMap[tuple.Attribute]; ok {
								//append to list if the list already exists but check for redundancy
								if !listContains(e.NumericMap[tuple.Attribute], num) {
									e.NumericMap[tuple.Attribute] = append(e.NumericMap[tuple.Attribute], num)
									e.NumericMap[tuple.Attribute] = bubbleSort(e.NumericMap[tuple.Attribute])
								}
							} else {
								//creating the list if it does not yet exist
								e.NumericMap[tuple.Attribute] = []float64{num}
							}
						}
						e.AttributeMap[attributeName] = append(e.AttributeMap[attributeName], a)
						attributesCollected++
					}
				}
			}
		}
		num++
		//if num == 9 {break}

	}

	//after all the Attribute mappings are set up we
	//are free to create Attribute-Value blocks
	for _, v := range e.AttributeMap {
		for _, v1 := range v {
			t := Tuple{
				Attribute: v1.attribute,
				Value:     v1.value,
			}
			if isNumeric(t)  {
				numericList := e.NumericMap[t.Attribute]
				cP := createCutPoints(numericList)
				for i := range cP {
					low := fmt.Sprintf("%f", numericList[0])
					high1 := fmt.Sprintf("%f", cP[i])
					high2 := fmt.Sprintf("%f", numericList[len(numericList)-1])
					firstPartOfInterval := low + ".." + high1
					secondPartOfInterval := high1 + ".." + high2
					firstInterval := Tuple{
						t.Attribute,
						firstPartOfInterval,
					}
					remainderInterval := Tuple{
						t.Attribute,
						secondPartOfInterval,
					}
					setsCoveredByFirst := e.FindCasesForInterval(numericList[0], cP[i])
					setsCoveredBySecond := e.reduceDecisionsSet(setsCoveredByFirst,listOfCases)
					e.AttributeValueBlock[firstInterval] = setsCoveredByFirst
					e.AttributeValueBlock[remainderInterval] = setsCoveredBySecond

				}
				break

			} else {

				if _, ok := e.AttributeValueBlock[t]; ok {
					e.AttributeValueBlock[t] = append(e.AttributeValueBlock[t], v1.caseNum)
				} else {
					e.AttributeValueBlock[t] = []int{v1.caseNum}
				}
			}

		}

	}



}

func isNumeric(attribute Tuple) bool {

	lb := strings.Split(attribute.Value, "..")
	if lb[0] == attribute.Value {
		_, err := strconv.ParseFloat(lb[0], 64)
		return err == nil
	}
	return false

}

func listContains(list []float64, num float64) bool {
	for i := range list {
		if list[i] == num {
			return true
		}

	}
	return false
}

//pulled from https://tutorialedge.net/golang/implementing-bubble-sort-with-golang/
func bubbleSort(input []float64) []float64 {
	// n is the number of items in our list
	n := len(input)
	// set swapped to true
	swapped := true
	// loop
	for swapped {
		// set swapped to false
		swapped = false
		// iterate through all of the elements in our list
		for i := 1; i < n; i++ {
			// if the current element is greater than the next
			// element, swap them
			if input[i-1] > input[i] {
				// log that we are swapping values for posterity
				// swap values using Go's tuple assignment
				input[i], input[i-1] = input[i-1], input[i]
				// set swapped to true - this is important
				// if the loop ends and swapped is still equal
				// to false, our algorithm will assume the list is
				// fully sorted.
				swapped = true
			}
		}
	}
	return input
}

func createCutPoints(list []float64) []float64 {
	//length := len(list)

	cutPoints := make([]float64, 0)
	for i := range list{
		if i + 1 ==len(list){
			break
		}
		cutPoints = append(cutPoints,(list[i] + list[i+1]) / 2)
	}
	return cutPoints

}

func (e *Env) FindCasesForInterval(low, high float64) []int {

	casesCovered := make([]int, 0)
	for _, v := range e.AttributeMap {

		for _, v1 := range v {
			t := Tuple{
				v1.attribute,
				v1.value,
			}
			if isNumeric(t) {
				num, _ := strconv.ParseFloat(t.Value, 64)
				if num >= low && num <= high {
					casesCovered = append(casesCovered, v1.caseNum)
				}

			}

		}
	}
	return casesCovered
}
