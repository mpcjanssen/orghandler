package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kardianos/osext"
)

const taskCreatedPrefix = "Created task "

func main() {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	logFile := filepath.Join(folderPath, "orghandler.log")
	f, err := os.Create(logFile)
	if err != nil {
		os.Exit(-1)
	}
	defer f.Close()
	logger := log.New(f, "", log.Ldate)
	orgURL, err := url.Parse(os.Args[1])
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("url: " + os.Args[1])

	title := orgURL.Query().Get("title")
	body := strings.TrimSpace(orgURL.Query().Get("body"))
	source := orgURL.Query().Get("url")

	logger.Println(source)
	logger.Println(title)
	logger.Println("")
	logger.Println(body)
	taskwarriorAdd(logger, f, title, source, body)
	//nextcloudAdd(logger, f, title, source, body)
}

//
//func nextcloudAdd(logger *log.Logger, f *os.File,  title string,  source string, body string) {
//	ncFile, err := os.OpenFile("/Users/mark.janssen/Nextcloud/todo/todo.txt", os.O_APPEND|os.O_WRONLY, 0644)
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer ncFile.Close()
//	date := time.Now().Format("2006-01-02")
//	fmt.Fprintln(ncFile,date, title, source, "+orgcapture" )
//	logger.Println("Task created")
//}

const taskPath = "/usr/local/bin/task"

func taskwarriorAdd(logger *log.Logger, f *os.File, title string, source string, body string) {
	cmd := exec.Command(taskPath, "add", title)
	cmd.Stderr = f
	out, err := cmd.Output()
	if err != nil {
		logger.Fatal(err)
	}
	strOut := string(out)
	logger.Println(string(out))
	if !strings.HasPrefix(strOut, taskCreatedPrefix) {
		logger.Fatal("unexpected result")
	}
	strOut = strings.TrimPrefix(strOut, taskCreatedPrefix)
	taskID := strings.TrimRight(strOut, ".\n")
	_, err = strconv.Atoi(taskID)
	if err != nil {
		logger.Fatal(err)
	}
	// Add url and body as annotation
	annotateTask(logger, taskID, source)
	if body != "" {
		annotateTask(logger, taskID, body)
	}
	cmd = exec.Command(taskPath, taskID, "sync")
	result, _ := cmd.Output()
	cmd.Wait()
	logger.Println(string(result))
}

func annotateTask(logger *log.Logger, taskID string, annotation string) {
	cmd := exec.Command(taskPath, taskID, "annotate", annotation)
	result, err := cmd.Output()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(string(result))
}
