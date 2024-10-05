package utils

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

func RandomizeSlice(slice []string) []string {

	if len(slice) <= 1 {
		return slice
	}

	// Fisher-Yates shuffle
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}

func SplitString(s string, sep string) []string {
	//Trim the string
	s = strings.TrimSpace(s)

	slice := strings.Split(s, sep)

	//Trim all the elements
	slice = TrimAll(slice)

	return slice
}

func ReadNamesFromFile(filePath string) ([]string, error) {
	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file line by line
	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	names = TrimAll(names)

	return names, nil
}

func TrimAll(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}
