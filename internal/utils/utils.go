package utils

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

func RandomizeNonSelectedSlice(slice []string, selected map[int]struct{}) ([]string, map[int]struct{}) {
	// Create a new slice to store the non-selected elements
	var nonSelected []string
	var selectedSlice []string

	// Iterate over the slice
	for i, s := range slice {
		// If the element is not selected, add it to the new slice
		if _, ok := selected[i]; !ok {
			nonSelected = append(nonSelected, s)
			continue
		}
		selectedSlice = append(selectedSlice, s)
	}

	// We loop from selected & add them to the selected slice
	newMap := make(map[int]struct{})

	for i := 0; i < len(selectedSlice); i++ {
		newMap[i] = struct{}{}
	}

	// Randomize the new slice
	return append(selectedSlice, RandomizeSlice(nonSelected)...), newMap
}

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
