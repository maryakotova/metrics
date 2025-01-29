package filetransfer

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/maryakotova/metrics/internal/models"
)

type FileReader struct {
	file   *os.File
	reader *bufio.Reader
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (fr *FileReader) Close() error {
	return fr.file.Close()
}

func (fr *FileReader) ReadMetrics() (metrics []*models.Metrics, err error) {

	for {
		data, err := fr.reader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		metric := models.Metrics{}
		err = json.Unmarshal(data, &metric)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

func (fw *FileWriter) WriteMetrics(metrics *[]models.Metrics) error {

	if len(*metrics) == 0 {
		return nil
	}

	for _, value := range *metrics {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		if _, err := fw.writer.Write(data); err != nil {
			return err
		}
		if err := fw.writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	return fw.writer.Flush()
}

func (fw *FileWriter) WriteMetric(metric *models.Metrics) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	if _, err := fw.writer.Write(data); err != nil {
		return err
	}
	if err := fw.writer.WriteByte('\n'); err != nil {
		return err
	}

	return fw.writer.Flush()
}
