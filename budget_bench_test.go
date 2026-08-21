package runesafe_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/cplieger/runesafe/v2"
)

// This file exists because Budget and the Budgeted pair are sold entirely on a
// COST property. Their contract is not a different output — for a value carrying
// no unsafe rune they are byte-identical to the Capped pair, cut for cut — it is
// that the cap sits AHEAD of the sanitizer, so walking a multi-megabyte upstream
// value inside a memory-limited process (CWE-400 work amplification) stops being
// possible. A library whose whole product is a bound on what untrusted input can
// cost should measure that bound rather than assert it in a doc comment, and
// nothing in the repo did until this landed.
//
// So the checks here are about SCALE, not about single numbers:
//
//   - The work bound: the budgeted order's cost must not grow with the size of
//     the raw value, because the caller does not control that size. Asserted on
//     allocated BYTES across inputs three orders of magnitude apart, with the
//     pre-sanitize order's own numbers logged beside it as the contrast. Bytes
//     rather than allocation counts, because the count is nearly blind to this:
//     see allocBytesPerRun.
//   - The value-count bound: once the allowance is spent, further writes must cost
//     nothing at all, whatever they contain. Budget.Write documents this ("further
//     writes append nothing") and it is what stops a hostile piece COUNT
//     amplifying an aggregate the per-value bound already covers.
//
// One thing these measurements show that the doc comments do not claim, and that a
// reader could easily assume the other way round: for a handful of small HONEST
// values the streaming Budget is the DEARER shape. Measured over eight 16-byte
// values under a 200-byte budget, streaming costs 7 allocations and 544 bytes where
// joining them and capping once costs 2 and 304 — the aggregate is assembled
// through a Builder and every write maps its own chunk. The two shapes cross over
// as the values grow, and by 64 KiB each the same comparison is 1.5 us against
// 2.5 ms. So the Budget is not an optimization for the honest case; it is a bound
// on the hostile one, and BenchmarkAggregate is where both halves of that show.

// benchGroups is the aggregate fixture: several upstream-controlled values headed
// for ONE log attribute under ONE shared byte budget, which is the shape Budget
// exists for. Two of them carry unsafe runes, and one of those is FIRST so that the
// measured path includes a rewrite rather than only the sanitizer's no-change fast
// path — with a clean value leading, the whole first write is absorbed by that fast
// path and the series would describe the wrong work.
func benchGroups(each int) []string {
	units := []string{
		"Group\x1bBeta-",
		"Group-Alpha-",
		"Grüppe-Gamma-",
		"Group\u202eDelta-",
		"Group-Epsilon-",
		"Group-Zeta-",
		"Group-Eta-",
		"Group-Theta-",
	}
	groups := make([]string, len(units))
	for i, u := range units {
		groups[i] = repeatTo(u, each)
	}
	return groups
}

// writeGroups runs the loop the README documents for Budget: the separator goes
// through Write too, and the caller stops when a write is REFUSED rather than when
// the budget looks full, because the refusal is what latches the truncation fact.
// BenchmarkAggregate and TestBudgetCostIsSetByTheBudgetNotTheValues share it so the
// benchmark and the assertion measure the same call shape.
// TestBudgetSpentAllowanceChargesNothingPerValue deliberately does not use it, for
// the reason given there: this loop's early exit would stand in for the guard that
// test is about.
func writeGroups(b *runesafe.Budget, groups []string) (string, bool) {
	for i, g := range groups {
		if i > 0 && !b.Write(", ") {
			break
		}
		if !b.Write(g) {
			break
		}
	}
	return b.Result()
}

// allocBytesPerRun is testing.AllocsPerRun's missing twin: mean BYTES allocated per
// call rather than mean allocations. It exists because the work bound this file
// gates is a bytes property, and the allocation COUNT is a poor witness for it —
// deleting Write's pre-cap entirely leaves the count unchanged for every input class
// but one, since a value the sanitizer merely rewrites needs exactly one buffer
// whatever its size. Bytes see it immediately: 368 against four million.
//
// runtime.MemStats.TotalAlloc is cumulative and monotonic, so the delta over a known
// number of calls is the mean directly. Background allocation from the runtime can
// add to it, which is why every assertion built on this compares magnitudes with
// wide headroom rather than exact equality; the differences it is asked to resolve
// are four orders of magnitude wide.
func allocBytesPerRun(runs int, f func()) float64 {
	f() // Warm up, so first-call lazy initialization is not charged to the mean.
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	for range runs {
		f()
	}
	runtime.ReadMemStats(&after)
	return float64(after.TotalAlloc-before.TotalAlloc) / float64(runs)
}

