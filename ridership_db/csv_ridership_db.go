package ridershipDB

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type CsvRidershipDB struct {
	idIdxMap      map[string]int
	csvFile       *os.File
	csvReader     *csv.Reader
	num_intervals int
}

func (c *CsvRidershipDB) Open(filePath string) error {
	c.num_intervals = 9

	// Create a map that maps MBTA's time period ids to indexes in the slice
	c.idIdxMap = make(map[string]int)
	for i := 1; i <= c.num_intervals; i++ {
		timePeriodID := fmt.Sprintf("time_period_%02d", i)
		c.idIdxMap[timePeriodID] = i - 1
	}

	// create csv reader
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	c.csvFile = csvFile
	c.csvReader = csv.NewReader(c.csvFile)

	return nil
}

func (c *CsvRidershipDB) GetRidership(lineId string) ([]int64, error) {
	rows, err := c.csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	values := make([]int64, c.num_intervals)
	for _, row := range rows[1:] {
		timePeriodID := ""
		if row[0] == lineId {
			if row[2] != timePeriodID {
				timePeriodID = row[2]
			}
			val, err := strconv.ParseInt(row[4], 10, 64)
			if err != nil {
				return nil, err
			}
			values[c.idIdxMap[timePeriodID]] += val
		}
	}

	return values[:], nil
}

func (c *CsvRidershipDB) Close() error {
	return c.csvFile.Close()
}
