package LEM2

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//"strings"
)

func DieOnError(err error) {

	if err != nil {
		panic(err)
	}
}



type attributeObject struct{
	attribute	string
	caseNum		int
	value 		string
}

type tuple struct{
	attribute string
	value	  string

}

func (a *attributeObject) String(){
	fmt.Printf("attribute: %s, value: %s case number: %d\n",a.attribute,a.value,a.caseNum)
}
func (t *tuple) String(){
	fmt.Printf("(%s, %s) ",t.attribute,t.value)
}

func isSpecialCharacter(s string) bool{
	list := []string{" ", "[", "]", "!", "/", ""}
	for _, v := range list{
		if s == v{
			return true
		}
	}
	return false
}


func LEM2() {
/*
	if err != nil{
		panic(err)
	}


	b := make([]byte,1)


	*/


	dir, err := os.Getwd()
	path := filepath.Join(dir,"/dataset/test.txt")
	f, err := os.Open(path)
	DieOnError(err)
	scanner := bufio.NewScanner(f)
	numAttributes := 0

	//b := make([]byte,1)
	num := 0
	attributeMap := make(map[string][]*attributeObject)
	decisionMap := make(map[string][]int)
	attributeList := make([]string,0)
	attributeValueBlock := make(map[tuple][]int)



	for scanner.Scan(){
		//r is the entire line read in
		r := bytes.NewBuffer(scanner.Bytes())
		//if num == 2 {fmt.Println(r)}

		for num == 0{

			//now one character (byte) is read from that line
			b, err  := r.ReadByte()
			DieOnError(err)
			s := string(b)
			fmt.Print(s)
			if s == "a"{
				numAttributes++
			}else if s == "d"  {
					fmt.Print("\n")
					break
			}

		}

		if num == 1{
			//fmt.Println(r.String())

			att := strings.SplitAfter(r.String()," ")
			for _, v := range att{
				if t := strings.TrimSpace(v); t != "d" && !isSpecialCharacter(t){
					attributeList = append(attributeList,t)
					}
			}

		}

		if num > 1{
			caseNum := num - 1

			att := strings.SplitAfter(r.String()," ")
			//l := len(att)
			for i, v := range att {
					t := strings.TrimSpace(v)
					if i == numAttributes {
						if _, ok := decisionMap[t]; ok {
							decisionMap[t] = append(decisionMap[t], caseNum)
							}else{
								decisionMap[t] = []int{caseNum}
						}
					}else if !isSpecialCharacter(t){
						a := new(attributeObject)
						if i  <= numAttributes-1 {
						attributeName := attributeList[i]
						a.attribute = attributeName
						a.caseNum = caseNum
						a.value = t
						attributeMap[attributeName] = append(attributeMap[attributeName],a)
						}
						//fmt.Print(v)
					}
				}
			}
		num++
		//if num == 9 {break}

		}

	for _, v := range attributeMap{
		for _, v1 := range v{
			//v1.String()
			t := tuple{
				attribute: v1.attribute,
				value: v1.value,
			}

			if _ ,ok := attributeValueBlock[t]; ok{
				attributeValueBlock[t] = append(attributeValueBlock[t],v1.caseNum)
			}else{
				attributeValueBlock[t] = []int{v1.caseNum}
			}
		}
	}


for i,v := range attributeValueBlock{
	i.String()
	fmt.Printf("%d \n",v)
}

}

