package LEM2

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func DieOnError(err error) {

	if err != nil {
		panic(err)
	}
}

func LEM2() {
/*
	if err != nil{
		panic(err)
	}


	b := make([]byte,1)


	*/

	dir, err := os.Getwd()
	path := filepath.Join(dir,"/dataset/wine.txt")
	f, err := os.Open(path)
	DieOnError(err)
	scanner := bufio.NewScanner(f)
	numAttributes := 0
	for scanner.Scan(){
		b := make([]byte,1)
		r := bytes.NewBuffer(scanner.Bytes())
		fmt.Println(r.String())

		//read the "<" character
		_, err = r.Read(b)
		DieOnError(err)

		//read the whitespace character
		_, err = r.Read(b)
		DieOnError(err)

		for true {
			_, err = r.Read(b)
			DieOnError(err)
			//if letter "a" then count the attribute
			if string(b) == "a" {
				numAttributes++
			}else if string(b) == "d"{
				break
			}


		}
		break

	}
	fmt.Println(numAttributes)
/*
	for true {
			_, err = f.Read(b)
			DieOnError(err)
			//if letter "a" then count the attribute
			if string(b) == "a" {
				numAttributes++
			}else if string(b) == "d"{
					break
			}


	}
	decisionIndex := numAttributes
	fmt.Println(decisionIndex)

	for true{
		_, err = f.Read(b)
		DieOnError(err)
		if string(b) == "\n"{
			fmt.Println(decisionIndex,"heyyyyyy")
			break
		}
		decisionIndex++
	}

	//fmt.Print(string(data))
	DieOnError(err)


 */
}
