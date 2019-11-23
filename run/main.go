package main

import (
	"fmt"
	"os"
	"path"

	"github.com/LEM2"
)

func main() {
	e := &LEM2.Env{

		AttributeMap:        make(map[string][]*LEM2.AttributeObject),
		DecisionMap:         make(map[string][]int),
		AttributeList:       make([]string, 0),
		AttributeValueBlock: make(map[LEM2.Tuple][]int),
		CasesCovered: 0,
	}



	fmt.Print("Enter the path to the input file: ")

	var input string
	var err error

	_, err = fmt.Scanln(&input)
	for err != nil{
		fmt.Print("Please enter a valid filepath: ")
		_, err = fmt.Scanln(&input)
	}




	e.FilePath = input
	//fmt.Println(et)
	fmt.Println("hey")

	list := e.Algorithm()
	dir,_ := os.Getwd()
	newFilePath := path.Join(dir,"MLEM2_output_"+path.Base(e.FilePath))
	e.FileEnv , err = os.Create(newFilePath)
	LEM2.DieOnError(err)
	list.String()
	_, err = e.FileEnv.Write(LEM2.OutputFile.Bytes())
	LEM2.DieOnError(err)
	LEM2.DieOnError(e.FileEnv.Close())

//	LEM2.OutputFile.Close()

//	outC := make(chan string)



	//e.OutputFile.Write(list.String())
	//e.RuleCheck(list,"no-recurrence-events")

	/*


	r, i1 := e.ColumnScan(t,i)
	fmt.Println(r)
	r,_ = e.ColumnScan(t,*i1)
	fmt.Println(r)

	r, _ = e.ColumnScan(t,*i1)
	fmt.Println(r)

	 */

}

