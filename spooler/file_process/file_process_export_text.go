package file_process

import (
	"axia4/config"
	"axia4/log"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func doExportText(filePath string, fileContentText string, overwrite bool) error {

	if config.File.Paths.FileExport == "" {
		return errConfigNoPathExport
	}
	if filePath == "" {
		return errPathEmpty
	}

	// define paths
	filePathTarget := filepath.Join(config.File.Paths.FileExport, filePath)

	log.Info(log.ContextFile, fmt.Sprintf("exporting text file to path '%s'", filePathTarget))

	if err := checkExportPath(filePathTarget, overwrite); err != nil {
		return err
	}

	file, err := os.Create(filePathTarget)
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	if _, err := bufWriter.WriteString(fileContentText); err != nil {
		return err
	}
	return bufWriter.Flush()
}
