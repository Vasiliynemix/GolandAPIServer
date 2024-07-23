package logging

import (
	"archive/zip"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"time"
)

type Zipper struct {
	files []*os.File
}

func InitZipper(files []*os.File) *Zipper {
	return &Zipper{files: files}
}

func (z *Zipper) ZipLogFiles(log *zap.Logger, maxSize int64, structDateFormat string) {
	files := z.files
	for _, file := range files {
		fileInfo, err := file.Stat()
		if err != nil {
			log.Error("Error getting file info", zap.Error(err))
			continue
		}

		if fileInfo.Size() >= maxSize {
			log.Info("Clearing log file and creating new zip file...", zap.String("file", file.Name()))

			zipFilePath := z.generateZipFileName(file.Name(), structDateFormat)
			if err = z.createZipFile(file.Name(), zipFilePath); err != nil {
				log.Error("Error creating zip file", zap.Error(err))
			}

			err = file.Truncate(0)
			if err != nil {
				log.Error("Error clearing log file", zap.Error(err))
				continue
			}
		}
	}
}

func (z *Zipper) generateZipFileName(sourceFilePath string, structDateFormat string) string {
	pathSlice := strings.Split(sourceFilePath, "/")
	fileName := pathSlice[len(pathSlice)-1]

	newFileName := fmt.Sprintf("%s_%s", time.Now().Format(structDateFormat), fileName)

	filePath := strings.Replace(sourceFilePath, fileName, newFileName, 1)
	zipFilePath := strings.Replace(filePath, ".log", ".zip", 1)
	return zipFilePath
}

func (z *Zipper) createZipFile(sourceFilePath, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer func(zipFile *os.File) {
		_ = zipFile.Close()
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
