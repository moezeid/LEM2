package LEM2

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DieOnError(err error) {

	if err != nil {
		panic(err)
	}
}
type Env struct{
	AttributeMap map[string][]*AttributeObject

	DecisionMap map[string][]int

	AttributeValueBlock map[Tuple][]int

	AttributeList []string

	CasesCovered int

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

func(e *Env) Parse() {

	dir, err := os.Getwd()
	path := filepath.Join(dir, "/dataset/breast.txt")
	f, err := os.Open(path)
	DieOnError(err)
	scanner := bufio.NewScanner(f)
	numAttributes := 0

	//b := make([]byte,1)
	num := 0

	for scanner.Scan() {
		//r is the entire line read in
		r := bytes.NewBuffer(scanner.Bytes())

		//this loop reads the first line that looks like " < a, a, ..., d >
		for num == 0 {

			//now one character (byte) is read from that line
			b, err := r.ReadByte()
			DieOnError(err)
			s := string(b)
			fmt.Print(s)
			if s == "a" {
				numAttributes++
			} else if s == "d" {
				fmt.Print("\n")
				break
			}

		}

		//this loop reads the second line which has the names of the attributes along with decisions
		if num == 1 {
			//fmt.Println(r.String())

			line := strings.SplitAfter(r.String(), " ")
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
			//v1.String()
			t := Tuple{
				Attribute: v1.attribute,
				Value:     v1.value,
			}

			if _, ok := e.AttributeValueBlock[t]; ok {
				e.AttributeValueBlock[t] = append(e.AttributeValueBlock[t], v1.caseNum)
			} else {
				e.AttributeValueBlock[t] = []int{v1.caseNum}
			}
		}
	}

	//prints [(a,v)]

	for i, v := range e.AttributeValueBlock {
		i.String()
		fmt.Printf(" covers cases: %d \n", v)
	}

	for i, v := range e.DecisionMap {
		fmt.Printf("(decision,%s): %d\n", i, v)
	}

}
