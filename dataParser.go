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

type attributeObject struct {
	attribute string
	caseNum   int
	value     string
}

type tuple struct {
	attribute string
	value     string
}

func (a *attributeObject) String() {
	fmt.Printf("attribute: %s, value: %s case number: %d\n", a.attribute, a.value, a.caseNum)
}
func (t *tuple) String() {
	fmt.Printf("(%s, %s) ", t.attribute, t.value)
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

func Parse() {

	dir, err := os.Getwd()
	path := filepath.Join(dir, "/dataset/bowl.txt")
	f, err := os.Open(path)
	DieOnError(err)
	scanner := bufio.NewScanner(f)
	numAttributes := 0

	//b := make([]byte,1)
	num := 0
	attributeMap := make(map[string][]*attributeObject)
	decisionMap := make(map[string][]int)
	attributeList := make([]string, 0)
	attributeValueBlock := make(map[tuple][]int)

	for scanner.Scan() {
		//r is the entire line read in
		r := bytes.NewBuffer(scanner.Bytes())
		//if num == 2 {fmt.Println(r)}

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
					attributeList = append(attributeList, t)
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
					if _, ok := decisionMap[t]; ok {
						decisionMap[t] = append(decisionMap[t], caseNum)
					} else {
						decisionMap[t] = []int{caseNum}
					}
				} else if !isSpecialCharacter(t) {
					a := new(attributeObject)
					if attributesCollected < numAttributes {
						attributeName := attributeList[attributesCollected]
						a.attribute = attributeName
						a.caseNum = caseNum
						a.value = t
						attributeMap[attributeName] = append(attributeMap[attributeName], a)
						attributesCollected++
					}
				}
			}
		}
		num++
		//if num == 9 {break}

	}

	//after all the attribute mappings are set up we
	//are free to create attribute-value blocks
	for _, v := range attributeMap {
		for _, v1 := range v {
			//v1.String()
			t := tuple{
				attribute: v1.attribute,
				value:     v1.value,
			}

			if _, ok := attributeValueBlock[t]; ok {
				attributeValueBlock[t] = append(attributeValueBlock[t], v1.caseNum)
			} else {
				attributeValueBlock[t] = []int{v1.caseNum}
			}
		}
	}

	//prints [(a,v)]

	for i, v := range attributeValueBlock {
		i.String()
		fmt.Printf(" covers cases: %d \n", v)
	}

	for i, v := range decisionMap {
		fmt.Printf("(decision,%s): %d\n", i, v)
	}

}
