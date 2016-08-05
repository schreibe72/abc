package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/schreibe72/abc/storage"
)

var s storage.StorageAttributes
var defaultPath string
var loggerPtr *log.Logger

type Arguments struct {
	configPath      string
	verbose         bool
	Key             string
	Account         string
	container       string
	blobname        string
	filename        string
	containerPrefix string
	blobPrefix      string
	workercount     int
	pipe            bool
	big             bool
	contentSetting  storage.ContentSetting
}

func check(err error) {
	if err != nil {
		loggerPtr.Fatal(err)
	}
}

func upload(a Arguments) {

	if a.verbose {
		loggerPtr.Printf("Account: %s\nContainer: %s\n", a.Account, a.container)

	}

	s.CreateContainer(a.container)
	if a.pipe {
		if a.verbose {
			loggerPtr.Printf("Blobname: %s\n", a.blobname)
			loggerPtr.Printf("Big: %t\n", a.big)
		}
		err := s.SaveBlob(os.Stdin, a.container, a.blobname, a.big, a.contentSetting)
		check(err)
	} else {
		f, err := os.Open(a.filename)
		check(err)
		defer f.Close()
		if a.contentSetting.ContentType == "" {
			a.contentSetting.ContentType = mime.TypeByExtension(filepath.Ext(a.filename))
		}
		if a.blobname == "" {
			a.blobname = filepath.Base(a.filename)
		}

		if a.verbose {
			loggerPtr.Printf("Upload %s as ContentType %s\n", a.blobname, a.contentSetting.ContentType)
		}
		a.big, err = FileIsTooBig(a.filename)
		check(err)
		if a.verbose {
			loggerPtr.Printf("Filename: %s\n", a.filename)
			loggerPtr.Printf("Blobname: %s\n", a.blobname)
			loggerPtr.Printf("Big: %t\n", a.big)
		}
		err = s.SaveBlob(f, a.container, a.blobname, a.big, a.contentSetting)
		check(err)
	}
}

func download(a Arguments) {

	if a.verbose {
		loggerPtr.Printf("Account: %s\nContainer: %s\n", a.Account, a.container)

	}

	if a.pipe {
		if a.verbose {
			loggerPtr.Println("Output to Pipe")
		}
		check(s.LoadBlob(os.Stdout, a.container, a.blobname))
	} else {
		if a.verbose {
			loggerPtr.Println("Output to File")
		}
		if a.filename == "" {
			a.filename = a.blobname
		}
		f, err := os.Create(a.filename)
		check(err)
		defer f.Close()
		check(s.LoadBlob(f, a.container, a.blobname))
		f.Close()
	}
}

func containerList(a Arguments) {

	l, err := s.ListContainer(a.containerPrefix)
	check(err)
	for _, c := range l {
		fmt.Println(c)
	}
}

func containerCreate(a Arguments) {
	err := s.CreateContainer(a.container)
	check(err)
}

func containerDelete(a Arguments) {
	err := s.DeleteContainer(a.container)
	check(err)
}

func containerShow(a Arguments) {
	l, err := s.ShowContainer(a.container, a.blobPrefix)
	check(err)
	for _, b := range l {
		fmt.Println(b)
	}
}

