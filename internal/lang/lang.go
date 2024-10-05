package lang

import (
	"errors"
	"slices"
)

type Label struct {
	// The text to display
	Text            map[string]string
	DefineLanguages []string
}

// define all available languages
const (
	English = "en"
	Spanish = "es"
)

func GetAvailableLanguages() []string {
	return []string{English, Spanish}
}

// Set the text for a given language
func (l *Label) SetText(lang string, text string) error {
	// validate the language
	if lang != English && lang != Spanish {
		return errors.New("unsupported language")
	}

	// set the text
	l.Text[lang] = text

	// add the language to the list of defined languages
	if !slices.Contains(l.DefineLanguages, lang) {
		l.DefineLanguages = append(l.DefineLanguages, lang)
	}

	return nil
}
