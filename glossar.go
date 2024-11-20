package main

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var selectedEntry *GlossaryEntry
var selectedEntryTags []string
var filteredGlossarList []GlossaryEntry
var glossarListSelectedIndex int = -1
var glossarList *widget.List
var glossarSearchEntry *widget.Entry

func SetupGlossar(w fyne.Window) *container.TabItem {

	filteredGlossarList = glossary

	glossarTermEntry := widget.NewEntry()
	glossarDefinitionEntry := widget.NewMultiLineEntry()
	glossarDefinitionEntry.Wrapping = fyne.TextWrapWord

	glossarEntryTagList := widget.NewList(
		func() int {
			return len(selectedEntryTags)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(selectedEntryTags[i])
		},
	)

	glossarList = widget.NewList(
		func() int {
			return len(filteredGlossarList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(filteredGlossarList[i].Term)
		},
	)

	glossarList.OnSelected = func(id widget.ListItemID) {
		selectedEntry = &filteredGlossarList[id]
		// tag list
		selectedEntryTags = selectedEntry.Tags
		glossarEntryTagList.Refresh()
		// input fields
		glossarTermEntry.SetText(selectedEntry.Term)
		glossarDefinitionEntry.SetText(selectedEntry.Definition)
	}

	glossarSearchEntry = widget.NewEntry()
	glossarSearchEntry.OnChanged = func(s string) {
		updateGlossarListSelection(glossary, s)
	}

	glossarSaveButton := widget.NewButton("Save", func() {
		term := strings.TrimSpace(glossarTermEntry.Text)
		definition := glossarDefinitionEntry.Text

		if term == "" {
			dialog.ShowInformation("Error", "Term cannot be empty or just spaces", w)
			return
		}

		if selectedEntry != nil {
			selectedEntry.Definition = definition
		} else {
			glossary = append(glossary, GlossaryEntry{Term: term, Definition: definition})
		}

		glossarTermEntry.SetText("")
		glossarDefinitionEntry.SetText("")
		selectedEntry = nil
		glossarList.Refresh()
		saveGlossary()
	})

	renameButton := widget.NewButton("Rename", func() {
		if selectedEntry == nil {
			dialog.ShowInformation("Error", "No entry selected to rename", w)
			return
		}

		renameEntry := widget.NewEntry()
		renameEntry.SetText(selectedEntry.Term)

		dialog.ShowCustomConfirm("Rename Term", "Rename", "Cancel", renameEntry, func(b bool) {
			if b {
				newTerm := strings.TrimSpace(renameEntry.Text)
				if newTerm == "" {
					dialog.ShowInformation("Error", "Term cannot be empty or just spaces", w)
					return
				}
				selectedEntry.Term = newTerm
				glossarTermEntry.SetText(newTerm)
				glossarList.Refresh()
				saveGlossary()
			}
		}, w)
	})

	clearButton := widget.NewButton("Clear", func() {
		glossarTermEntry.SetText("")
		glossarDefinitionEntry.SetText("")
		selectedEntry = nil
	})

	glossarListContainer := container.NewVScroll(glossarList)
	glossarListContainer.SetMinSize(fyne.NewSize(200, 480)) // Set minimum size to increase height

	glossarListBox := container.NewVBox(
		widget.NewLabel("Search Glossar:"),
		glossarSearchEntry,
		glossarListContainer,
	)

	glossarEntryTagListContainer := container.NewVScroll(glossarEntryTagList)
	glossarEntryTagListContainer.SetMinSize(fyne.NewSize(300, 200))

	background := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	background.SetMinSize(fyne.NewSize(300, 1)) // without this line, the sizing of the right hand side (form fields) will not work

	vbox := container.New(
		layout.NewVBoxLayout(),
		background,
		container.NewHBox(widget.NewLabel("Term:"), renameButton, clearButton),
		glossarTermEntry,
		widget.NewLabel("Definition:"),
		glossarDefinitionEntry,
		glossarSaveButton,
		widget.NewLabel("Tags:"),
		glossarEntryTagListContainer,
	)

	return container.NewTabItem("Glossar",
		container.NewHBox(glossarListBox,
			vbox,
		),
	)

}

func updateGlossarListSelection(glossaryEntries []GlossaryEntry, searchString string) {
	var filteredEntries []GlossaryEntry
	lowerSearchTerm := strings.ToLower(searchString)
	if strings.Contains(lowerSearchTerm, "tag:") {
		tagSearchString := lowerSearchTerm
		tagSearchString = strings.Replace(tagSearchString, "tag:", "", 1)

		for _, entry := range glossaryEntries {
			for _, tag := range entry.Tags {
				if strings.Contains(strings.ToLower(tag), strings.ToLower(tagSearchString)) {
					filteredEntries = append(filteredEntries, entry)
				}
			}
		}

	} else {
		for _, entry := range glossaryEntries {
			if strings.Contains(strings.ToLower(entry.Term), lowerSearchTerm) {
				filteredEntries = append(filteredEntries, entry)
			}
		}
	}

	filteredGlossarList = filteredEntries
	glossarList.Refresh()
}

func FocusGlossarSearch(w fyne.Window) {
	w.Canvas().Focus(glossarSearchEntry)
}
