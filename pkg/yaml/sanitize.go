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
	origDoc, origErr := analyzeDocument(original)
	origTreeText, origErrors := "", []ErrorNode(nil)
	origLintIssues := []LintIssue(nil)
	if origErr == nil {
		origTreeText = origDoc.TreeText
		origErrors = origDoc.ParseErrors
		origLintIssues = lintIssuesFromAnalysis(original, origDoc)
	} else {
		origTreeText, origErrors, _ = ParseTree(original)
		origLintIssues = Lint(original)
	}

	// Run up to maxIterations fix iterations.
	for iter := 0; iter < cfg.maxIterations; iter++ {
		doc, err := analyzeDocument(src)
		if err != nil {
			break
		}
		treeText := doc.TreeText
		errors := doc.ParseErrors
		lintIssues := lintIssuesFromAnalysis(src, doc)

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

		fixed, fixes := applyFixes(src, doc, lintIssues, &cfg)
		if len(fixes) == 0 {
			// No progress — stop.
			doc2, err := analyzeDocument(src)
			treeText2, errors2, lintIssues2 := "", []ErrorNode(nil), []LintIssue(nil)
			if err == nil {
				treeText2 = doc2.TreeText
				errors2 = doc2.ParseErrors
				lintIssues2 = lintIssuesFromAnalysis(src, doc2)
			} else {
				treeText2, errors2, _ = ParseTree(src)
				lintIssues2 = Lint(src)
			}
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

	doc, err := analyzeDocument(src)
	treeText, errors, lintIssues := "", []ErrorNode(nil), []LintIssue(nil)
	if err == nil {
		treeText = doc.TreeText
		errors = doc.ParseErrors
		lintIssues = lintIssuesFromAnalysis(src, doc)
	} else {
		treeText, errors, _ = ParseTree(src)
		lintIssues = Lint(src)
	}
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
