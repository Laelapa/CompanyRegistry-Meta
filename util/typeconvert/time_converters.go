package typeconvert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TimeToPgtypeTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:             t,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}