// TestSanitizeBudgetedWorkBoundDoesNotScaleWithInput gates the reason the Budgeted
// pair exists at all: it caps the RAW bytes on a rune boundary before the sanitizer
// walks them, so its cost is a function of the CAP and not of the value an upstream
// chose to send.
//
// The assertion is on allocated BYTES, because that is what the claim is about and
// the allocation count barely moves — see allocBytesPerRun. Measured on go1.27.0
// with a 64-byte cap over invalid UTF-8: the budgeted order allocates 368 bytes per
// call at every raw size from 1 KiB to 1 MiB, while the capped order runs about
// 3.1 KB, then 98 KB, then 4.0 MB over the same three sizes, because it materializes
// the sanitized form of everything it was handed. Deleting the pre-cap puts the
// budgeted order on exactly the capped order's numbers, which is how this assertion
// was seen to fail. That is the divergence the type's doc comment describes, in
// numbers.
//
// The tolerance is a factor of two across a 1024x increase in input, not equality:
// TotalAlloc picks up whatever else the runtime did during the loop, and a bound
// that loose still separates a flat cost from one tracking the payload by three
// orders of magnitude. The count is checked for constancy as well, but as the weaker
// of the two — invalid UTF-8 is the only class where it moves at all.
//
// The class under test is invalid UTF-8 because it is where the two orders separate
// furthest: sanitizing grows each invalid byte into the three-byte U+FFFD, so the
// pre-sanitize order pays for a buffer three times the size of the whole hostile
// payload while this order never materializes more than the cap.
func TestSanitizeBudgetedWorkBoundDoesNotScaleWithInput(t *testing.T) {
	const (
		runs      = 5
		limit     = 64
		tolerance = 2.0
	)
	sizes := []int{1024, 32768, 1048576}
	budgetedBytes := make([]float64, 0, len(sizes))
	cappedBytes := make([]float64, 0, len(sizes))
	budgetedAllocs := make([]float64, 0, len(sizes))
	for _, n := range sizes {
		in := invalidUTF8(n)
		budgeted := func() {
			sinkString, sinkBool = runesafe.SanitizeSingleLineBudgeted(in, limit, "...")
		}
		budgetedBytes = append(budgetedBytes, allocBytesPerRun(runs, budgeted))
		budgetedAllocs = append(budgetedAllocs, testing.AllocsPerRun(runs, budgeted))
		cappedBytes = append(cappedBytes, allocBytesPerRun(runs, func() {
			sinkString, sinkBool = runesafe.SanitizeSingleLineCapped(in, limit, "...")
		}))
	}
	for i, got := range budgetedBytes {
		if got > budgetedBytes[0]*tolerance {
			t.Errorf("SanitizeSingleLineBudgeted(invalid UTF-8, %d bytes, %d, \"...\") "+
				"allocated %.0f bytes per run but %.0f at %d bytes, want no more than %.0fx: "+
				"the cap sits ahead of the sanitizer, so the raw value's size must not "+
				"reach the cost", sizes[i], limit, got, budgetedBytes[0], sizes[0], tolerance)
		}
		if budgetedAllocs[i] != budgetedAllocs[0] {
			t.Errorf("SanitizeSingleLineBudgeted(invalid UTF-8, %d bytes, %d, \"...\") "+
				"allocated %v times per run but %v at %d bytes, want the same count: a "+
				"bounded chunk cannot need more buffer regrowths than a smaller one",
				sizes[i], limit, budgetedAllocs[i], budgetedAllocs[0], sizes[0])
		}
	}
	t.Logf("cap %d over %v bytes of raw invalid UTF-8: budgeted %.0f bytes/run, capped "+
		"%.0f — the capped order pays for the whole value", limit, sizes, budgetedBytes, cappedBytes)
}

