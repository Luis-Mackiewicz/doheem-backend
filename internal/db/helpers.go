package db

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDFromString(s string) pgtype.UUID {
	var u pgtype.UUID
	if s != "" {
		u.Scan(s)
	}
	return u
}

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	v, _ := u.Value()
	return v.(string)
}

func NumericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(strconv.FormatFloat(f, 'f', -1, 64)); err != nil {
		n.Valid = false
		return n
	}
	return n
}

func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	v, _ := n.Value()
	s, ok := v.(string)
	if !ok {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func TextFromStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func TextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func UUIDFromStringPtrOrZero(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{}
	}
	return UUIDFromString(*s)
}

func UUIDFromStringPtr(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{}
	}
	return UUIDFromString(*s)
}

func TimeFromTimestamptz(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func TimestamptzFromTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func TimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func DateFromTime(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func DateFromTimePtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func DateToTimePtr(t pgtype.Date) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func Int2FromInt16Ptr(i *int16) pgtype.Int2 {
	if i == nil {
		return pgtype.Int2{}
	}
	return pgtype.Int2{Int16: *i, Valid: true}
}

func Int2ToInt16Ptr(i pgtype.Int2) *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

func JSONToMap(b []byte) map[string]interface{} {
	if b == nil {
		return nil
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m
}

func MapToJSON(m map[string]interface{}) []byte {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	return b
}

func UUIDToStringPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	s := UUIDToString(u)
	return &s
}
