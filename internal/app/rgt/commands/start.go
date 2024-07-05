package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/creack/pty"
	"github.com/fsnotify/fsnotify"
	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

//TODO: run only changed test files
//TODO: migrate to gotestsum --watch ??
//TODO: add proper colors in the out
//TODO: discover `files` and extension and based on file change run different test runner - i have two languages in the project pythong golang when i change go code i get golang test run when i change python i get python runner
//TODO: First run should happen at start - we dont need to wait for file change
//TODO: BUGFIX when there is no library and io cant run go test ./... then the screen is empty

//TODO: global configuration and per project configuration

//TODO: todo run only changed files?
//TODO: collect multiple file changes -> filter out by GO -> remove duplicates either files or folders
//TODO: support other languages, python, rust, bash, go, lua

//TODO add viper and cobra

var watcher *fsnotify.Watcher
var testRunner string
var testNames string
var runTestsIntheSubFolder bool
var testType string
var supportedTypes = [...]string{"golang", "python"}

func init() {
	startCmd.Flags().StringVar(&testRunner, "test-runner", "default", "Specifies which test runner to use.")
	startCmd.Flags().StringVar(&testType, "test-type", "golang", fmt.Sprintf("Specifies which test runner to run supported. %s", supportedTypes))
	startCmd.Flags().StringVar(&testNames, "test-name", "", fmt.Sprintf("Language/tool specific value to filter out tests to run. %s", supportedTypes))
	startCmd.Flags().BoolVar(&runTestsIntheSubFolder, "sub-folder-only", false, "If set true will run only tests from the folder the file that is changed is.")
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start auto test runner",
	Run: func(cmd *cobra.Command, args []string) {
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		//TODO: ability to override the behaviour and just run all the tests
		//TODO: ignore vendor + .gitignore folders
		//TODO: only add folders to watch if they have `go` file in there
		//TODO: handle exit of the app

		// go through each subfolder and add it to watcher
		if err := filepath.Walk(".", watchDir); err != nil {
			fmt.Println("ERROR", err)
		}

		done := make(chan bool)
		lastFileWritten := ""
		goFuncStarted := false
		fmt.Printf("Started rgt using `%s` test runner.\n", testRunner)
		if runTestsIntheSubFolder {
			fmt.Print("sub-folder-only mode\n")
		}
		fmt.Println("Waiting for file changes.")
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						//TODO: only watch on golang files and configure ability to ignore files
						lastFileWritten = event.Name
						// We want to start goroutine block it for 1-2-3 seconds
						// And then listen to files - there might be many events for 1 file
						// We are only interested in one last one
						// so we keep collecting lastFileWritten events but start only one goroutine
						// its waiting with pointer to lastFileWritten so that its dynamic but goroutine will wait
						if !goFuncStarted {
							go func(fPath *string) {
								goFuncStarted = true

								screen.Clear()
								screen.MoveTopLeft()
								s := spinner.New(spinner.CharSets[35], 100*time.Millisecond) // Build our new spinner
								s.Color("red", "bold")
								s.Start()

								var cmd *exec.Cmd

								if testType == "golang" {
									//TODO: support for any test runner from config like gotestsum - hacky mess at the moment
									if testRunner == "default" {
										if runTestsIntheSubFolder {
											cmd = exec.Command("go", "test", "-run="+testNames)
											cmd.Dir = extractDir(lastFileWritten)
										} else {
											cmd = exec.Command("go", "test", "./...", "-run="+testNames)
										}
									} else if testRunner == "gotestsum" {
										cmd = exec.Command(testRunner)
										if runTestsIntheSubFolder {
											cmd.Dir = extractDir(lastFileWritten)
										}
									}
								} else {
									if testRunner == "default" {
										cmd = exec.Command("pytest", fmt.Sprintf("./%s", *fPath))
									} else {
										cmd = exec.Command(testRunner, fmt.Sprintf("./%s", *fPath))
									}
								}

								if cmd == nil {
									log.Fatalf("incorrect test runner specified '%s' supported [go, gotestsum, pytest, python]", testRunner)
								}

								// Make sure we get the error
								ptmx, err := pty.Start(cmd)
								if err != nil {
									fmt.Printf("Failed to start command: %s\n", err)
								}
								defer func() { _ = ptmx.Close() }()

								var buf bytes.Buffer

								_, _ = io.Copy(&buf, ptmx)
								s.Stop()

								fmt.Println(buf.String())
								goFuncStarted = false
							}(&lastFileWritten)
						}
					}

					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()

		<-done
	},
}

func extractDir(fPath string) string {
	return filepath.Dir(fPath)
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}