// TestBudgetSpentAllowanceChargesNothingPerValue gates the aggregate's other
// bound. Budget.Write's contract is that once anything has been dropped, or the
// allowance is spent, a further write appends nothing and is refused whole; the
// cost consequence is that the sanitizer is never even entered for it. That is
// what makes a hostile piece COUNT harmless: an upstream returning 8000 groups
// instead of 8 must not cost 1000x more.
//
// Two things about the shape of this test were found by breaking the guard and
// watching it stay green, and both are the reason it is written this way.
//
// It writes every value unconditionally instead of going through writeGroups.
// The documented loop stops at the first refusal, so a test driven through it
// measures the loop's early exit rather than Write's guard — with the guard
// removed it still reported a constant. The guard exists for the caller that does
// NOT stop, which Write's own doc comment describes as the one that must attempt
// the write it cannot fit, so that is the caller this test models.
//
// And the values are 1 KiB each rather than the 20 bytes of the first draft.
// Sanitizing a 20-byte value produces a buffer the compiler can prove does not
// escape, so it lands on the stack, AllocsPerRun sees nothing, and a refused write
// doing full work measures as free. Only above the stack-allocation threshold does
// the regression become visible.
//
// Measured on go1.27.0 with the guard intact: a constant 3 allocations for a
// 32-byte budget whether 8 or 8000 values are written into it. With the sanitize
// call moved above the guard — an easy edit, and one no functional test would
// notice, since the OUTPUT is identical either way — it becomes 10, 82, 802 and 8002:
// one allocation per value.
func TestBudgetSpentAllowanceChargesNothingPerValue(t *testing.T) {
	const runs = 20
	writeCounts := []int{8, 80, 800, 8000}
	counts := make([]float64, 0, len(writeCounts))
	value := repeatTo("Group\x1bName\u202eSuffix ", 1024)
	for _, writes := range writeCounts {
		counts = append(counts, testing.AllocsPerRun(runs, func() {
			b := runesafe.NewSingleLineBudget(32, "...")
			for range writes {
				sinkBool = b.Write(value) // deliberately ignoring the refusal
			}
			sinkString, sinkBool = b.Result()
		}))
	}
	for i, got := range counts {
		if got != counts[0] {
			t.Errorf("Budget(32, \"...\") over %d values of 1 KiB each allocated %v times per run "+
				"but %v over %d values, want the same count: a spent allowance must refuse "+
				"a value without sanitizing it, so the piece count cannot amplify the cost",
				writeCounts[i], got, counts[0], writeCounts[0])
		}
	}
	t.Logf("a spent 32-byte budget is a constant %v allocations from %d to %d written values",
		counts[0], writeCounts[0], writeCounts[len(writeCounts)-1])
}

// TestBudgetCostIsSetByTheBudgetNotTheValues is the aggregate twin of the
// work-bound test above: eight writes into one budget, with each value three orders
// of magnitude larger, must cost the same. Every Write caps its argument to the
// REMAINING budget before mapping it, so what the sanitizer sees is bounded by the
// budget and not by the values — the property that makes it safe to pass an
// upstream list straight into a Budget without measuring it first.
//
// Bytes again rather than counts, for the reason on allocBytesPerRun: with Write's
// pre-cap deleted this test reported a constant allocation count for a fixture of
// clean values, and a constant count for one the sanitizer merely rewrites, because
// a rewrite is one buffer whatever its size. Only the byte figure moves with the
// input, and it moves by four orders of magnitude.
//
// Measured on go1.27.0 with a 200-byte budget over eight invalid-UTF-8 values: 1,040
// bytes per call whether each value is 1 KiB or 1 MiB. Every size here exceeds the
// budget, so the first write fills it and the remaining seven are refused outright,
// which is why the number does not move at all. The same call over eight 16-byte
// values that all fit costs 552 bytes in 12 allocations — fewer bytes but four times
// as many of them, because the aggregate is grown across all eight writes instead of
// being filled by the first. Cost tracks the budget the caller chose either way, and
// an upstream cannot raise it by sending more.
func TestBudgetCostIsSetByTheBudgetNotTheValues(t *testing.T) {
	const (
		runs      = 10
		limit     = 200
		tolerance = 2.0
	)
	// Eight invalid-UTF-8 values per case; see above for why the class matters.
	values := func(each int) []string {
		out := make([]string, 8)
		for i := range out {
			out[i] = invalidUTF8(each)
		}
		return out
	}
	sizes := []int{1024, 32768, 1048576}
	bytesPerRun := make([]float64, 0, len(sizes))
	for _, each := range sizes {
		groups := values(each)
		bytesPerRun = append(bytesPerRun, allocBytesPerRun(runs, func() {
			sinkString, sinkBool = writeGroups(runesafe.NewSingleLineBudget(limit, "..."), groups)
		}))
	}
	for i, got := range bytesPerRun {
		if got > bytesPerRun[0]*tolerance {
			t.Errorf("Budget(%d, \"...\") over 8 values of %d bytes each allocated %.0f bytes "+
				"per run but %.0f at %d bytes each, want no more than %.0fx: each write caps "+
				"to the remaining budget before sanitizing, so value size must not reach the "+
				"cost", limit, sizes[i], got, bytesPerRun[0], sizes[0], tolerance)
		}
	}
	fitting := values(16)
	t.Logf("a %d-byte budget over 8 values allocates %.0f bytes/run at %v bytes each, "+
		"against %.0f when all eight fit at 16 bytes each",
		limit, bytesPerRun, sizes,
		allocBytesPerRun(runs, func() {
			sinkString, sinkBool = writeGroups(runesafe.NewSingleLineBudget(limit, "..."), fitting)
		}))
}

