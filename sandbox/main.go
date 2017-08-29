package main

import (
	"errors"
	"fmt"
)

/*
	"fmt"
	"math/rand"
	"time"
*/

/*
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
*/

/*
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
*/

var (
	//ErrSrcNotExist src file doesnt exist
	ErrSrcNotExist = errors.New("Source file does not exist")

	//ErrSrcNotRegularFile src file is not a regular file
	ErrSrcNotRegularFile = errors.New("Source file is not a regular file")

	//ErrDstNotRegularFile dst file is not a regular file
	ErrDstNotRegularFile = errors.New("Destination file is not a regular file")
)

/*
func commandOutput(cmdLine string) (string, error) {
	fmt.Println("commandOutput ENTER")
	fmt.Println("Cmdline:", cmdLine)

	cmd := exec.Command("bash", "-c", cmdLine)
	if cmd == nil {
		fmt.Println("Error creating cmd")
		fmt.Println("commandOutput LEAVE")
		return "", errors.New("Cmd is nil")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error getting output:", err)
		fmt.Println("commandOutput LEAVE")
		return "", err
	}

	output := strings.TrimSpace(string(out))

	fmt.Println("commandOutput Succeeded")
	fmt.Println(output)
	fmt.Println("commandOutput LEAVE")
	return output, nil
}

func regexMatch(haystack string, regex string) ([]string, error) {
	fmt.Println("RegexMatch ENTER")
	fmt.Println("haystack:", haystack)
	fmt.Println("regexp:", regex)

	r, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("Rexexp Failed:", err)
		fmt.Println("RegexMatch LEAVE")
		return nil, err
	}

	strings := r.FindStringSubmatch(haystack)
	if strings == nil {
		fmt.Println("Rexexp Failed:", err)
		fmt.Println("RegexMatch LEAVE")
		return nil, errors.New("Regex failed")
	}

	fmt.Println("RegexMatch:", strings)
	fmt.Println("RegexMatch LEAVE")
	return strings, nil
}
*/

/*
func doesDependencyExist(serviceName string, depName string) (bool, error) {
	fmt.Println("doesDependencyExist ENTER")
	fmt.Println("serviceName:", serviceName)
	fmt.Println("depName:", depName)

	fileName := "/etc/init.d/" + serviceName
	fmt.Println("fileName:", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Failed on file Open:", err)
		fmt.Println("doesDependencyExist LEAVE")
		return false, err
	}
	defer file.Close()

	r, err := regexp.Compile(depName)
	if err != nil {
		fmt.Println("regexp is invalid")
		fmt.Println("doesDependencyExist LEAVE")
		return false, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fmt.Println("Line:", line)
		if len(line) == 0 {
			continue
		}

		strings := r.FindStringSubmatch(line)
		if strings != nil || len(strings) == 1 {
			fmt.Println("Match found:", line)
			fmt.Println("doesDependencyExist LEAVE")
			return true, nil
		}
	}

	fmt.Println("Dependency was not found")
	fmt.Println("doesDependencyExist LEAVE")

	return false, nil
}

func makeTmpFileWithNewDep(serviceName string, depName string) error {
	fmt.Println("makeTmpFileWithNewDep ENTER")
	fmt.Println("serviceName:", serviceName)
	fmt.Println("depName:", depName)

	fileName := "/etc/init.d/" + serviceName
	fmt.Println("fileName:", fileName)

	fileNameTmp := "/tmp/" + serviceName + ".tmp"
	fmt.Println("fileNameTmp:", fileNameTmp)

	sfi, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Src Stat Failed:", err)
		fmt.Println("makeTmpFileWithNewDep LEAVE")
		return ErrSrcNotExist
	}
	if !sfi.Mode().IsRegular() {
		//cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		fmt.Println("Src file is not regular")
		fmt.Println("makeTmpFileWithNewDep LEAVE")
		return ErrSrcNotRegularFile
	}
	dfi, err := os.Stat(fileNameTmp)
	if err == nil {
		if !(dfi.Mode().IsRegular()) {
			fmt.Println("Dst file is not regular")
			fmt.Println("makeTmpFileWithNewDep LEAVE")
			return ErrDstNotRegularFile
		}
	}

	//Copy the file
	in, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open SRC file:", err)
		fmt.Println("makeTmpFileWithNewDep LEAVE")
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(fileNameTmp, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Failed to open DST file:", err)
		fmt.Println("makeTmpFileWithNewDep LEAVE")
		return err
	}
	defer out.Close()

	r, err := regexp.Compile("Required-Start:")
	if err != nil {
		fmt.Println("regexp is invalid")
		fmt.Println("makeTmpFileWithNewDep LEAVE")
		return err
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fmt.Println("Line:", line)
		if len(line) == 0 {
			continue
		}

		str := r.FindStringSubmatch(line)
		if str != nil {
			fmt.Println("Match found:", line)
			newLine := line + " scini"
			out.WriteString(newLine + "\n")
		} else {
			out.WriteString(line + "\n")
		}
	}

	err = out.Sync()
	if err != nil {
		fmt.Println("Failed to flush file:", err)
		fmt.Println("makeTmpFileWithNewDep LEAVE")
		return err
	}

	fmt.Println("makeTmpFileWithNewDep Succeeded")
	fmt.Println("makeTmpFileWithNewDep LEAVE")

	return nil
}

func copyFileEx(src string, dst string, mode os.FileMode) error {
	fmt.Println("CopyFile ENTER")
	fmt.Println("SRC:", src)
	fmt.Println("DST:", dst)

	sfi, err := os.Stat(src)
	if err != nil {
		fmt.Println("Src Stat Failed:", err)
		fmt.Println("CopyFile LEAVE")
		return ErrSrcNotExist
	}
	if !sfi.Mode().IsRegular() {
		//cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		fmt.Println("Src file is not regular")
		fmt.Println("CopyFile LEAVE")
		return ErrSrcNotRegularFile
	}
	dfi, err := os.Stat(dst)
	if err == nil {
		if !(dfi.Mode().IsRegular()) {
			fmt.Println("Dst file is not regular")
			fmt.Println("CopyFile LEAVE")
			return ErrDstNotRegularFile
		}
		if os.SameFile(sfi, dfi) {
			fmt.Println("Src and Dst files are the same")
			fmt.Println("CopyFile LEAVE")
			return nil
		}
	}

	//Copy the file
	in, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open SRC file:", err)
		fmt.Println("CopyFile LEAVE")
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, mode)
	if err != nil {
		fmt.Println("Failed to open DST file:", err)
		fmt.Println("CopyFile LEAVE")
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		fmt.Println("Failed to copy file:", err)
		fmt.Println("CopyFile LEAVE")
		return err
	}

	err = out.Sync()
	if err != nil {
		fmt.Println("Failed to flush file:", err)
		fmt.Println("CopyFile LEAVE")
		return err
	}

	fmt.Println("CopyFile succeeded")
	fmt.Println("CopyFile LEAVE")
	return nil
}
*/

