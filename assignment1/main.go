package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	homeDir, err := getHomeDir()
	if err != nil {
		log.Fatal("Couldn't get home dir")
	}

	rootchPath := homeDir + "/rootch"
	binPath := rootchPath + "/bin"

	libDirPath := rootchPath + "/lib/x86_64-linux-gnu"
	lib64DirPath := rootchPath + "/lib64"

	devDirPath := rootchPath + "/dev"

	x86_64libs := []string{"/libtinfo.so.5", "/libdl.so.2", "/libc.so.6"}
	lib64 := "/ld-linux-x86-64.so.2"

	log.Printf("Creating chroot file environment under the '%s'\n", rootchPath)

	/* create home folder for jail */
	if _, err := os.Stat(rootchPath); os.IsNotExist(err) {
		os.Mkdir(rootchPath, 0775)
	}

	/* create bin folder for jail*/
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0775)
	}

	/* create lib folders for jail */
	if _, err := os.Stat(libDirPath); os.IsNotExist(err) {
		os.MkdirAll(libDirPath, 0775)
	}

	if _, err := os.Stat(lib64DirPath); os.IsNotExist(err) {
		os.MkdirAll(lib64DirPath, 0775)
	}

	if _, err := os.Stat(devDirPath); os.IsNotExist(err) {
		os.MkdirAll(devDirPath, 0775)
	}

	/* copy /bin/bash to jail */
	copyCmd := exec.Command("cp", "/bin/bash", binPath)
	err = copyCmd.Run()
	if err != nil {
		log.Fatalf("Error while copying bash: %s", err)
	}

	/* copy lib files */
	for _, value := range x86_64libs {
		copyLibsCmd := exec.Command("cp", "/lib/x86_64-linux-gnu"+value, libDirPath+value)
		err = copyLibsCmd.Run()
		if err != nil {
			log.Fatalf("Error while copying lib files: %s", err)
		}
	}

	copyLib64Cmd := exec.Command("cp", "/lib64"+lib64, lib64DirPath+lib64)
	err = copyLib64Cmd.Run()
	if err != nil {
		log.Fatalf("Error while copying lib64 file: %s", err)
	}

	copyDevCmd := exec.Command("cp", "/dev/null", devDirPath+"/null")
	err = copyDevCmd.Run()
	if err != nil {
		log.Fatalf("Error while copying /dev/null file: %s")
	}

	log.Println("Successfully created chroot environment.")

	err = unix.Chroot(rootchPath)
	if err != nil {
		log.Fatalf("Error while executing chroot: %s", err)
	}

	fmt.Println("Listing all files under the chroot:")
	listFiles("/")

	err = os.Chdir("/bin")
	if err != nil {
		log.Fatalf("Error while changing dir: %s", err)
	}

	/*
		isChrootCmd := exec.Command("./bash", "./ischroot")
		err = isChrootCmd.Run()
		if err != nil {
			log.Fatalf("Error while running ischroot.sh: %s", err)
		}*/

	fmt.Println("is Chrooted: ", isChrooted())

	os.Exit(0)
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func listFiles(dir string) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}

/* https://golang.org/src/syscall/exec_linux_test.go
 * Check if we are in a chroot by checking if the inode of / is
 * different from 2 (there is no better test available to non-root on
 * linux */
func isChrooted() bool {
	root, err := os.Stat("/")
	if err != nil {
		log.Fatalf("cannot stat /: %v", err)
	}

	return root.Sys().(*syscall.Stat_t).Ino != 2
}
