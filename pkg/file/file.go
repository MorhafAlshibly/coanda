package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindMyRootDir() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	lastDir := workingDirectory
	myUniqueRelativePath := ".git"
	for {
		currentPath := fmt.Sprintf("%s/%s", lastDir, myUniqueRelativePath)
		fi, err := os.Stat(currentPath)
		if err == nil {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				return strings.Replace(currentPath, myUniqueRelativePath, "", 1)
			}
		}
		newDir := filepath.Dir(lastDir)
		if newDir == "/" || newDir == lastDir {
			return ""
		}
		lastDir = newDir
	}
}
