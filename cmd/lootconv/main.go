package main

import (
	"MagicItemShop/loader"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	path := "./content/Items.csv"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir("./content", 0755)
		err = loader.WriteItemsCSV(path)
		if err != nil {
			fmt.Printf("Error writing Items CSV: %v", err)
			os.Exit(1)
		}
	}

	csvData, err := loader.ImportCSV("./content/Items.csv")
	if err != nil {
		fmt.Printf("Error importing CSV: %v", err)
		os.Exit(1)
	}
	err = loader.SaveJSON("./content/Items.json", csvData)
	if err != nil {
		fmt.Printf("Error saving JSON: %v", err)
		os.Exit(1)
	}
	fmt.Println("CSV to JSON conversion successful")

	myApp := app.New()
	masterWindow := myApp.NewWindow("Magic Item Populator")
	masterWindow.SetMaster()
	masterWindow.SetContent(BuildUI(myApp))
	masterWindow.Resize(fyne.NewSize(400, 200))

	masterWindow.Show()

	myApp.Run()
}
