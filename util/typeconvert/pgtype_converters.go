package typeconvert

import "github.com/jackc/pgx/v5/pgtype"

func PtrStringToPgtypeText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func PgtypeTextToPtrString(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
