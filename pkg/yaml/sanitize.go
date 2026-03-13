package yamlsanitize

// Sanitize attempts to fix common YAML errors heuristically, iterating until
// the tree is clean or no more fixes can be applied.
func Sanitize(src string, opts ...Option) Result {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	original := src
	var allFixes []Fix

	// Capture original parse state before any fixes.
	origTreeText, origErrors, _ := ParseTree(original)
	origLintIssues := Lint(original)

	// Run up to maxIterations fix iterations.
	for iter := 0; iter < cfg.maxIterations; iter++ {
		treeText, errors, err := ParseTree(src)
		if err != nil {
			break
		}
		lintIssues := Lint(src)

		if len(errors) == 0 && len(lintIssues) == 0 {
			return Result{
				Original:           original,
				Sanitized:          src,
				TreeText:           treeText,
				OriginalTreeText:   origTreeText,
				Errors:             errors,
				OriginalErrors:     origErrors,
				LintIssues:         lintIssues,
				OriginalLintIssues: origLintIssues,
				Fixes:              allFixes,
				ParseClean:         true,
				LintClean:          true,
			}
		}

		fixed, fixes := applyFixes(src, errors, lintIssues, &cfg)
		if len(fixes) == 0 {
			// No progress — stop.
			treeText2, errors2, _ := ParseTree(src)
			lintIssues2 := Lint(src)
			return Result{
				Original:           original,
				Sanitized:          src,
				TreeText:           treeText2,
				OriginalTreeText:   origTreeText,
				Errors:             errors2,
				OriginalErrors:     origErrors,
				LintIssues:         lintIssues2,
				OriginalLintIssues: origLintIssues,
				Fixes:              allFixes,
				ParseClean:         len(errors2) == 0,
				LintClean:          len(lintIssues2) == 0,
			}
		}
		allFixes = append(allFixes, fixes...)
		src = fixed
	}

	treeText, errors, _ := ParseTree(src)
	lintIssues := Lint(src)
	return Result{
		Original:           original,
		Sanitized:          src,
		TreeText:           treeText,
		OriginalTreeText:   origTreeText,
		Errors:             errors,
		OriginalErrors:     origErrors,
		LintIssues:         lintIssues,
		OriginalLintIssues: origLintIssues,
		Fixes:              allFixes,
		ParseClean:         len(errors) == 0,
		LintClean:          len(lintIssues) == 0,
	}
}
