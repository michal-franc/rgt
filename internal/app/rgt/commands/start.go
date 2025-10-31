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

// Map of test types to their file extensions
var testTypeExtensions = map[string][]string{
	"golang": {".go"},
	"python": {".py"},
}

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
		// If --test-type flag was not explicitly set, run interactive mode
		if !cmd.Flags().Changed("test-type") {
			detectAndPromptTestType()
		}

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

		runTests(lastFileWritten)
		fmt.Println("Waiting for file changes...")
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						// Filter by file extension based on test type
						if !shouldProcessFile(event.Name, testType) {
							continue
						}

						lastFileWritten = event.Name
						// We want to start goroutine block it for 1-2-3 seconds
						// And then listen to files - there might be many events for 1 file
						// We are only interested in one last one
						// so we keep collecting lastFileWritten events but start only one goroutine
						// its waiting with pointer to lastFileWritten so that its dynamic but goroutine will wait
						if !goFuncStarted {
							go func(fPath *string) {
								goFuncStarted = true

								runTests(*fPath)
								fmt.Println("Waiting for file changes...")

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

func runTests(fPath string) {
	screen.Clear()
	screen.MoveTopLeft()

	if fPath != "" {
		fmt.Printf("File changed: %s\n\n", fPath)
	}

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond) // Build our new spinner
	s.Color("red", "bold")
	s.Start()

	var cmd *exec.Cmd

	if testType == "golang" {
		//TODO: support for any test runner from config like gotestsum - hacky mess at the moment
		if testRunner == "default" {
			if runTestsIntheSubFolder {
				cmd = exec.Command("go", "test", "-run="+testNames)
				cmd.Dir = extractDir(fPath)
			} else {
				cmd = exec.Command("go", "test", "./...", "-run="+testNames)
			}
		} else if testRunner == "gotestsum" {
			cmd = exec.Command(testRunner)
			if runTestsIntheSubFolder {
				cmd.Dir = extractDir(fPath)
			}
		}
	} else {
		if testRunner == "default" {
			cmd = exec.Command("pytest", fmt.Sprintf("./%s", fPath))
		} else {
			cmd = exec.Command(testRunner, fmt.Sprintf("./%s", fPath))
		}
	}

	if cmd == nil {
		log.Fatalf("incorrect test runner specified '%s' supported [go test, gotestsum, pytest]", testRunner)
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

// shouldProcessFile checks if a file should trigger tests based on its extension and test type
func shouldProcessFile(filePath string, testType string) bool {
	fileExt := filepath.Ext(filePath)

	// Handle "all" test type - check against all known extensions
	if testType == "all" {
		for _, extensions := range testTypeExtensions {
			for _, ext := range extensions {
				if fileExt == ext {
					return true
				}
			}
		}
		return false
	}

	extensions, ok := testTypeExtensions[testType]

	// If test type not found, process all files (backward compatible)
	if !ok {
		return true
	}

	for _, ext := range extensions {
		if fileExt == ext {
			return true
		}
	}
	return false
}

// detectProjectFileTypes scans the current directory for .go and .py files
func detectProjectFileTypes() map[string]bool {
	found := map[string]bool{
		"golang": false,
		"python": false,
	}

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories (but not current dir) and common ignore paths
		if info.IsDir() && path != "." {
			if (len(info.Name()) > 0 && info.Name()[0] == '.') || info.Name() == "vendor" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
		}

		ext := filepath.Ext(path)
		if ext == ".go" {
			found["golang"] = true
		} else if ext == ".py" {
			found["python"] = true
		}

		// Early exit if both types found
		if found["golang"] && found["python"] {
			return filepath.SkipAll
		}

		return nil
	})

	return found
}

// detectAndPromptTestType detects available file types and prompts user for selection
func detectAndPromptTestType() {
	fmt.Println("Detecting project files...")
	detected := detectProjectFileTypes()

	foundTypes := []string{}
	if detected["golang"] {
		foundTypes = append(foundTypes, "Golang")
	}
	if detected["python"] {
		foundTypes = append(foundTypes, "Python")
	}

	// No test files found
	if len(foundTypes) == 0 {
		fmt.Println("❌ No Go or Python files found in project.")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	// Only one type found - auto-select but confirm
	if len(foundTypes) == 1 {
		typeFound := "golang"
		if detected["python"] {
			typeFound = "python"
		}
		fmt.Printf("✓ Project checked: found %s files only\n", foundTypes[0])
		fmt.Printf("Auto-selecting %s test runner.\n", foundTypes[0])
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		testType = typeFound
		return
	}

	// Both types found - show menu
	fmt.Printf("✓ Project checked: found %s and %s files\n", foundTypes[0], foundTypes[1])
	fmt.Println("\nPlease select test runner:")
	fmt.Println("  1) Golang")
	fmt.Println("  2) Python")
	fmt.Println("  3) All (watch both)")
	fmt.Print("\nYour choice (1-3): ")

	var choice int
	for {
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Scanf("%s", new(string)) // Clear input buffer
			fmt.Print("Invalid input. Please enter 1, 2, or 3: ")
			continue
		}

		switch choice {
		case 1:
			testType = "golang"
			fmt.Println("→ Selected: Golang")
			return
		case 2:
			testType = "python"
			fmt.Println("→ Selected: Python")
			return
		case 3:
			testType = "all"
			fmt.Println("→ Selected: All (watching both Go and Python files)")
			return
		default:
			fmt.Print("Invalid choice. Please enter 1, 2, or 3: ")
		}
	}
}
