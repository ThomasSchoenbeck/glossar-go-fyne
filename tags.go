package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SetupTags(w fyne.Window) *container.TabItem {
	tagList := widget.NewList(
		func() int {
			return len(tags)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(tags[i])
		},
	)

	tagListContainer := container.NewVScroll(tagList)
	tagListContainer.SetMinSize(fyne.NewSize(400, 200)) // Set minimum size to increase height

	return container.NewTabItem("Tags", container.NewVBox(
		tagListContainer,
	))
}
