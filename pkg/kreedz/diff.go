package kreedz

import "github.com/sergi/go-diff/diffmatchpatch"

func Diff(src, dst string) (diffs []diffmatchpatch.Diff) {
	dmp := diffmatchpatch.New()
	wSrc, wDst, warray := dmp.DiffLinesToRunes(src, dst)
	diffs = dmp.DiffMainRunes(wSrc, wDst, false)
	diffs = dmp.DiffCharsToLines(diffs, warray)
	return diffs
}