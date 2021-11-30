package core

import (
	"bytes"
	"fmt"
	"os"
)

const pathPreFix = "./saves/maps/"
const mapSaveFileSufix = ".mapsave"

func saveBufferToFile(path string, buffer *bytes.Buffer) {
	f, err := os.Create(path)
	checkErr(err)

	_, err = f.Write(buffer.Bytes())
	checkErr(err)

	defer checkErr(f.Close())
}

func saveMapBufferToFile(name string, buffer *bytes.Buffer) {
	saveBufferToFile(pathPreFix+name+mapSaveFileSufix, buffer)
}

func loadBufferFromFile(path string) *bytes.Buffer {
	buf, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("File %s not found.", path)
		return nil
	}
	return bytes.NewBuffer(buf)
}

func loadMapBufferFromFile(name string) *bytes.Buffer {
	return loadBufferFromFile(pathPreFix + name + mapSaveFileSufix)
}
