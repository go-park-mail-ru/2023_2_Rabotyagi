package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/myerrors"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
)

// GetLastValSeq returns id of last record not at all, because sequence auto increment even if unsuccessful insert in table
func GetLastValSeq(ctx context.Context, tx pgx.Tx, logger *zap.SugaredLogger, nameTable pgx.Identifier) (uint64, error) {
	sanitizedNameTable := nameTable.Sanitize()
	SQLGetLastValSeq := fmt.Sprintf(`SELECT last_value FROM %s;`, sanitizedNameTable)
	seqRow := tx.QueryRow(ctx, SQLGetLastValSeq)

	var count uint64

	if err := seqRow.Scan(&count); err != nil {
		logger.Errorf("error in GetLastValSeq: %+v", err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return count, nil
}
