package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var selectedEntry *GlossaryEntry
var selectedEntryTags []string
var filteredGlossarList []GlossaryEntry
var glossarListSelectedIndex int = -1
var glossarList *widget.List
var glossarSearchEntry *widget.Entry
var glossarSearchLabel *widget.Label
var glossarTermEntry *widget.Entry
var glossarDefinitionEntry *widget.Entry
var glossarEntryTagList *widget.List
var renameButton *widget.Button
var deleteButton *widget.Button
var addTagButton *widget.Button
var addTagEntry *widget.Entry

var UpdatedTags []string

var selectedGlossarTag string
var glossarTagDeleteButton *widget.Button
var glossarSearchEntryClearButton *widget.Button

func SetupGlossar(w fyne.Window) *container.TabItem {

	filteredGlossarList = glossary

	glossarTermEntry = widget.NewEntry()
	glossarDefinitionEntry = widget.NewMultiLineEntry()
	glossarDefinitionEntry.Wrapping = fyne.TextWrapWord

	glossarEntryTagList = widget.NewList(
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

	glossarTagDeleteButton = widget.NewButton("Delete", func() {

		for i, tag := range selectedEntryTags {
			if tag == selectedGlossarTag {
				selectedEntryTags = append(selectedEntryTags[:i], selectedEntryTags[i+1:]...)
				glossarEntryTagList.Refresh()
				selectedGlossarTag = ""
				break
			}
		}

	})

	glossarTagDeleteButton.Disable()

	glossarEntryTagList.OnSelected = func(id widget.ListItemID) {
		selectedGlossarTag = selectedEntryTags[id]
		glossarTagDeleteButton.Enable()
	}

	renameButton = widget.NewButton("Rename", func() {
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
	renameButton.Disable()

	deleteButton = widget.NewButton("Delete", func() {

		dialog.ShowConfirm("Delete", fmt.Sprintf("do you realy want to delete the entry: %s?", selectedEntry.Term), func(b bool) {
			if b {
				for i, entry := range glossary {
					if entry.Id == selectedEntry.Id {
						glossary = append(glossary[:i], glossary[i+1:]...)
						break
					}
				}
				updateGlossarListSelection(glossary, glossarSearchEntry.Text)
				saveGlossary()
				clearSelection()
			}
		}, w)

	})
	deleteButton.Importance = widget.DangerImportance
	deleteButton.Disable()

	glossarList.OnSelected = func(id widget.ListItemID) {
		selectedEntry = &filteredGlossarList[id]
		deleteButton.Enable()
		renameButton.Enable()
		// tag list
		selectedEntryTags = selectedEntry.Tags
		glossarEntryTagList.Refresh()
		// input fields
		glossarTermEntry.SetText(selectedEntry.Term)
		glossarDefinitionEntry.SetText(selectedEntry.Definition)
	}

	// Glossar Search <start>
	glossarSearchEntry = widget.NewEntry()
	glossarSearchEntryClearButton = widget.NewButton("", func() {
		glossarSearchEntry.SetText("")
		glossarSearchEntryClearButton.Hide()
	})
	glossarSearchEntryClearButton.Icon = myTheme.Icon(myTheme{}, theme.IconNameDelete)
	glossarSearchEntryClearButton.Importance = widget.LowImportance
	glossarSearchEntryClearButton.Resize(fyne.NewSize(30, 30))
	glossarSearchEntryClearButton.Move(fyne.NewPos(243, 3))
	glossarSearchEntryClearButton.Hide()

	glossarSearchEntry.OnChanged = func(s string) {
		if s == "" {
			glossarSearchEntryClearButton.Hide()
		} else {
			glossarSearchEntryClearButton.Show()
		}
		updateGlossarListSelection(glossary, s)
	}
	// Glossar Search <end>

	glossarSaveButton := widget.NewButton("Save", func() {
		term := strings.TrimSpace(glossarTermEntry.Text)
		definition := glossarDefinitionEntry.Text

		if term == "" {
			dialog.ShowInformation("Error", "Term cannot be empty or just spaces", w)
			return
		}
		if UpdatedTags != nil {
			tags = UpdatedTags
		}

		if selectedEntry != nil {
			selectedEntry.Definition = definition
			selectedEntry.Tags = selectedEntryTags

			if selectedEntry.Id == "" {
				for i, entry := range glossary {
					if entry.Term == selectedEntry.Term {
						glossary[i].Id = newNanoId(w)
						glossary[i].Definition = selectedEntry.Definition
						glossary[i].Tags = selectedEntry.Tags
					}
				}
			} else {

				for i, entry := range glossary {
					if entry.Id == selectedEntry.Id {
						glossary[i].Definition = selectedEntry.Definition
						glossary[i].Tags = selectedEntry.Tags
					}
				}
			}
		} else {
			glossary = append(glossary, GlossaryEntry{Id: newNanoId(w), Term: term, Definition: definition, Tags: selectedEntryTags})
		}

		clearSelection()
		updateGlossarListSelection(glossary, glossarSearchEntry.Text)
		saveGlossary()
	})
	glossarSaveButton.Importance = widget.HighImportance

	clearButton := widget.NewButton("Clear", func() {
		clearSelection()
	})

	addTagEntry = widget.NewEntry()
	addTagButton = widget.NewButton("Add", func() {
		dialog.ShowForm("New Tag", "add", "cancel", []*widget.FormItem{widget.NewFormItem("", addTagEntry)}, func(b bool) {
			if b {
				newTag := addTagEntry.Text
				selectedEntryTags = append(selectedEntryTags, newTag)
				glossarEntryTagList.Refresh()
				tagInList := false
				for _, tag := range tags {
					if tag == newTag {
						tagInList = true
					}
				}
				if !tagInList {
					UpdatedTags = tags
					UpdatedTags = append(UpdatedTags, newTag)
				}
			}
			addTagEntry.SetText("")
		}, w)
	})

	glossarListContainer := container.NewVScroll(glossarList)
	glossarListContainer.SetMinSize(fyne.NewSize(275, 480)) // Set minimum size to increase height

	glossarSearchLabel = widget.NewLabel(fmt.Sprintf("Search Glossar: %d of %d", len(filteredGlossarList), len(glossary)))

	glossarListBox := container.NewVBox(
		glossarSearchLabel,
		container.NewStack(
			glossarSearchEntry,
			// container.NewVBox(layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), glossarSearchEntryClearButton)),
			container.NewWithoutLayout(layout.NewSpacer(), glossarSearchEntryClearButton),
		),
		glossarListContainer,
	)

	glossarEntryTagListContainer := container.NewVScroll(glossarEntryTagList)
	glossarEntryTagListContainer.SetMinSize(fyne.NewSize(300, 200))

	// background := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	background := canvas.NewRectangle(myTheme.Color(myTheme{}, theme.ColorNameBackground, theme.VariantDark))
	background.SetMinSize(fyne.NewSize(300, 1)) // without this line, the sizing of the right hand side (form fields) will not work

	vbox := container.New(
		layout.NewVBoxLayout(),
		background,
		container.NewHBox(widget.NewLabel("Term:"), renameButton, clearButton, deleteButton),
		glossarTermEntry,
		widget.NewLabel("Definition:"),
		glossarDefinitionEntry,
		glossarSaveButton,
		container.NewHBox(widget.NewLabel("Tags:"), addTagButton, glossarTagDeleteButton),
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
	glossarSearchLabel.SetText(fmt.Sprintf("Search Glossar: %d of %d", len(filteredGlossarList), len(glossary)))
	glossarList.Refresh()
}

func FocusGlossarSearch(w fyne.Window) {
	w.Canvas().Focus(glossarSearchEntry)
}

func clearSelection() {
	glossarTermEntry.SetText("")
	glossarDefinitionEntry.SetText("")
	selectedEntry = nil
	selectedEntryTags = []string{}
	glossarEntryTagList.Refresh()
	deleteButton.Disable()
	renameButton.Disable()
	UpdatedTags = nil
	glossarTagDeleteButton.Disable()
	selectedGlossarTag = ""
}

func newNanoId(w fyne.Window) string {
	id, err := gonanoid.New()
	if err != nil {
		dialog.ShowError(err, w)
	}
	return id
}
