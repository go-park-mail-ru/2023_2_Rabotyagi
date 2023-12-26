package repository

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

func OrderByClauseForProductList(premiumCoefficient, nonPremiumCoefficient,
	soldByUserCoefficient, viewsCoefficient uint16) string { //nolint: gofumpt
	return fmt.Sprintf(`CASE
		WHEN premium_status = %d THEN ((%d * views + %d * (SELECT COUNT(*) FROM product p2
		WHERE p2.saler_id = product.saler_id)) * %d) 
		ELSE ((%d * views + %d * (SELECT COUNT(*) FROM product p2
		WHERE p2.saler_id = product.saler_id)) * %d) 
		END DESC`, statuses.IntStatusPremiumSucceeded, premiumCoefficient, soldByUserCoefficient, premiumCoefficient,
		viewsCoefficient, soldByUserCoefficient, nonPremiumCoefficient)
}
