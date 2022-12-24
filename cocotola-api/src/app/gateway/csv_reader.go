package gateway

import (
	"encoding/csv"
	"errors"
	"io"

	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func ReadCSV(fileReader io.Reader, fn func(i int, line []string) error) error {
	csvReader := csv.NewReader(fileReader)
	var i = 1
	for {
		var line []string
		line, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return liberrors.Errorf("csvReader.Read. err: %w", err)
		}

		if err := fn(i, line); err != nil {
			return err
		}
		i++
	}
	return nil
}
