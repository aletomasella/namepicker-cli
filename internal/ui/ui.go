package ui

import (
	"fmt"

	"github.com/aletomasella/namepicker-cli/internal/lang"
	"github.com/aletomasella/namepicker-cli/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	FILE   = "file"
	RANDOM = "random"
	MANUAL = "manual"
)

type Model struct {
	names            []string
	cursor           int
	selectedNames    map[int]struct{}
	selectedLanguage string
	selectedSource   string
	languages        []string
	sources          []string
	filePath         string
	inputNames       string
}

type Header lang.Label

type Footer lang.Label

func inicializeNameChoices() map[string][]string {
	nameChoices := make(map[string][]string)

	nameSeed := []string{"John", "Jane", "Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy", "Kevin", "Laura", "Michael", "Nancy", "Oliver", "Peggy", "Quincy", "Rita", "Steve", "Tina", "Ursula", "Victor", "Wendy", "Xavier", "Yvonne", "Zack"}

	// Randomize the names
	randomSeed := utils.RandomizeSlice(nameSeed)

	nameChoices[RANDOM] = randomSeed
	nameChoices[FILE] = make([]string, 0)
	nameChoices[MANUAL] = make([]string, 0)

	return nameChoices
}

func InitialModel() Model {

	nameChoices := inicializeNameChoices()

	// Default name choices are the random names

	return Model{
		names:         nameChoices[RANDOM],
		selectedNames: make(map[int]struct{}),
		cursor:        0,
		languages:     lang.GetAvailableLanguages(),
		sources:       []string{RANDOM, FILE, MANUAL},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// The r key randomizes the order of the choices in the list
		case "r":
			m.names = utils.RandomizeSlice(m.names)

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.names)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			// If the user has not selected a language yet
			if m.selectedLanguage == "" {
				m.selectedLanguage = m.languages[m.cursor]
				return m, nil
			}

			// If the user has not selected a source yet
			if m.selectedSource == "" {
				m.selectedSource = m.sources[m.cursor]
				m.names = inicializeNameChoices()[m.selectedSource]
				return m, nil
			}

			_, ok := m.selectedNames[m.cursor]
			if ok {
				delete(m.selectedNames, m.cursor)
			} else {
				m.selectedNames[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {

	// Default language is English
	// First the user must select the language from the available languages
	// Then the user must select from where do they want to get the names from:
	// - Random names
	// - Names from a file
	// - Write the names manually separated by commas

	// view
	view := ""

	// The header
	header := Header{
		Text: map[string]string{
			lang.English: "Random Name Picker\n\n",
			lang.Spanish: "Selector de Nombres\n\n",
		},
		DefineLanguages: lang.GetAvailableLanguages(),
	}

	// The footer
	footer := Footer{
		Text: map[string]string{
			lang.English: "Press q to quit.\n",
			lang.Spanish: "Presiona q para salir.\n",
		},
		DefineLanguages: lang.GetAvailableLanguages(),
	}

	if m.selectedLanguage == "" {

		view += header.Text[lang.English]

		// Label to select the language
		languageLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Select the language:\n",
				lang.Spanish: "Selecciona el idioma:\n",
			},
			DefineLanguages: lang.GetAvailableLanguages(),
		}

		view += languageLabel.Text[lang.English]

		// Iterate over the available languages
		for i, l := range m.languages {
			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected

			view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, l)

		}
		view += footer.Text[lang.English]

		return view

	}

	if m.selectedSource == "" {
		view += header.Text[m.selectedLanguage]

		// Label to select the source
		sourceLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Select the source:\n",
				lang.Spanish: "Selecciona la fuente:\n",
			},
			DefineLanguages: lang.GetAvailableLanguages(),
		}

		view += sourceLabel.Text[m.selectedLanguage]

		// Iterate over the available sources

		for i, s := range m.sources {
			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected

			view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, s)
		}

		view += footer.Text[m.selectedLanguage]

		return view
	}

	view += header.Text[m.selectedLanguage]

	// Iterate over our choices
	for i, choice := range m.names {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selectedNames[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	view += footer.Text[m.selectedLanguage]

	// Send the UI for rendering
	return view
}
