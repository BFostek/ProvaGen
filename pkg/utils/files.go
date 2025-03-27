package utils

import (
	"fmt"
	"os"
	"path"
)

func CreateStructure(dir, slug string) error {
	createDirectoryIfNotExists(dir)
	createDirectoryIfNotExists(path.Join(dir, slug))
	return nil

}
func createDirectoryIfNotExists(dir string) error {
	println(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		fmt.Printf("Created directory: %s\n", dir)
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking directory %s: %v", dir, err)
	}
	return nil
}

func CreateFileWithContent(fpath string, content string) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fpath, err)
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", fpath, err)
	}

	fmt.Printf("Created file: %s\n", fpath)
	return nil
}
