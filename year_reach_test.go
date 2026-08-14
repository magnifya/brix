package gronx

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// yearCase describes an expected prev/next result for an expression carrying
// an explicit year segment. When want is empty an error is expected; when
// unreachable is true that error must be the "unreachable year" error.
type yearCase struct {
	expr        string
	ref         time.Time
	incl        bool
	want        string
	unreachable bool
}

func parseRef(s string) time.Time {
	ref, _ := time.Parse(FullDateFormat, s)
	return ref
}

func assertYearCase(t *testing.T, label string, got time.Time, err error, c yearCase) {
	t.Helper()
	if c.want == "" {
		if err == nil {
			t.Errorf("%s: expected an error, got %s", label, got.Format(FullDateFormat))
			return
		}
		if c.unreachable && !strings.HasPrefix(err.Error(), "unreachable year") {
			t.Errorf("%s: expected unreachable year error, got %q", label, err.Error())
		}
		return
	}
	if err != nil {
		t.Errorf("%s: expected %s, got error %v", label, c.want, err)
		return
	}
	if actual := got.Format(FullDateFormat); actual != c.want {
		t.Errorf("%s: expected %s, got %s", label, c.want, actual)
	}
}

// TestPrevTickSameYear is the regression for reverse lookups whose year segment
// matches the reference year: a matching moment earlier in that year must be
// returned instead of an "unreachable year" error, for both incl modes.
func TestPrevTickSameYear(t *testing.T) {
	ref := parseRef("2024-11-08 22:18:16")
	for _, incl := range []bool{true, false} {
		prev, err := PrevTickBefore("30 15 4 11 * 2024", ref, incl)
		assertYearCase(t, fmt.Sprintf("prev incl=%v", incl), prev, err,
			yearCase{want: "2024-11-04 15:30:00"})
	}
}

func TestPrevTickYearReachability(t *testing.T) {
	ref := parseRef("2024-11-08 22:18:16")
	cases := []yearCase{
		// same year, trigger near year start
		{"0 0 1 1 * 2024", parseRef("2024-06-15 12:00:00"), true, "2024-01-01 00:00:00", false},
		{"0 0 1 1 * 2024", parseRef("2024-06-15 12:00:00"), false, "2024-01-01 00:00:00", false},
		// same year, trigger near year end, reference after the trigger
		{"0 0 31 12 * 2024", parseRef("2024-12-31 23:59:00"), true, "2024-12-31 00:00:00", false},
		{"0 0 31 12 * 2024", parseRef("2024-12-31 23:59:00"), false, "2024-12-31 00:00:00", false},
		// year earlier than the reference year
		{"30 15 4 11 * 2020", ref, true, "2020-11-04 15:30:00", false},
		{"30 15 4 11 * 2020", ref, false, "2020-11-04 15:30:00", false},
		// year later than the reference year is unreachable backwards
		{"30 15 4 11 * 2026", ref, true, "", true},
		{"30 15 4 11 * 2026", ref, false, "", true},
		// year list: most recent past year wins
		{"30 15 4 11 * 2022,2026", ref, true, "2022-11-04 15:30:00", false},
		{"30 15 4 11 * 2022,2026", ref, false, "2022-11-04 15:30:00", false},
		// year list including the reference year
		{"30 15 4 11 * 2020,2024", ref, true, "2024-11-04 15:30:00", false},
		{"30 15 4 11 * 2020,2024", ref, false, "2024-11-04 15:30:00", false},
		// year list entirely in the future is unreachable backwards
		{"30 15 4 11 * 2026,2028", ref, true, "", true},
		{"30 15 4 11 * 2026,2028", ref, false, "", true},
		// year range including the reference year
		{"30 15 4 11 * 2020-2024", ref, true, "2024-11-04 15:30:00", false},
		{"30 15 4 11 * 2020-2024", ref, false, "2024-11-04 15:30:00", false},
		// year range entirely in the past
		{"30 15 4 11 * 2020-2022", ref, true, "2022-11-04 15:30:00", false},
		{"30 15 4 11 * 2020-2022", ref, false, "2022-11-04 15:30:00", false},
		// year range entirely in the future is unreachable backwards
		{"30 15 4 11 * 2026-2028", ref, true, "", true},
		{"30 15 4 11 * 2026-2028", ref, false, "", true},
	}
	for i, c := range cases {
		got, err := PrevTickBefore(c.expr, c.ref, c.incl)
		assertYearCase(t, fmt.Sprintf("prev #%d incl=%v %s", i, c.incl, c.expr), got, err, c)
	}
}

func TestNextTickYearReachability(t *testing.T) {
	ref := parseRef("2024-11-08 22:18:16")
	cases := []yearCase{
		// same year, a future trigger exists: must not be flagged unreachable
		{"30 15 4 12 * 2024", ref, true, "2024-12-04 15:30:00", false},
		{"30 15 4 12 * 2024", ref, false, "2024-12-04 15:30:00", false},
		// year later than the reference year
		{"30 15 4 11 * 2026", ref, true, "2026-11-04 15:30:00", false},
		{"30 15 4 11 * 2026", ref, false, "2026-11-04 15:30:00", false},
		// year earlier than the reference year is unreachable forwards
		{"30 15 4 11 * 2020", ref, true, "", true},
		{"30 15 4 11 * 2020", ref, false, "", true},
		// same year whose only trigger is in the past: no future match, must error
		{"30 15 4 11 * 2024", ref, true, "", false},
		{"30 15 4 11 * 2024", ref, false, "", false},
		// year list: past year skipped, future year picked
		{"30 15 4 11 * 2024,2026", ref, true, "2026-11-04 15:30:00", false},
		{"30 15 4 11 * 2024,2026", ref, false, "2026-11-04 15:30:00", false},
		// year range spanning past and future years
		{"30 15 4 11 * 2024-2026", ref, true, "2025-11-04 15:30:00", false},
		{"30 15 4 11 * 2024-2026", ref, false, "2025-11-04 15:30:00", false},
	}
	for i, c := range cases {
		got, err := NextTickAfter(c.expr, c.ref, c.incl)
		assertYearCase(t, fmt.Sprintf("next #%d incl=%v %s", i, c.incl, c.expr), got, err, c)
	}
}
