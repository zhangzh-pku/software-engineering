package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	_, err = os.Stat(dst)
	if err == nil {
		err = os.Remove(dst)
		if err != nil {
			return fmt.Errorf("files not removed: %v", err)
		}
	}
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	return nil
}

func CopyDir(src, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %q: %v\n", path, err)
		}

		dstPath := filepath.Join(dst, path[len(src):])
		if info.IsDir() {
			err := os.MkdirAll(dstPath, info.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		} else {
			err := CopyFile(path, dstPath)
			if err != nil {
				return fmt.Errorf("failed to copy file: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %v: %v", src, err)
	}

	return nil
}

func CopyFiles(sourceDir, destinationDir string) error {
	fileInfos, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %v", err)
	}

	for _, fileInfo := range fileInfos {
		sourcePath := filepath.Join(sourceDir, fileInfo.Name())
		destinationPath := filepath.Join(destinationDir, fileInfo.Name())

		err := CopyFile(sourcePath, destinationPath)
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %v", sourcePath, err)
		}
	}

	return nil
}
