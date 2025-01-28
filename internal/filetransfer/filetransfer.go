package filetransfer

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/maryakotova/metrics/internal/models"
)

type FileReader struct {
	file *os.File
	// decoder *json.Decoder
	reader *bufio.Reader
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file: file,
		// decoder: json.NewDecoder(file),
		reader: bufio.NewReader(file),
	}, nil
}

func (fr *FileReader) Close() error {
	return fr.file.Close()
}

func (fr *FileReader) ReadMetrics() (metrics []*models.Metrics, err error) {
	data, err := fr.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

type FileWriter struct {
	file *os.File
	// encoder *json.Encoder
	writer *bufio.Writer
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file: file,
		// encoder: json.NewEncoder(file),
		writer: bufio.NewWriter(file),
	}, nil
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

func (fw *FileWriter) WriteMetrics(metrics *[]models.Metrics) error {
	data, err := json.Marshal(metrics)
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
