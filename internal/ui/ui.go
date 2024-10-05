package ui

import (
	"fmt"

	"github.com/aletomasella/namepicker-cli/internal/lang"
	"github.com/aletomasella/namepicker-cli/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
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
	filePath         textinput.Model
	inputNames       textinput.Model
	readingPath      bool
	readingInput     bool
	err              error
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
		readingInput:  false,
		readingPath:   false,
		inputNames:    textinput.New(),
		filePath:      textinput.New(),
		err:           nil,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		//Are we reading input?
		if m.inputNames.Focused() {

			// If we are reading input, we need to handle the input
			// and update the model accordingly
			if msg.Type == tea.KeyEnter {
				m.readingInput = false
				m.names = utils.SplitString(m.inputNames.Value(), ",")
				return m, nil
			}

			m.inputNames, cmd = m.inputNames.Update(msg)
			return m, cmd

		}

		// //Are we reading the file path?
		if m.filePath.Focused() {
			// If we are reading input, we need to handle the input
			// and update the model accordingly

			if msg.Type == tea.KeyEnter {
				m.readingPath = false
				names, err := utils.ReadNamesFromFile(m.filePath.Value())

				if err != nil {
					m.err = err
					return m, nil
				}

				m.names = names
				m.readingPath = false

				return m, nil
			}

			m.filePath, cmd = m.filePath.Update(msg)
			return m, cmd
		}

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// The r key randomizes the order of the choices in the list
		case "r":
			m.names = utils.RandomizeSlice(m.names)

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

	case error:
		m.err = msg
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

	// Check if an error occurred
	if m.err != nil {
		view += header.Text[m.selectedLanguage]

		view += fmt.Sprintf("Error: %s\n", m.err.Error())

		view += footer.Text[m.selectedLanguage]
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

	// If selected source is FILE  we need to show the input field

	if m.selectedSource == FILE && m.filePath.Value() == "" {
		view += header.Text[m.selectedLanguage]

		m.readingPath = true
		// Label to select the source

		sourceLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Enter the file path:\n",
				lang.Spanish: "Ingresa la ruta del archivo:\n",
			},
			DefineLanguages: lang.GetAvailableLanguages(),
		}

		view += sourceLabel.Text[m.selectedLanguage]

		m.filePath.Placeholder = "names.txt"
		m.filePath.Focus()

		view += m.filePath.View()

		view += "\n"

		view += footer.Text[m.selectedLanguage]

		return view
	}

	// If selected source is MANUAL we need to show the input field

	if m.selectedSource == MANUAL && m.inputNames.Value() == "" {
		view += header.Text[m.selectedLanguage]

		m.readingInput = true

		// Label to select the source
		sourceLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Enter the names separated by commas:\n",
				lang.Spanish: "Ingresa los nombres separados por comas:\n",
			},
			DefineLanguages: lang.GetAvailableLanguages(),
		}

		view += sourceLabel.Text[m.selectedLanguage]

		m.inputNames.Placeholder = "John, Jane, Alice"
		m.inputNames.Focus()

		view += m.inputNames.View()

		view += "\n"

		view += footer.Text[m.selectedLanguage]

		return view
	}

	view += header.Text[m.selectedLanguage]

	// Check if names are empty
	if len(m.names) == 0 {

		emptyLabel := lang.Label{
			Text: map[string]string{
				lang.English: "No names available\n",
				lang.Spanish: "No hay nombres disponibles\n",
			},
			DefineLanguages: lang.GetAvailableLanguages(),
		}

		view += emptyLabel.Text[m.selectedLanguage]

		view += footer.Text[m.selectedLanguage]

		return view
	}

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
