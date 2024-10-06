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
	cursor              int
	selectedNames       map[int]struct{}
	selectedLanguage    string
	selectedSource      string
	names               []string
	languages           []string
	sources             []string
	filePath            textinput.Model
	inputNames          textinput.Model
	inputNamesData      string
	filePathData        string
	currentChoiceLength int
	err                 error
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
	namesTi := textinput.New()
	namesTi.CharLimit = 156
	namesTi.Width = 20
	namesTi.Placeholder = "John, Jane, Alice"
	pathTi := textinput.New()
	pathTi.CharLimit = 156
	pathTi.Width = 20
	pathTi.Placeholder = "names.txt"

	return Model{
		names:               nameChoices[RANDOM],
		selectedNames:       make(map[int]struct{}),
		cursor:              0,
		languages:           lang.GetAvailableLanguages(),
		sources:             []string{RANDOM, FILE, MANUAL},
		inputNames:          namesTi,
		filePath:            pathTi,
		err:                 nil,
		currentChoiceLength: len(lang.GetAvailableLanguages()),
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

		//actual key pressed
		switch msg.String() {

		// The r key randomizes the order of the choices in the list
		case "r":
			m.names, m.selectedNames = utils.RandomizeNonSelectedSlice(m.names, m.selectedNames)

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < m.currentChoiceLength-1 {
				m.cursor++
			}

		case "enter":

			// We handle the enter key differently depending on the current state of the model

			if m.inputNames.Focused() {
				m.inputNames.Blur()
				m.inputNamesData = m.inputNames.Value()
				m.names = utils.SplitString(m.inputNamesData, ",")
				m.currentChoiceLength = len(m.names)
				return m, cmd
			}

			if m.filePath.Focused() {
				m.filePath.Blur()
				m.filePathData = m.filePath.Value()
				names, err := utils.ReadNamesFromFile(m.filePathData)

				if err != nil {
					m.err = err
					return m, nil
				}

				m.names = names
				m.currentChoiceLength = len(m.names)
				return m, cmd
			}

			// If the user has not selected a language yet
			if m.selectedLanguage == "" {
				m.selectedLanguage = m.languages[m.cursor]
				m.currentChoiceLength = len(m.sources)
				return m, nil
			}

			// If the user has not selected a source yet
			if m.selectedSource == "" {
				m.selectedSource = m.sources[m.cursor]

				// If the source is random, we need to randomize the names
				if m.selectedSource == RANDOM {
					m.currentChoiceLength = len(m.names)
					m.names = utils.RandomizeSlice(m.names)
					return m, nil
				}

				// If the source is a file, we need to ask the user for the file path
				if m.selectedSource == FILE {
					m.filePath.Focus()
					m.filePath.Cursor.Blink = true
					return m, nil
				}

				// If the source is manual, we need to ask the user for the names
				if m.selectedSource == MANUAL {
					m.inputNames.Focus()
					m.inputNames.Cursor.Blink = true
					return m, nil
				}

				m.names = inicializeNameChoices()[m.selectedSource]
				m.currentChoiceLength = len(m.names)
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

	// //Are we reading inputs?
	// If we are reading input, we need to handle the input
	// and update the model accordingly

	switch {
	case m.inputNames.Focused():

		m.inputNames, cmd = m.inputNames.Update(msg)
		return m, cmd

	case m.filePath.Focused():

		m.filePath, cmd = m.filePath.Update(msg)
		return m, cmd
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
	}

	// The footer
	footer := Footer{
		Text: map[string]string{
			lang.English: "Press Esc to quit.\n",
			lang.Spanish: "Presiona Esc para salir.\n",
		},
	}

	if m.selectedLanguage == "" {

		view += header.Text[lang.English]

		// Label to select the language
		languageLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Select the language:\n",
				lang.Spanish: "Selecciona el idioma:\n",
			},
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

		return Render(view)

	}

	// Check if an error occurred
	if m.err != nil {
		view += header.Text[m.selectedLanguage]

		view += fmt.Sprintf("Error: %s\n\n", m.err.Error())

		view += footer.Text[m.selectedLanguage]

		return Render(view)
	}

	if m.selectedSource == "" {
		view += header.Text[m.selectedLanguage]

		// Label to select the source
		sourceLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Select the source:\n",
				lang.Spanish: "Selecciona la fuente:\n",
			},
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

		return Render(view)
	}

	if m.filePath.Focused() {

		view += header.Text[m.selectedLanguage]

		// Label to select the source

		sourceLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Enter the file path:\n",
				lang.Spanish: "Ingresa la ruta del archivo:\n",
			},
		}

		view += sourceLabel.Text[m.selectedLanguage]

		view += "\n"

		view += fmt.Sprintf("%s\n", m.filePath.View())

		view += footer.Text[m.selectedLanguage]

		return Render(view)
	}

	if m.inputNames.Focused() {

		view += header.Text[m.selectedLanguage]

		// Label to select the source
		sourceLabel := lang.Label{
			Text: map[string]string{
				lang.English: "Enter the names separated by commas:\n",
				lang.Spanish: "Ingresa los nombres separados por comas:\n",
			},
		}

		view += sourceLabel.Text[m.selectedLanguage]

		view += "\n"

		view += fmt.Sprintf("%s\n", m.inputNames.View())

		view += footer.Text[m.selectedLanguage]

		return Render(view)
	}

	view += header.Text[m.selectedLanguage]

	// Check if names are empty
	if m.names == nil || len(m.names) == 0 {

		emptyLabel := lang.Label{
			Text: map[string]string{
				lang.English: "No names available\n",
				lang.Spanish: "No hay nombres disponibles\n",
			},
		}

		view += emptyLabel.Text[m.selectedLanguage]

		view += footer.Text[m.selectedLanguage]

		return Render(view)
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

	// Add label for "r" key
	randomizeLabel := lang.Label{
		Text: map[string]string{
			lang.English: "\nPress 'r' to randomize the names\n",
			lang.Spanish: "\nPresiona 'r' para mezclar los nombres\n",
		},
	}

	view += randomizeLabel.Text[m.selectedLanguage]

	view += footer.Text[m.selectedLanguage]

	// Send the UI for rendering
	return Render(view)
}