// BenchmarkWorkBound is the series that makes the pre-cap order's whole reason for
// existing visible: the same cap, the same marker, the same input, in the two
// orders the package offers. capped sanitizes the value and then bounds the
// result; budgeted bounds the raw bytes and then sanitizes the chunk.
//
// Read it as two slopes, not four numbers. budgeted should be FLAT across the two
// sizes — it never walks past the cap — while capped should scale with the input,
// and its B/op series (which the fleet's reducer publishes and alerts on
// separately) is where the difference is stark rather than marginal. A change that
// moved the cap back behind the sanitizer would flatten the pair together, which is
// the regression this exists to catch.
//
// The fixture is invalid UTF-8 because it is the worst case for the capped order
// (three output bytes per input byte) and the realistic one for a truncated
// upstream response. No SetBytes: budgeted deliberately does not process the bytes
// it is handed, so a throughput figure derived from the fixture's length would
// describe the fixture and not the work.
func BenchmarkWorkBound(b *testing.B) {
	const (
		limit  = 64
		marker = "..."
	)
	for _, order := range []struct {
		name string
		call func(string) (string, bool)
	}{
		{"capped", func(s string) (string, bool) {
			return runesafe.SanitizeSingleLineCapped(s, limit, marker)
		}},
		{"budgeted", func(s string) (string, bool) {
			return runesafe.SanitizeSingleLineBudgeted(s, limit, marker)
		}},
	} {
		b.Run(order.name, func(b *testing.B) {
			for _, n := range []int{typicalBytes, 1048576} {
				in := invalidUTF8(n)
				b.Run(fmt.Sprintf("bytes_%d", n), func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						sinkString, sinkBool = order.call(in)
					}
				})
			}
		})
	}
}

// BenchmarkAggregate puts the streaming Budget beside the one-shot call it
// replaces — join the untrusted values, then sanitize and cap the result — so both
// shapes land in the same chart at the same fixture.
//
// The honest reading is that streaming wins here by three orders of magnitude —
// 1.5 us against 2.5 ms, 624 bytes against 533 KB — because join_then_capped
// materializes the whole untrusted aggregate before any bound applies, while every
// Write caps its argument to what is left of the budget first. That gap is the
// type's entire reason for existing, and the two series are meant to be read
// together: streaming's numbers do not move when the values grow, join_then_capped's
// track whatever the upstream decided to send.
//
// The fixture sits deliberately at the hostile end. The file's header comment has
// the other end, where the same comparison over eight 16-byte values goes the other
// way and streaming is the dearer shape; publishing only this end would misrepresent
// what adopting the type costs.
func BenchmarkAggregate(b *testing.B) {
	const (
		limit  = 200
		marker = "..."
	)
	groups := benchGroups(largeBytes)
	joined := strings.Join(groups, ", ")
	b.Run("streaming", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkString, sinkBool = writeGroups(runesafe.NewSingleLineBudget(limit, marker), groups)
		}
	})
	b.Run("join_then_capped", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkString, sinkBool = runesafe.SanitizeSingleLineCapped(joined, limit, marker)
		}
	})
}
