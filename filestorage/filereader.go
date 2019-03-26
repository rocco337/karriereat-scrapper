package filestorage
import "os"

import "bufio"
import "fmt"
type FileReader struct{
	reader *bufio.Reader
	scanner *bufio.Scanner
}

func (fileReader *FileReader) Init(filename string){
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file: %v\n",err)
		os.Exit(1)
	}

	fileReader.scanner = bufio.NewScanner(f)
	fileReader.scanner.Split(bufio.ScanLines)	
}

// func (fileReader *FileReader) ReadAll(){
// 	for fileReader.scanner.Scan() {
// 		fmt.Println(fileReader.scanner.Text())
// 	}
// }

func (fileReader *FileReader) ReadLine() (string, bool){
	if fileReader.scanner.Scan() ==true{
		return fileReader.scanner.Text(), false
	}

	return "", true
}

