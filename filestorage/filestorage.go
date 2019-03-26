package filestorage

import "os"
type FileStorage struct{
	FileWriter os.File	
}

func (fileStorage *FileStorage) Init(filename string, shouldClearFileContent bool){
	if shouldClearFileContent == true {
		os.Remove(filename)
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	
	fileStorage.FileWriter = *f
}

func (fileStorage *FileStorage) AppendLine(content string){
	const sep = "\n"
	fileStorage.FileWriter.WriteString(content + sep)
}