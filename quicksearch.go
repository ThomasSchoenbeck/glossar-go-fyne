package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

var listMaxLength int = 5
var selectedIndex int = -1
var filtered []GlossaryEntry
var quickSearchEntryFocused bool = false
var quicksearchEntryClearButton *widget.Button
var appWindowActive bool = false

func SetupQuicksearch(a fyne.App) fyne.Window {
	quickSearchWindow := a.NewWindow("Quick Search")
	a.Lifecycle().SetOnEnteredForeground(func() {
		appWindowActive = true
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		appWindowActive = false
	})
	resultsContainer := container.NewVBox()
	resultsScroll := container.NewVScroll(resultsContainer)
	resultsScroll.SetMinSize(fyne.NewSize(400, 200)) // Ensure space for at least 5 items

	searchEntry = newCustomEntry(func() {
		quickSearchEntryFocused = true
		filtered = glossary
		updateSelection(resultsContainer, selectedIndex, filtered)
	}, func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyDown {
			if selectedIndex < len(filtered)-1 {
				selectedIndex++
			}
		} else if key.Name == fyne.KeyUp {
			if selectedIndex > 0 {
				selectedIndex--
			}
		}
		updateSelection(resultsContainer, selectedIndex, filtered)
	}, func() {
		quickSearchEntryFocused = false
	},
	)
	suggestionLabel := widget.NewRichTextFromMarkdown("")
	suggestionLabel.Hide()

	quicksearchEntryClearButton = widget.NewButton("", func() {
		searchEntry.SetText("")
		quicksearchEntryClearButton.Hide()
	})
	quicksearchEntryClearButton.Icon = myTheme.Icon(myTheme{}, theme.IconNameDelete)
	quicksearchEntryClearButton.Importance = widget.LowImportance
	quicksearchEntryClearButton.Resize(fyne.NewSize(30, 30))
	quicksearchEntryClearButton.Move(fyne.NewPos(365, 3))
	quicksearchEntryClearButton.Hide()

	searchEntry.OnChanged = func(s string) {
		resultsContainer.Objects = nil
		filtered = []GlossaryEntry{}
		for _, entry := range glossary {
			if s == "" || containsIgnoreCase(entry.Term, s) {
				filtered = append(filtered, entry)
			}
		}

		for i, entry := range filtered {
			if i >= listMaxLength {
				break
			}
			highlightedTerm := highlightMatch(entry.Term, s)
			result := widget.NewRichTextFromMarkdown(highlightedTerm + ": " + entry.Definition)
			resultsContainer.Add(result)
		}

		if s == "" {
			quicksearchEntryClearButton.Hide()
		} else {
			quicksearchEntryClearButton.Show()
		}

		resultsContainer.Refresh()
	}

	quickSearchWindow.SetContent(container.NewVBox(
		container.NewStack(
			searchEntry,
			container.NewWithoutLayout(suggestionLabel),
			container.NewWithoutLayout(quicksearchEntryClearButton),
		),
		resultsScroll,
	))

	// Register global hotkey for ALT + SPACE and ESC to hide/show quick search window
	go func() {
		hook.Register(hook.KeyDown, []string{"ctrl", "alt", "space"}, func(e hook.Event) {
			if quickSearchVisible {
				quickSearchWindow.Hide()
				quickSearchVisible = false
			} else {
				quickSearchWindow.Show()
				quickSearchVisible = true
				FocusSearchEntry(quickSearchWindow)
			}
		})
		hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
			if quickSearchVisible && !quickSearchEntryFocused && appWindowActive {
				quickSearchWindow.Hide()
				quickSearchVisible = false
			}
			if quickSearchEntryFocused && searchEntry.Text == "" {
				quickSearchWindow.Canvas().Focus(nil)
			}
			if quickSearchEntryFocused && searchEntry.Text != "" {
				searchEntry.SetText("")
			}
		})

		s := hook.Start()
		<-hook.Process(s)
	}()

	return quickSearchWindow
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func updateSelection(container *fyne.Container, index int, filtered []GlossaryEntry) {
	container.Objects = nil
	start := 0
	if index >= listMaxLength {
		start = index - 4
	}
	end := start + listMaxLength
	if end > len(filtered) {
		end = len(filtered)
	}
	for i := start; i < end; i++ {
		richText := widget.NewRichTextFromMarkdown(filtered[i].Term + ": " + filtered[i].Definition)
		richText.Wrapping = fyne.TextWrapWord
		if i == index {
			richText.ParseMarkdown("**" + filtered[i].Term + "**: " + filtered[i].Definition)
		}
		container.Add(richText)
	}
	container.Refresh()
}

func highlightMatch(term, query string) string {
	lowerTerm := strings.ToLower(term)
	lowerQuery := strings.ToLower(query)
	start := strings.Index(lowerTerm, lowerQuery)
	if start == -1 {
		return term
	}
	end := start + len(query)
	return term[:start] + "**" + term[start:end] + "**" + term[end:]
}

func FocusSearchEntry(w fyne.Window) {
	w.Canvas().Focus(searchEntry)
}
