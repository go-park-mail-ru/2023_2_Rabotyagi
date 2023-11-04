package utils

import (
	"context"
	"fmt"
	"log"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/jackc/pgx/v5"
)

// GetLastValSeq returns id of last record not at all, because sequence auto increment even if unsuccessful insert in table
func GetLastValSeq(ctx context.Context, tx pgx.Tx, nameTable pgx.Identifier) (uint64, error) {
	sanitizedNameTable := nameTable.Sanitize()
	SQLGetLastValSeq := fmt.Sprintf(`SELECT last_value FROM %s;`, sanitizedNameTable)
	seqRow := tx.QueryRow(ctx, SQLGetLastValSeq)

	var count uint64

	if err := seqRow.Scan(&count); err != nil {
		log.Printf("error in GetLastValSeq: %+v", err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return count, nil
}
