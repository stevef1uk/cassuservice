package handler

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)


const FILETOPROCESS = "restapi/configure_simple.go"  // Name of file to update created by go-swagger during generation


func tempFile() *os.File {

	file, err := ioutil.TempFile(os.TempDir(), "stdin")
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func input(in *os.File)  {
	if in != nil {
		os.Stdin = in
	}
}

func output(in *os.File)  {
	if in != nil {
		os.Stdout = in
	}
}

// Function to create the temporary file to include updated contents of the generated handler
func openFile(  debug bool, genDir string ) *os.File {

	// Create the directory if not already there
	fullfileName := genDir + "/" + MAINFILE
	if debug { fmt.Println("Dir Name = ", fullfileName )}
	// Create data dir if it doesn't already exist
	_, err := os.Stat(fullfileName)
	if ( err != nil ) {
		if os.IsNotExist(err) {
			// This file shoud exist
			if debug { fmt.Println("Have you generated the code as the file" + FILETOPROCESS + " is not foumd" ) }
			panic("You need to use go-swagger to generate the go-swagger files before trying to generate the Cassandra Handler code")
		}
	}
	myfile, err := os.Open(fullfileName )
	if ( err != nil ) {
		panic(err)
	}

	return myfile
}

// Create the final file deleting it if it exists by linking the passed temporary file to it and delete temporary file
func createFile( generatedCodePath string, tmpFile string  ) {
	fullfileName := generatedCodePath + "/" + FILETOPROCESS
	var err = os.Remove(fullfileName)
	if err != nil {
		panic(err)
	}
	err = os.Link(tmpFile, fullfileName)
	if err != nil {
		panic(err)
	}

	err = os.Remove(tmpFile)
	if err != nil {
		panic(err)
	}
}

// Only enable debug here when in difficulty as the debug strings will end up in the generated file causing compilation issues
func Parse( debug bool, generatedCodePath string, tableName string, endPointNameOverRide string) {
	reprocessing := false
	genString := ""
	if (endPointNameOverRide != "" ) {
		genString = strings.Title(endPointNameOverRide)
	} else {
		genString = tableName
	}
	if debug {fmt.Println("handler.Parse  genString = ", genString)
	}
	handlerString := "api.Get" + genString + "Handler = operations.Get" + genString

	fileout := tempFile()
	defer fileout.Close()
	if debug { fmt.Println("created file: ", fileout.Name())
	}
	//fmt.Println("created file: " , fileout.Name() )

	input(openFile(debug, generatedCodePath))
	output(fileout)

	skip := false
	reader := bufio.NewReader(os.Stdin)
	if debug { fmt.Println("Parsing input read = ", handlerString)}

for {
		text, err := reader.ReadString('\n')
		if ( err != nil ) {
			if debug {fmt.Println("Err = ", err) }
			break
		}

		// Stop this being run twice!
		if ( strings.Contains(text, "data.Search(params)")  ) {
			log.Println("Already updated ", FILETOPROCESS)
			reprocessing = true
			break
		}

		if ( skip ) {
			skip = false
			continue
		}

		// Remove newline character first
		text = strings.Replace(text, "\n", "", -1)

		// Add Cassandra Shutdown hook
		if ( strings.Contains(text, "api.ServerShutdown = func() {}") ) {
			if debug { fmt.Print("Found Shutdown")}
			fmt.Println(`
        api.ServerShutdown = func() {
        data.Stop()
    }
`)
		} else if ( strings.Contains(text, "func setupMiddlewares(handler") ) {
			if debug { fmt.Print("Found setupMiddlewares") }
			fmt.Println(text)
			fmt.Println(`
        data.SetUp()
`)
			} else if ( strings.Contains(text, handlerString) ) {
				if debug {
					fmt.Print("Found hadlerString")
				}
				fmt.Println(text)
				skip = true
				//fmt.Println(text)

				fmt.Println(`
	return data.Search(params)
`)
				} else if ( strings.Contains(text, "restapi/operations") ) {
					if debug {
						fmt.Print("Found import")
					}
					tmpStr := strings.Replace(text, "restapi/operations", "data", -1)
					fmt.Println(text)
					fmt.Println(tmpStr + `
`)
				} else {
					fmt.Println(text)
				}
	}
	if ( reprocessing == false ) {

		createFile(generatedCodePath, fileout.Name())
	}
}