package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const tmpDir = "./files"

func main() {

	welcomeData := []byte("Wena los cabros.")

	path := filepath.Join(tmpDir, "/welcome.txt")

	err := ioutil.WriteFile(path, welcomeData, 0777)

	if err != nil {
		fmt.Println(err)
	}

}