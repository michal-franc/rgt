### rgt -  Red Green Test
Auto test runner that monitor files and run tests whenever they are changed 

### Motivation
I am just lazy and like to start tests asap whenever I change code to verify instantaneously if what I am doing makes sense.

#### Future
Support for more types of files and test runners so that if I change python a python test runner starts etc.

### Current limitations
- very basic and not optimized
- only `GO` supported at the moment
- watches all the files - even `.git` file changes

### Usage

Run `rgt` in your main folder to watch all the files in sub folders
```
rgt
```

Commands
```
rgt start
      --go-test-runner string   Specifies which go test runner to use. (default "go", supports gotestsum)
      --sub-folder-only         If set true will run only tests from the folder the file that is changed in.
```

### Installation
It was tested on `linux`.  
Download binary from github releases and unpack it to your `PATH` folder.
Using Go Get
```
go get -u github.com/michal-franc/rgt/cmd/rgt
```
