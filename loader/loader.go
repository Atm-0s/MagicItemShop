package loader

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func ImportCSV(path string) ([]MagicItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	index := map[string]int{}
	for i, h := range header {
		index[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"name", "rarity", "page_number", "table"}
	for _, col := range required {
		if _, ok := index[col]; !ok {
			return nil, fmt.Errorf("missing required column: %v", col)
		}
	}
	_, hasTableNum := index["table_number"]

	var items []MagicItem
	rowNum := 1
	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		rowNum++
		if err != nil {
			return nil, fmt.Errorf("error at row %d: %w", rowNum, err)
		}

		name := strings.TrimSpace(row[index["name"]])
		if name == "" {
			return nil, fmt.Errorf("error at row %d: empty name", rowNum)
		}

		rar, err := ParseRarity(row[index["rarity"]])
		if err != nil {
			return nil, fmt.Errorf("error at row %d: %w", rowNum, err)
		}

		costStr := row[index["cost_gp"]]
		var cost int
		if strings.TrimSpace(costStr) == "" {
			cost = DefaultCostForRarity(rar)
		} else {
			cost, err = atoiField(costStr)
			if err != nil {
				return nil, fmt.Errorf("error at row %d: invalid cost_gp: %w", rowNum, err)
			}
		}

		page, err := atoiField(row[index["page_number"]])
		if err != nil {
			return nil, fmt.Errorf("error at row %d: invalid page_number: %w", rowNum, err)
		}

		tab, err := ParseTable(row[index["table"]])
		if err != nil {
			return nil, fmt.Errorf("error parsing table at row %d: %w", rowNum, err)
		}

		var tabN []int
		if hasTableNum {
			raw := strings.TrimSpace(row[index["table_number"]])
			if raw != "" {
				tabN, err = parseIntList(raw)
				if err != nil {
					return nil, fmt.Errorf("error at row %d (table_number): %w", rowNum, err)
				}
			}
		}

		items = append(items, MagicItem{
			Name:        name,
			Rarity:      rar,
			CostGP:      cost,
			PageNumber:  page,
			Table:       tab,
			TableNumber: tabN,
		})
	}
	return items, nil
}

func atoiField(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// Supports "12", "3-5", "1,4,7-9", "2; 8 ; 10-12"
func parseIntList(s string) ([]int, error) {
	s = strings.ReplaceAll(s, ";", ",")
	parts := strings.Split(s, ",")
	var out []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, "-") {
			b := strings.SplitN(p, "-", 2)
			if len(b) != 2 {
				return nil, fmt.Errorf("bad range %q", p)
			}
			start, err1 := atoiField(b[0])
			end, err2 := atoiField(b[1])
			if err1 != nil || err2 != nil || start > end {
				return nil, fmt.Errorf("bad range %q", p)
			}
			for i := start; i <= end; i++ {
				out = append(out, i)
			}
		} else {
			n, err := atoiField(p)
			if err != nil {
				return nil, fmt.Errorf("bad number %q", p)
			}
			out = append(out, n)
		}
	}
	return out, nil
}

func LoadJSON(path string) ([]MagicItem, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var items []MagicItem
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func SaveJSON(path string, items []MagicItem) error {
	b, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func WriteItemsCSV(path string) error {
	header := []string{"Name", "Rarity", "Cost_GP", "Table", "Table_Number", "Page_Number"}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

// Functions below are to populate a shop with randomly chosen items.

func ShopPopulator(quantity int, rarity Rarity) ([]MagicItem, error) {
	allItems, err := LoadJSON("./content/Items.json")
	var selectedItems []MagicItem
	if err != nil {
		return nil, err
	}
	for i := 0; i < quantity; i++ {
		n := rand.Intn(100) + 1
		for _, item := range allItems {
			if item.Rarity == rarity {
				tableNumMap := make(map[int]struct{}, len(item.TableNumber))
				for _, num := range item.TableNumber {
					tableNumMap[num] = struct{}{}
				}
				if _, ok := tableNumMap[n]; ok {
					selectedItems = append(selectedItems, item)
					break
				}
			}
		}
	}
	return selectedItems, nil
}

func Capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
