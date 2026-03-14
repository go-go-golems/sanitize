package jsonsanitize

// Sanitize attempts to fix common JSON errors heuristically, iterating until
// the document is clean enough or no more fixes can be applied.
func Sanitize(src string, opts ...Option) Result {
	result, err := SanitizeWithOptions(src, opts...)
	if err != nil {
		return Result{
			Original:         src,
			Sanitized:        src,
			ParseClean:       false,
			LintClean:        false,
			StrictParseClean: false,
		}
	}
	return result
}

// SanitizeWithOptions attempts to fix common JSON issues using validated options.
func SanitizeWithOptions(src string, opts ...Option) (Result, error) {
	cfg, err := buildConfig(opts...)
	if err != nil {
		return Result{}, err
	}
	return sanitizeWithConfig(src, cfg), nil
}

func sanitizeWithConfig(src string, cfg config) Result {
	original := src
	var allFixes []Fix

	origDoc, origErr := analyzeDocument(original)
	origTreeText, origErrors := "", []ErrorNode(nil)
	origLintIssues := []LintIssue(nil)
	if origErr == nil {
		origTreeText = origDoc.TreeText
		origErrors = origDoc.ParseErrors
		origLintIssues = lintIssuesFromAnalysis(origDoc, &cfg)
	} else {
		origTreeText, origErrors, _ = ParseTree(original)
		origLintIssues = Lint(original)
	}

	for iter := 0; iter < cfg.maxIterations; iter++ {
		doc, err := analyzeDocument(src)
		if err != nil {
			break
		}

		lintIssues := lintIssuesFromAnalysis(doc, &cfg)
		if len(doc.ParseErrors) == 0 && len(lintIssues) == 0 && doc.StrictParseError == nil {
			return Result{
				Original:           original,
				Sanitized:          src,
				TreeText:           doc.TreeText,
				OriginalTreeText:   origTreeText,
				Errors:             doc.ParseErrors,
				OriginalErrors:     origErrors,
				LintIssues:         lintIssues,
				OriginalLintIssues: origLintIssues,
				Fixes:              allFixes,
				ParseClean:         true,
				LintClean:          true,
				StrictParseClean:   true,
			}
		}

		fixed, fixes := applyFixes(src, doc, &cfg)
		if len(fixes) == 0 || fixed == src {
			return finalizeResult(original, src, origTreeText, origErrors, origLintIssues, allFixes, cfg)
		}

		allFixes = append(allFixes, fixes...)
		src = fixed
	}

	return finalizeResult(original, src, origTreeText, origErrors, origLintIssues, allFixes, cfg)
}

func finalizeResult(original, current, origTreeText string, origErrors []ErrorNode, origLintIssues []LintIssue, fixes []Fix, cfg config) Result {
	doc, err := analyzeDocument(current)
	treeText, errors, lintIssues := "", []ErrorNode(nil), []LintIssue(nil)
	parseClean, strictClean := false, false
	if err == nil {
		treeText = doc.TreeText
		errors = doc.ParseErrors
		lintIssues = lintIssuesFromAnalysis(doc, &cfg)
		parseClean = len(errors) == 0
		strictClean = doc.StrictParseError == nil
	} else {
		treeText, errors, _ = ParseTree(current)
		lintIssues = Lint(current)
	}

	return Result{
		Original:           original,
		Sanitized:          current,
		TreeText:           treeText,
		OriginalTreeText:   origTreeText,
		Errors:             errors,
		OriginalErrors:     origErrors,
		LintIssues:         lintIssues,
		OriginalLintIssues: origLintIssues,
		Fixes:              fixes,
		ParseClean:         parseClean,
		LintClean:          len(lintIssues) == 0,
		StrictParseClean:   strictClean,
	}
}
