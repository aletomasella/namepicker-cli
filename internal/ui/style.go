package ui

import "github.com/charmbracelet/lipgloss"

// var style = lipgloss.NewStyle().
// 	Bold(true).
// 	Foreground(lipgloss.Color("#FAFAFA")).
// 	Background(lipgloss.Color("#7D56F4")).
// 	PaddingTop(2).
// 	PaddingLeft(4).
// 	Width(22)

var style = lipgloss.NewStyle().
	Bold(true).
	//	BorderForeground(lipgloss.Color("63")).
	// blue foreground
	Foreground(lipgloss.Color("35")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22).
	BorderStyle(lipgloss.RoundedBorder())

// // Set a rounded, yellow-on-purple border to the top and left
// var style = lipgloss.NewStyle().
// 	BorderStyle(lipgloss.RoundedBorder()).
// 	BorderForeground(lipgloss.Color("228")).
// 	BorderBackground(lipgloss.Color("63")).
// 	BorderTop(true).
// 	BorderLeft(true)

// Make your own border
var myCuteBorder = lipgloss.Border{
	Top:         "._.:*:",
	Bottom:      "._.:*:",
	Left:        "|*",
	Right:       "|*",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

func Style() lipgloss.Style {
	return style
}

func Render(view string) string {
	return Style().Render(view)
}
