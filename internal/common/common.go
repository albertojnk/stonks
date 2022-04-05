package common

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateUUID generate uuid
func GenerateUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

//GetEnv get environment variable passing a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

//ReadCsv function.
func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	reader := csv.NewReader(f)
	reader.Comma = ';'
	lines, err := reader.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

// ParseIfIsNaNOrInf parse to zero if value is NaN or Inf.
func ParseIfIsNaNOrInf(value float64) float64 {
	if math.IsNaN(value) {
		return 0.0
	}

	if math.IsInf(value, 1) {
		return 0.0
	}

	if math.IsInf(value, -1) {
		return 0.0
	}

	if math.IsInf(value, 0) {
		return 0.0
	}
	return value
}

func Unique(slice []interface{}) []interface{} {
	keys := make(map[interface{}]bool)
	list := []interface{}{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetFirstOfMonth() (time.Time, time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	firstNextMonth := firstOfMonth.AddDate(0, 1, 0)

	return firstOfMonth, firstNextMonth
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
