package SaveFile

import (
	"archive/zip"
	"fmt"
)

func ProcessSave(filename string) (*zip.ReadCloser, error) {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return archive, nil
}
