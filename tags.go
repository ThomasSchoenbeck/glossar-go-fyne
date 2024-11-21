package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var tagSearchLabel *widget.Label
var tagSearchEntry *widget.Entry

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
	tagSearchEntry.OnChanged = func(s string) {
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
						deleteTag(i, selectedTag)
						updateTagList(searchEntry.Text)
						tagList.Refresh()
						selectedTag = ""
						searchEntry.SetText("")
						tagDeleteButton.Disable()
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
		tagSearchEntry,
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

func deleteTag(deleteIndex int, tagToDelete string) {
	tags = append(tags[:deleteIndex], tags[deleteIndex+1:]...)

	// delete from the entire glossary
	for i, entry := range glossary {
		for j, tag := range entry.Tags {
			if tag == tagToDelete {
				glossary[i].Tags = append(glossary[i].Tags[:j], glossary[i].Tags[j+1:]...)
			}
		}
	}

	saveGlossary()

}
