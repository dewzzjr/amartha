package files

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

func Reader(path string) (<-chan []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(f)

	c := make(chan []string)
	go func() {
		defer close(c)
		record, err := reader.Read()
		for record != nil {
			c <- record
			record, err = reader.Read()
		}
		if err != nil && err != io.EOF {
			log.Error().Err(err).Str("path", path).Msg("error read csv file")
		}
	}()

	return c, nil
}