func main() {
	/*
		test := []string{"string1", "string2", "string3"}

		rand.Seed(time.Now().UnixNano())

		for i := 0; i < 10; i++ {
			random := rand.Intn(len(test))
			fmt.Println(test[random])
		}
	*/

	/*
		test := "ping 1"
		iindex := strings.Index(test, " ")
		value := test[iindex+1:]
		fmt.Println("value: ", value)
	*/

	/*
		for i := 0; i < len(os.Args); i++ {
			fmt.Println("Arg", (i + 1), ":", os.Args[i])
		}

		fmt.Println("hello")
		fmt.Println("how")
		fmt.Println("are")
		fmt.Println("you")
	*/

	/*
		outputCmd := "fdisk -l | grep \\/dev\\/"
		output, err := commandOutput(outputCmd)
		if err != nil {
			fmt.Println("Failed to get device list. Err:", err)
			return
		}

		buffer := bytes.NewBufferString(output)
		for {
			str, err := buffer.ReadString('\n')

			needles, errRegex := regexMatch(str, "Disk (/dev/.*):")
			if errRegex != nil {
				fmt.Println("RegexMatch Failed. Err:", err)
				if err == io.EOF {
					break
				}
				continue
			}
			if len(needles) != 2 {
				fmt.Println("Incorrect size")
				if err == io.EOF {
					break
				}
				continue
			}
			device := needles[1]
			fmt.Println("Device Found:", device)
			if err == io.EOF {
				break
			}
		}
	*/

	/*
		found, err := doesDependencyExist("rexray", "scini")
		if err != nil {
			fmt.Println("doesDependencyExist Failed. Err:", err)
			return
		}
		if found {
			fmt.Println("Dependency already exists!")
			return
		}

		err = makeTmpFileWithNewDep("rexray", "scini")
		if err != nil {
			fmt.Println("makeTmpFileWithNewDep Failed. Err:", err)
			return
		}

		err = copyFileEx("/tmp/rexray.tmp", "/etc/init.d/rexray", 0666)
		if err != nil {
			fmt.Println("CopyFile Failed. Err:", err)
			return
		}

		fmt.Println("AddDependentService Succeeded")
	*/

	/*
		device := "/dev/xvdg"
		lenDev := len(device) - 1
		fmt.Println("lenDev:", lenDev)

		highestLetter := "f"
		if highestLetter[0] < device[lenDev] {
			highestLetter = device[lenDev:]
		}
		highestLetter = string(highestLetter[0] + 1)
		fmt.Println("highestLetter:", highestLetter)

		newDevice := "/dev/xvd" + highestLetter
		fmt.Println("newDevice:", newDevice)
	*/

	cnt := uint(1)
	fmt.Println("Value:", (cnt%1 == 0))

	cnt = cnt + 1
	fmt.Println("Value:", (cnt%1 == 0))
}
