package detect_js_changes

import (
	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/sergi/go-diff/diffmatchpatch"
	"strings"
)

func beautify(src string) *string {
	options := jsbeautifier.DefaultOptions()
	return jsbeautifier.BeautifyFile(src, options)
}

func lineDiff(src1, src2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(src1, src2)
	diffs := dmp.DiffMain(a, b, false)
	result := dmp.DiffCharsToLines(diffs, c)
	return result
}

func isChange(diff diffmatchpatch.Diff, ignoreKeywords []string) bool {
	switch diff.Type {
	case diffmatchpatch.DiffEqual:
		return false
	default:
		result := true
		for _, keyword := range ignoreKeywords {
			if strings.Contains(diff.Text, keyword) {
				result = false
			}
		}
		return result
	}
}

func Detect(file1 string, file2 string, ignoreKeywords []string) bool {
	src1 := beautify(file1)
	src2 := beautify(file2)
	diffs := lineDiff(*src1, *src2)
	for _, diff := range diffs {
		switch isChange(diff, ignoreKeywords) {
		case true:
			return true
		case false:
			continue
		}
	}
	return false
}
