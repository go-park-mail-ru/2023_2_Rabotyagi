package repository

import "fmt"

func OrderByClauseForProductList(premiumCoefficient, nonPremiumCoefficient,
	soldByUserCoefficient, viewsCoefficient uint16) string { //nolint: gofumpt
	return fmt.Sprintf(`CASE
		WHEN premium = true THEN ((%d * views + %d * (SELECT COUNT(*) FROM product p2
		WHERE p2.saler_id = product.saler_id)) * %d) 
		ELSE ((%d * views + %d * (SELECT COUNT(*) FROM product p2
		WHERE p2.saler_id = product.saler_id)) * %d) 
		END DESC`, premiumCoefficient, soldByUserCoefficient, premiumCoefficient,
		viewsCoefficient, soldByUserCoefficient, nonPremiumCoefficient)
}
