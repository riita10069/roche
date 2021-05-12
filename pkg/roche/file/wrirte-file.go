package file

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/jennifer/jen"
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

func Append(content string, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		return errors.New("can not open " + filename)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, content)
	if err != nil {
		return errors.New(err.Error())
	}

	fmt.Println("Execute code append to following file")
	fmt.Println("     ✔ ", filename)
	return nil
}

func JenniferToFile(f *jen.File, filename string) {
	buf := &bytes.Buffer{}
	err := f.Render(buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	CreateAndWrite(buf.String(), filename)

}