func blobDelete(a Arguments) {
	check(s.DeleteBlob(a.container, a.blobname))
}
func printCmds() {
	fmt.Println("usage: abc <command> [<args>]")
	fmt.Println("The most commonly used abc commands are: ")
	fmt.Println(" upload              upload a file")
	fmt.Println(" download            download a file")
	fmt.Println(" container list      list all containers")
	fmt.Println(" container show      shows all blobs of a container")
	fmt.Println(" container create    create a container")
	fmt.Println(" container delete    delete a container")
	fmt.Println(" blob delete         delete a blob in container")
	fmt.Println(" save credentials    save azure credential to file")
}
func parseFlags(cmd string, args []string) (Arguments, error) {
	var a Arguments
	u, err := user.Current()
	if err != nil {
		return a, err
	}
	home := u.HomeDir
	//home := os.Getenv("HOME")
	a.configPath = fmt.Sprintf("%s/.azure/abc-config.json", home)
	subflags := flag.NewFlagSet("subcommand", flag.ExitOnError)

	subflags.BoolVar(&a.verbose, "v", false, "Verbose info")
	subflags.StringVar(&a.Key, "k", "", "a Azure Key (if stored - mandatory)")
	subflags.StringVar(&a.Account, "a", "", "a Azure Storage Account (if stored - mandatory)")

	switch cmd {
	case "upload":
		subflags.StringVar(&a.blobname, "n", "", "The Blob File Name (required for pipe)")
		subflags.IntVar(&a.workercount, "worker", 10, "download Worker Count")
		subflags.BoolVar(&a.pipe, "pipe", false, "incoming Pipe")
		subflags.BoolVar(&a.big, "big", false, "spilt file which are bigger than 195GB in part Blockblobs")
		subflags.StringVar(&a.contentSetting.ContentType, "contentType", "", "Contenttype for the uploaded file")
		subflags.StringVar(&a.contentSetting.CacheControl, "cacheControl", "", "CacheControl for the uploaded file")
		subflags.StringVar(&a.contentSetting.ContentLanguage, "contentLanguage", "", "ContentLanguage for the uploaded file")
		subflags.StringVar(&a.contentSetting.ContentEncoding, "contentEncoding", "", "ContentEncoding for the uploaded file")
		subflags.StringVar(&a.filename, "f", "", "Filename to upload (required if no pipe)")
		subflags.StringVar(&a.container, "c", "", "a Azure Container (required)")
	case "download":
		subflags.StringVar(&a.blobname, "n", "", "The Blob File Name (required)")
		subflags.IntVar(&a.workercount, "worker", 10, "download Worker Count")
		subflags.BoolVar(&a.pipe, "pipe", false, "outgoing Pipe")
		subflags.StringVar(&a.filename, "f", "", "Filename to download")
		subflags.StringVar(&a.container, "c", "", "a Azure Container (required)")
	case "container_create", "container_delete":
		subflags.StringVar(&a.container, "c", "", "a Azure Container (required)")
	case "container_show":
		subflags.StringVar(&a.container, "c", "", "a Azure Container (required)")
		subflags.StringVar(&a.blobPrefix, "bp", "", "a Azure Blob Prefix")
	case "container_list":
		subflags.StringVar(&a.containerPrefix, "cp", "", "a Azure Container Prefix")
	case "blob_delete":
		subflags.StringVar(&a.container, "c", "", "a Azure Container (required)")
		subflags.StringVar(&a.blobname, "n", "", "The Blob File Name (required)")
	case "save_credentials":
		// no special Flags needed
	default:
		return a, errors.New("No valid Command")
	}
	a.load()
	subflags.Parse(args)
	if len(subflags.Args()) > 0 || len(args) == 0 {
		subflags.PrintDefaults()
		os.Exit(2)
		//return a, errors.New("No valid Flags")
	}
	return a, nil
}

func main() {

	loggerPtr = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	if len(os.Args) == 1 {
		printCmds()
		return
	}
	var cmd string
	var args []string
	switch os.Args[1] {
	case "upload", "download":
		args = os.Args[2:]
		cmd = os.Args[1]
	case "container", "blob", "save":
		args = os.Args[3:]
		cmd = fmt.Sprintf("%s_%s", os.Args[1], os.Args[2])
	default:
		printCmds()
		fmt.Println("")
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
	arguments, err := parseFlags(cmd, args)
	if err != nil {
		printCmds()
		fmt.Println("")
		fmt.Printf("%q is not valid command.\n", strings.Replace(cmd, "_", " ", 1))
		os.Exit(2)
	}

	s = storage.StorageAttributes{
		Key:         arguments.Key,
		Account:     arguments.Account,
		WorkerCount: arguments.workercount,
		Logger:      loggerPtr,
		Verbose:     arguments.verbose}
	s.NewStorageClient()

	switch cmd {
	case "upload":
		upload(arguments)
	case "download":
		download(arguments)
	case "save_credentials":
		arguments.save()
	case "container_list":
		containerList(arguments)
	case "container_create":
		containerCreate(arguments)
	case "container_delete":
		containerDelete(arguments)
	case "container_show":
		containerShow(arguments)
	case "blob_delete":
		blobDelete(arguments)
	}
}
