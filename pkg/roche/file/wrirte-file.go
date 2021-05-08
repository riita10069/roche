package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

func makeDirectory (filename string) {
	paths := strings.Split(filename, "/")
	paths = paths[:len(paths)-1]
	filepath := strings.Join(paths, "/")
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		defaultUmask := syscall.Umask(0)
		os.MkdirAll(filepath, 0777)
		os.Chmod(filepath, 0777)
		syscall.Umask(defaultUmask)
	}
	fmt.Println("Execute Make Directory")
	fmt.Println("     ✔ ", filepath)
}

func CreateAndWrite(content string, filename string) error {
	if strings.Contains(filename, "/") {
		makeDirectory(filename)
	}
	err := ioutil.WriteFile(filename, []byte(content), 0664)
	if err != nil {
		return errors.New("cannot write" + filename)
	}
	fmt.Println("Execute code generate to following file")
	fmt.Println("     ✔ ", filename)
	return nil
}
