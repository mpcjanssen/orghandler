package main

import (
	"path/filepath"
	"os"
	"log"
	"github.com/kardianos/osext"
	"net/url"
	"os/exec"
	"strings"
	"strconv"
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
	orgURL,err := url.Parse(os.Args[1])
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

	cmd := exec.Command("task", "add", title)
	cmd.Stderr = f
	out, err :=cmd.Output()
	if err != nil {
		logger.Fatal(err)
	}
	strOut := string(out)
	logger.Println(string(out))
	if !strings.HasPrefix(strOut, taskCreatedPrefix) {
		logger.Fatal("unexpected result")
	}
	strOut = strings.TrimPrefix(strOut , taskCreatedPrefix)
	taskId := strings.TrimRight(strOut, ".\n")
	_, err = strconv.Atoi(taskId)
	if err != nil {
		logger.Fatal(err)
	}
	// Add url and body as annotation
	annotateTask(logger, taskId, source)

	if body != "" {
		annotateTask(logger, taskId, body)
	}
	cmd = exec.Command("task", taskId, "sync")
	result, _ := cmd.Output()
	cmd.Wait()
	logger.Println(string(result))

}

func annotateTask(logger *log.Logger, taskId string, annotation string)  {
	cmd := exec.Command("task", taskId, "annotate", annotation)
	result, err := cmd.Output()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(string(result))
}