### rgt -  Red Green Test
Auto test runner that monitor files and run tests whenever they are changed

### Motivation
I am just lazy and like to start tests asap whenever I change code to verify instantaneously if what I am doing makes sense.

### How it works

1. You run `rgt` (or `rgt start`) in your project directory
2. It starts watching all files recursively using fsnotify
3. When you save any file, it automatically runs your tests
4. Shows a spinner while tests run
5. Displays the test results with colors preserved

**Basic usage:**
```bash
cd your-project
rgt                    # Just this! Defaults to 'start' command
```

**With options:**
```bash
rgt --sub-folder-only          # Only run tests in the changed file's subfolder
rgt --test-name TestFoo        # Only run specific test
rgt --test-runner gotestsum    # Use gotestsum instead of go test
rgt --test-type python         # For Python projects with pytest
```

The "Red Green" in the name is a reference to the TDD workflow: Red (failing test) → Green (passing test) → Refactor.


#### Future
Support for more types of files and test runners so that if I change python a python test runner starts etc.

### Current limitations
- very basic and not optimized
- only `GO`, `Python` supported at the moment
- watches all the files - even `.git` file changes

### Installation
It was tested on `linux`.  
Download binary from github releases and unpack it to your `PATH` folder.
Using Go Get
```
go install github.com/michal-franc/rgt/cmd/rgt@latest
```
