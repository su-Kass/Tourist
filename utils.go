package main

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"runtime"
	"path/filepath"
	"strconv"
)

type UserData struct {
	HmeDir	     string
	CurDir	     string
	TmpDir	     string
	AppDir	     string
	SysType	     string
	CpuNum	     string
	Username     string
}

func ReturnArguments() (flag string, data string) {
	if len(os.Args) > 1 {
		flag = os.Args[1]
		MakeOutput("2", flag+"\n")
	}
	if len(os.Args) > 2 {
		data = os.Args[2]
		MakeOutput("2", data+"\n")
	}
	return flag, data
}


func ErrHandle(err error, exit bool) bool {
	var errBool bool = false
	if err != nil {
		if exit {
			panic(err)
		}
		errBool = true
	}
	return errBool
}

// I didn't wanna import fmt for a singular function + allows me to format things; prints text
//
// 1=red 2=green 3=yellow 4=blue, etc
func MakeOutput(color string, message string) {
	var code string = "\033[3" + color + "m"
	var out string = code + message + "\033[0m"
	os.Stdout.Write([]byte(out))
}

func TalkCat(text string) {
	MakeOutput("4", "/ᐠ˵- ⩊ -˵マ: ")
	MakeOutput("3", text+"\n")
}

func SearchForFiletype(folder, filetype string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(folder, "*."+filetype))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no file matching *.%s in %s", filetype, folder)
	}
	return filepath.Base(matches[0]), nil
}

//Runs a command in the user's shell; hopefully used only to run system programs, however can be used to do anything
//
// Has command and then arguments:
//
// echo <--- command
//
// 'Hello World!' <--- argument
func RunShell(command string) {
	cmd := exec.Command("sh", "-c", command)

	if err := cmd.Run(); err != nil {
		MakeOutput("1", err.Error()+"\n")
	}

	MakeOutput("2", "/ᐠ˵- ⩊ -˵マ : ")
	MakeOutput("3", "Run, "+command)
	MakeOutput("2", "(˵◝ ⩊  ◜˵マ \n")
}

// Simply makes a directory, dont question why it exists.
func MakeDir(Dir string) {
	os.Mkdir(Dir, 0755)
}

func MoveFile(src, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("could not move %s to %s: %w", src, dst, err)
	}
	return nil
}

// Copies a file somewhere else
//
// File, Directory
func CopyFile(file string, dir string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return os.WriteFile(dir, data, 0755)
}

// Makes a file at the directory
//
// Returns an error
func MakeFile(fileName string) error {
	_, err := os.Create(fileName)
	return err
}

// Writes a line for a given file
//
// Returns error
func WriteLn(fileName string, line int, text string) error{
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")

	for len(lines) < line {
		lines = append(lines, "")
	}
	lines[line-1] = text
	return os.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

// Reads the given line for the given file.ReadLn
//
// Returns data + success bool
func ReadLn(fileName string, line int) (string, bool){
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", false
	}
	lines := strings.Split(string(data), "\n")

	if line < 1 || line > len(lines) {
		return "", false
	}
	return lines[line-1], true
}

// ReturnData returns a map of the user's data needed for setting up stuff
//
// returns map: HOME CURRENT SYSTYPE TEMPDIR CPUNUM USERNAME
func ReturnData() UserData {
	var userData UserData
	var homeDir, _ = os.UserHomeDir()
	var currentDir, _ = os.Getwd()

	var systemType string = runtime.GOOS
	var tempFolder string = os.TempDir()
	var cpuNum string = strconv.Itoa(runtime.NumCPU())

	var userName string
	var applicationDir string
	switch runtime.GOOS {
		case "linux":
			applicationDir = filepath.Join(homeDir, ".local", "share")
			userName = os.Getenv("USER")
		case "windows":
			applicationDir = filepath.Join(os.Getenv("LOCALAPPDATA"))
			userName = os.Getenv("USERNAME")
		default: applicationDir = "UNKNOWN!"
	}
	userData.HmeDir = homeDir
	userData.CurDir = currentDir
	userData.SysType = systemType
	userData.TmpDir = tempFolder
	userData.CpuNum = cpuNum
	userData.Username = userName
	userData.AppDir = applicationDir
	return userData
}
