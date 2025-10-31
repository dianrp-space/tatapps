package notification

import (
	"fmt"
	"sort"
	"strings"

	"tatapps/internal/models"
)

// LowStockEntry describes an inventory warning row used in notifications.
type LowStockEntry struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	WarehouseID   *uint   `json:"warehouse_id,omitempty"`
	WarehouseName string  `json:"warehouse_name,omitempty"`
	Quantity      float64 `json:"quantity"`
	MinStock      float64 `json:"min_stock"`
	Unit          string  `json:"unit"`
	Aggregated    bool    `json:"aggregated"`
	ItemIDs       []uint  `json:"item_ids,omitempty"`
}

// ComputeLowStockEntries rebuilds the dashboard low stock aggregation so other
// components (manual checks, scheduler) stay consistent.
func ComputeLowStockEntries(items []models.InventoryItem) []LowStockEntry {
	aggregatedByCategory := make(map[string]*LowStockEntry)
	var lowStock []LowStockEntry

	for _, item := range items {
		quantity := item.Quantity
		minStock := item.MinStock
		hasSN := strings.TrimSpace(item.SN) != ""
		unit := strings.TrimSpace(item.Unit)
		warehouseName := strings.TrimSpace(item.Warehouse.Name)

		if hasSN {
			category := strings.TrimSpace(item.Category)
			if category == "" {
				category = "Uncategorized"
			}

			key := fmt.Sprintf("%d|%s", item.WarehouseID, strings.ToLower(category))
			entry, exists := aggregatedByCategory[key]
			if !exists {
				entry = &LowStockEntry{
					Category:   category,
					Unit:       unit,
					Aggregated: true,
					ItemIDs:    make([]uint, 0, 1),
				}
				if item.WarehouseID != 0 {
					wid := item.WarehouseID
					entry.WarehouseID = &wid
					entry.WarehouseName = warehouseName
				}
				aggregatedByCategory[key] = entry
			}

			entry.Quantity += quantity
			if minStock > entry.MinStock {
				entry.MinStock = minStock
			}
			if entry.Unit == "" {
				entry.Unit = unit
			}
			if entry.WarehouseName == "" {
				entry.WarehouseName = warehouseName
			}
			entry.ItemIDs = append(entry.ItemIDs, item.ID)
		} else {
			if minStock > 0 && quantity <= minStock {
				entry := LowStockEntry{
					ID:         item.ID,
					Name:       item.Name,
					Category:   item.Category,
					Quantity:   quantity,
					MinStock:   minStock,
					Unit:       unit,
					Aggregated: false,
					ItemIDs:    []uint{item.ID},
				}
				if item.WarehouseID != 0 {
					wid := item.WarehouseID
					entry.WarehouseID = &wid
					entry.WarehouseName = warehouseName
				}
				lowStock = append(lowStock, entry)
			}
		}
	}

	for _, entry := range aggregatedByCategory {
		if entry.MinStock <= 0 {
			continue
		}
		if entry.Quantity <= entry.MinStock {
			displayWarehouse := entry.WarehouseName
			if displayWarehouse == "" {
				displayWarehouse = "Semua Gudang"
			}
			if entry.Name == "" {
				entry.Name = entry.Category
			}
			entry.Name = fmt.Sprintf("%s (%s)", entry.Category, displayWarehouse)
			if entry.Unit == "" {
				entry.Unit = "unit"
			}
			lowStock = append(lowStock, *entry)
		}
	}

	sort.Slice(lowStock, func(i, j int) bool {
		return lowStock[i].Quantity < lowStock[j].Quantity
	})

	return lowStock
}

// BuildLowStockMessage renders a WhatsApp/email-friendly digest matching the
// manual low stock check formatting.
func BuildLowStockMessage(entries []LowStockEntry) string {
	if len(entries) == 0 {
		return ""
	}

	formatQuantity := func(value float64) string {
		if value == 0 {
			return "0"
		}
		if strings.HasSuffix(fmt.Sprintf("%.2f", value), ".00") {
			return fmt.Sprintf("%.0f", value)
		}
		return fmt.Sprintf("%.2f", value)
	}

	message := "🔔 *LOW STOCK ALERT* 🔔\n\n"
	for idx, entry := range entries {
		unit := entry.Unit
		if unit == "" {
			unit = "unit"
		}
		warehouse := entry.WarehouseName
		if warehouse == "" {
			warehouse = "Semua Gudang"
		}
		title := entry.Name
		if title == "" {
			title = entry.Category
		}
		if title == "" {
			title = "Item"
		}

		message += fmt.Sprintf("%d. *%s*\n", idx+1, title)
		message += fmt.Sprintf("   • Gudang: %s\n", warehouse)
		message += fmt.Sprintf("   • Stok tersisa: %s %s\n", formatQuantity(entry.Quantity), unit)
		message += fmt.Sprintf("   • Minimum: %s %s\n", formatQuantity(entry.MinStock), unit)
		if entry.Category != "" && !entry.Aggregated {
			message += fmt.Sprintf("   • Kategori: %s\n", entry.Category)
		}
		if entry.Aggregated {
			message += fmt.Sprintf("   • Akumulasi %d item SN unik\n", len(entry.ItemIDs))
		}
		message += "\n"
	}

	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	message += "⚠️ Segera lakukan restocking untuk item di atas!"

	return message
}
