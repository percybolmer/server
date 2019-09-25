package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
)

type newServer struct {
	PackageName string
	ServerName  string
	Location    string
	Db          bool
	DbDriver    string
	Gorm        bool
}

func main() {
	var d newServer

	flag.StringVar(&d.PackageName, "package", "", "The name for the generated package")
	flag.StringVar(&d.ServerName, "server", "", "The name of the server object to generate")
	flag.BoolVar(&d.Db, "db", false, "Set true if its a database server")
	flag.BoolVar(&d.Gorm, "gorm", false, "Set true if the server should use gorm as db handler")
	flag.StringVar(&d.Location, "location", "", "Location of the generated output file/Files")
	flag.StringVar(&d.DbDriver, "driver", "postgres", "The database driver to use, defaults to postgres")
	flag.Parse()

	finfo := getFileStat(d.Location)
	if !finfo.IsDir() {
		log.Fatal("Please specify a directory to generate files in, not a file")
		return
	}
	// If Gorm is set, then DB is not needed to be set
	if d.Gorm {
		d.Db = true
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.go", d.Location, d.ServerName))
	checkError(err)
	testf, err := os.Create(fmt.Sprintf("%s/%s_test.go", d.Location, d.ServerName))
	checkError(err)
	dockerf, err := os.Create(fmt.Sprintf("%s/%s", d.Location, "Dockerfile"))
	checkError(err)

	temps, err := template.ParseFiles("templates/serverTemplate.gohtml", "templates/db_funcs.gohtml")
	checkError(err)

	tests, err := template.ParseFiles("templates/tests/testTemplate.gohtml", "templates/tests/db_tests.gohtml")
	checkError(err)

	dock, err := template.ParseFiles("templates/docker.gohtml")
	checkError(err)

	err = temps.Execute(f, d)
	checkError(err)
	err = tests.Execute(testf, d)
	checkError(err)
	err = dock.Execute(dockerf, d)
	checkError(err)

	cmd := exec.Command("gofmt", "-w", d.Location)
	err = cmd.Run()
	checkError(err)

}

// getFileStat will get filestats or create a Directory,
// will trigger error if something goes wrong
func getFileStat(location string) os.FileInfo {
	finfo, err := os.Stat(location)
	if os.IsNotExist(err) {
		err := os.Mkdir(location, 0755)
		checkError(err)
		finfo, err = os.Stat(location)
		checkError(err)
		return finfo
	}
	return finfo
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
