package main

import (
	"MagicItemShop/loader"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func GenerateButton(myApp fyne.App, quantity *widget.Entry, rarity loader.Rarity) func() {
	return func() {
		qty, err := strconv.Atoi(quantity.Text)
		if err != nil || qty <= 0 {
			errorWindow := myApp.NewWindow("Error")
			errorWindow.SetContent(widget.NewLabel("Quantity must be a positive integer number"))
			errorWindow.Show()
			return
		}
		items, err := loader.ShopPopulator(qty, rarity)
		if err != nil {
			errorWindow := myApp.NewWindow("Error")
			errorWindow.SetContent(widget.NewLabel(fmt.Sprintf("Error generating items: %v", err)))
			errorWindow.Show()
			return
		}

		var gridItems []fyne.CanvasObject
		headers := []fyne.CanvasObject{
			widget.NewLabel("Name"),
			widget.NewLabel("Rarity"),
			widget.NewLabel("Cost GP"),
			widget.NewLabel("Table"),
			widget.NewLabel("Page Number"),
		}
		gridItems = append(gridItems, headers...)
		for _, items := range items {
			gridItems = append(gridItems,
				widget.NewLabel(items.Name),
				widget.NewLabel(string(items.Rarity)),
				widget.NewLabel(fmt.Sprintf("%d", items.CostGP)),
				widget.NewLabel(string(items.Table)),
				widget.NewLabel(fmt.Sprintf("%d", items.PageNumber)),
			)
		}

		grid := container.NewGridWithColumns(5, gridItems...)
		w := myApp.NewWindow(string(rarity))
		w.SetContent(grid)
		w.Resize(fyne.NewSize(600, 400))
		w.Show()

	}
}

func BuildUI(myApp fyne.App) *fyne.Container {
	quantity := widget.NewEntry()
	quantity.SetPlaceHolder("Enter quantity (Only positive integers are valid)")
	commonButton := widget.NewButton("Generate Common Items", GenerateButton(myApp, quantity, loader.Common))
	uncommonButton := widget.NewButton("Generate Uncommon Items", GenerateButton(myApp, quantity, loader.Uncommon))
	rareButton := widget.NewButton("Generate Rare Items", GenerateButton(myApp, quantity, loader.Rare))
	veryRareButton := widget.NewButton("Generate Very Rare Items", GenerateButton(myApp, quantity, loader.VeryRare))
	legendaryButton := widget.NewButton("Generate Legendary Items", GenerateButton(myApp, quantity, loader.Legendary))
	rarityButtons := container.NewVBox(commonButton, uncommonButton, rareButton, veryRareButton, legendaryButton)
	return container.NewVBox(quantity, rarityButtons)
}
