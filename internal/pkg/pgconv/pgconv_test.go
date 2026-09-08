package pgconv_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/pkg/pgconv"
)

func TestUUIDConversions(t *testing.T) {
	orig := uuid.New()

	pg := pgconv.UUID(orig)
	if !pg.Valid || pg.Bytes != orig {
		t.Fatalf("expected valid pgtype with same bytes, got %v", pg)
	}

	back := pgconv.ToUUID(pg)
	if back != orig {
		t.Fatalf("expected %s, got %s", orig, back)
	}

	optPg := pgconv.OptUUID(&orig)
	if !optPg.Valid || optPg.Bytes != orig {
		t.Fatalf("expected valid opt UUID")
	}

	optBack := pgconv.ToOptUUID(optPg)
	if optBack == nil || *optBack != orig {
		t.Fatalf("expected %s, got %v", orig, optBack)
	}

	nilPg := pgconv.OptUUID(nil)
	if nilPg.Valid {
		t.Fatalf("expected invalid pgtype for nil UUID")
	}
	if pgconv.ToOptUUID(nilPg) != nil {
		t.Fatalf("expected nil from invalid pgtype")
	}
}

func TestTextConversions(t *testing.T) {
	str := "Hello POS"
	pg := pgconv.Text(str)
	if !pg.Valid || pg.String != str {
		t.Fatalf("expected %s, got %v", str, pg)
	}

	if pgconv.ToText(pg) != str {
		t.Fatalf("expected %s", str)
	}

	optPg := pgconv.OptText(&str)
	if !optPg.Valid || *pgconv.ToOptText(optPg) != str {
		t.Fatalf("expected valid opt text")
	}

	nilPg := pgconv.OptText(nil)
	if nilPg.Valid || pgconv.ToOptText(nilPg) != nil {
		t.Fatalf("expected nil for nil string")
	}
}

func TestTimeConversions(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	pg := pgconv.Timestamptz(now)
	if !pg.Valid || !pg.Time.Equal(now) {
		t.Fatalf("expected %v, got %v", now, pg.Time)
	}

	if !pgconv.ToTime(pg).Equal(now) {
		t.Fatalf("expected equal time")
	}

	optPg := pgconv.OptTimestamptz(&now)
	if !optPg.Valid || !pgconv.ToOptTime(optPg).Equal(now) {
		t.Fatalf("expected valid opt time")
	}

	nilPg := pgconv.OptTimestamptz(nil)
	if nilPg.Valid || pgconv.ToOptTime(nilPg) != nil {
		t.Fatalf("expected nil for nil time")
	}
}

func TestNumericConversions(t *testing.T) {
	testCases := []int64{
		0,
		1,
		-1,
		125000,
		9999999999,
		-50000,
	}

	for _, tc := range testCases {
		num := pgconv.Int64ToNumeric(tc)
		if !num.Valid {
			t.Fatalf("expected valid numeric for %d", tc)
		}

		back, err := pgconv.NumericToInt64(num)
		if err != nil {
			t.Fatalf("unexpected error converting numeric: %v", err)
		}
		if back != tc {
			t.Fatalf("expected %d, got %d", tc, back)
		}

		safeBack := pgconv.NumericToInt64Safe(num)
		if safeBack != tc {
			t.Fatalf("expected safe %d, got %d", tc, safeBack)
		}
	}

	// Test invalid / null numeric
	invalidNum := pgtype.Numeric{Valid: false}
	val, err := pgconv.NumericToInt64(invalidNum)
	if err != nil || val != 0 {
		t.Fatalf("expected 0 and no error for invalid numeric")
	}
}
