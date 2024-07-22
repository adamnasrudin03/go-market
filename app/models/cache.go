package models

import "fmt"

func KeyCacheProductDetail(id uint64) string {
	return fmt.Sprintf("product_detail_%d", id)
}
