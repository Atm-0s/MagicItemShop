A go program that uses Fyne to display a desired quantity of magic items by rarity.
The aim with this is to make generating the items for a shop for quick lookup.
When built and run it will look for a content/Items.csv from root and create one with the correct header if none are found.
All you need to do then is data entry for :-
. The Name of the item: Name.
. The Rarity of the item: Common, Uncommon, Rare, Very Rare and Legendary.
. Cost in GP (Also calculated based on rarity so you can ,, instead): Cost_GP.
. The table it's from: Table.
. The d100 number you roll to choose the item: Table_Number.
. The page number you'll find the item on: Page_Number.

When entering multiple tables make sure to add 100 to the roll number on the table i.e 0-100 for Arcana, 101-200 for Armaments, 201-300 for Implements and 301-400 for Relics

The CSV gets converted to a JSON which is used to build a slice of MagicItem structs. The UI has you enter a desired quantity
then press the button for the rarity you want, such as Rare, which generates a new window with a list of the items and the info
about them. Multiple windows can be opened and the program closes when the main window is closed.

I tested this a bunch but I don't really think I can share the csv file I've made, that stuff must be sourced ethically.
