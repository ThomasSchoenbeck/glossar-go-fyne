package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var tagSearchLabel *widget.Label
var tagSearchEntry *widget.Entry
var tagSearchEntryClearButton *widget.Button
var tagDeleteButton *widget.Button
var tagList *widget.List
var selectedTag string
var filteredTagList []string

func SetupTags(w fyne.Window) *container.TabItem {
	filteredTagList = tags
	tagList = widget.NewList(func() int {
		return len(filteredTagList)
	},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(filteredTagList[i])
		},
	)

	tagSearchEntry = widget.NewEntry()
	tagSearchEntryClearButton = widget.NewButton("", func() {
		tagSearchEntry.SetText("")
		tagSearchEntryClearButton.Hide()
	})
	tagSearchEntryClearButton.Icon = myTheme.Icon(myTheme{}, theme.IconNameDelete)
	tagSearchEntryClearButton.Importance = widget.LowImportance
	tagSearchEntryClearButton.Resize(fyne.NewSize(30, 30))
	tagSearchEntryClearButton.Move(fyne.NewPos(560, 3))
	tagSearchEntryClearButton.Hide()
	tagSearchEntry.OnChanged = func(s string) {
		if s == "" {
			tagSearchEntryClearButton.Hide()
		} else {
			tagSearchEntryClearButton.Show()
		}
		updateTagList(s)
	}

	tagListScroll := container.NewVScroll(tagList)
	tagListScroll.SetMinSize(fyne.NewSize(400, 200)) // Set minimum size to increase height

	tagSearchLabel = widget.NewLabel(fmt.Sprintf("Search Tags: (%d of %d)", len(filteredTagList), len(tags)))

	tagDeleteButton = widget.NewButton("Delete", func() {
		dialog.ShowConfirm("Delete", fmt.Sprintf("do you realy want to delete this tag: %s? It will also be deleted in every glossar Entry!", selectedTag), func(b bool) {
			if b {
				for i, tag := range tags {
					if tag == selectedTag {
						tagToRemove := selectedTag
						counter := deleteTag(i, tagToRemove)
						updateTagList(searchEntry.Text)
						tagList.Refresh()
						selectedTag = ""
						searchEntry.SetText("")
						tagDeleteButton.Disable()
						dialog.ShowInformation("Delete finished", fmt.Sprintf("Deletion has removed the tag %s in %d glossar entries", tagToRemove, counter), w)
						break
					}
				}
			}
		}, w)
	})
	tagDeleteButton.Importance = widget.DangerImportance
	tagDeleteButton.Disable()

	tagList.OnSelected = func(id widget.ListItemID) {

		newSelectedTag := filteredTagList[id]

		if newSelectedTag == selectedTag {
			tagList.UnselectAll()
			selectedTag = ""
			tagDeleteButton.Disable()
		} else {

			selectedTag = filteredTagList[id]
			tagDeleteButton.Enable()
		}
	}

	return container.NewTabItem("Tags", container.NewVBox(
		container.NewHBox(tagSearchLabel, tagDeleteButton),
		container.NewStack(
			tagSearchEntry,
			container.NewWithoutLayout(layout.NewSpacer(), tagSearchEntryClearButton),
		),
		tagListScroll,
	))
}

func updateTagList(searchString string) {
	// container.Objects = nil
	filteredTagList = []string{}
	lowerSearchTerm := strings.ToLower(searchString)
	if len(lowerSearchTerm) > 0 {

		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), lowerSearchTerm) {
				filteredTagList = append(filteredTagList, tag)
			}
		}
	} else {

		if len(filteredTagList) == 0 {
			filteredTagList = tags
		}
	}

	tagSearchLabel.SetText(fmt.Sprintf("Search Tags: (%d of %d)", len(filteredTagList), len(tags)))
	tagList.Refresh()
}

func RefreshTagList() {
	updateTagList(tagSearchEntry.Text)
}

func deleteTag(deleteIndex int, tagToDelete string) int {
	tags = append(tags[:deleteIndex], tags[deleteIndex+1:]...)

	counter := 0
	// delete from the entire glossary
	for i, entry := range glossary {
		for j, tag := range entry.Tags {
			if tag == tagToDelete {
				glossary[i].Tags = append(glossary[i].Tags[:j], glossary[i].Tags[j+1:]...)
				counter++
			}
		}
	}

	saveGlossary()

	return counter
}
