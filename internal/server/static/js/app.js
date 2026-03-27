const editor = document.getElementById('editor');
const formatSelect = document.getElementById('format-select');
const exampleSelect = document.getElementById('example-select');
const exampleMeta = document.getElementById('example-meta');
const statusBadge = document.getElementById('status-badge');
const spinner = document.getElementById('spinner');
const grammarTag = document.getElementById('grammar-tag');
const inputTitle = document.getElementById('input-title');
const inputFormatBadge = document.getElementById('input-format-badge');
const treeOutput = document.getElementById('tree-output');
const treeBadge = document.getElementById('tree-badge');
const strictBadge = document.getElementById('strict-badge');
const treeToggle = document.getElementById('tree-toggle');
const inputBadge = document.getElementById('input-badge');
const issuesSection = document.getElementById('issues-section');
const issuesHeading = document.getElementById('issues-heading');
const issueItems = document.getElementById('issue-items');
const sanitizedTitle = document.getElementById('sanitized-title');
const sanitizedOut = document.getElementById('sanitized-output');
const sanitizedBadge = document.getElementById('sanitized-badge');
const fixesSection = document.getElementById('fixes-section');
const fixItems = document.getElementById('fix-items');
const copyBtn = document.getElementById('copy-btn');

const formatConfig = {
  yaml: {
    label: 'YAML',
    placeholder: 'Paste or type broken YAML here…',
    emptyExample: '— type your own YAML —',
    grammar: 'tree-sitter-yaml',
  },
  json: {
    label: 'JSON',
    placeholder: 'Paste or type malformed JSON or prose-wrapped JSON here…',
    emptyExample: '— type your own JSON —',
    grammar: 'tree-sitter-json',
  },
};

let debounceTimer = null;
let lastSanitized = '';
let lastResult = null;
let showingOriginalTree = true;
let currentFormat = 'yaml';
let allExamples = [];

async function loadExamples() {
  try {
    const res = await fetch('/api/examples');
    allExamples = await res.json();
    populateExampleSelect();
  } catch (e) {
    console.error('Failed to load examples', e);
  }
}

function examplesForCurrentFormat() {
  return allExamples.filter((ex) => ex.format === currentFormat);
}

function populateExampleSelect() {
  const config = formatConfig[currentFormat];
  const examples = examplesForCurrentFormat();

  exampleSelect.innerHTML = '';
  const emptyOpt = document.createElement('option');
  emptyOpt.value = '';
  emptyOpt.textContent = config.emptyExample;
  exampleSelect.appendChild(emptyOpt);

  examples.forEach((ex, i) => {
    const opt = document.createElement('option');
    opt.value = String(i);
    opt.textContent = ex.name;
    opt.title = ex.description || '';
    exampleSelect.appendChild(opt);
  });

  exampleSelect.value = '';
  setExampleMeta(null);
}

function setFormat(format) {
  currentFormat = format in formatConfig ? format : 'yaml';
  const config = formatConfig[currentFormat];

  formatSelect.value = currentFormat;
  grammarTag.textContent = config.grammar;
  editor.placeholder = config.placeholder;
  inputTitle.textContent = `Input ${config.label}`;
  sanitizedTitle.textContent = `Sanitized ${config.label}`;
  inputFormatBadge.textContent = currentFormat;

  populateExampleSelect();
  resetUI();
}

function setExampleMeta(example) {
  if (!example) {
    exampleMeta.textContent = 'custom input';
    exampleMeta.className = 'meta-pill';
    return;
  }

  const parts = [example.source || 'builtin'];
  if (example.category) {
    parts.push(example.category);
  }
  if (example.filename) {
    parts.push(example.filename);
  }
  exampleMeta.textContent = parts.join(' · ');
  exampleMeta.className = `meta-pill format-${example.format}`;
}

formatSelect.addEventListener('change', () => {
  setFormat(formatSelect.value);
  if (editor.value.trim()) {
    triggerAnalysis();
  }
});

exampleSelect.addEventListener('change', () => {
  const idx = exampleSelect.value;
  if (idx === '') {
    setExampleMeta(null);
    if (!editor.value.trim()) {
      resetUI();
    }
    return;
  }

  const ex = examplesForCurrentFormat()[parseInt(idx, 10)];
  if (!ex) {
    return;
  }

  editor.value = ex.input;
  setExampleMeta(ex);
  triggerAnalysis();
});

editor.addEventListener('input', () => {
  exampleSelect.value = '';
  setExampleMeta(null);
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(triggerAnalysis, 320);
});

async function triggerAnalysis() {
  const src = editor.value;
  if (!src.trim()) {
    resetUI();
    return;
  }

  setLoading(true);
  try {
    const res = await fetch('/api/sanitize', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ format: currentFormat, input: src }),
    });
    if (!res.ok) {
      throw new Error(await res.text());
    }
    renderResult(await res.json());
  } catch (e) {
    lastSanitized = '';
    lastResult = null;
    treeOutput.textContent = 'Error: ' + e.message;
    strictBadge.className = 'badge badge-quiet';
    strictBadge.textContent = currentFormat === 'json' ? 'strict invalid' : 'strict n/a';
  } finally {
    setLoading(false);
  }
}

function renderResult(data) {
  lastSanitized = data.sanitized || '';
  lastResult = data;

  const origErrors = data.original_errors || [];
  const origLintIssues = data.original_lint_issues || [];
  const origParserLintIssues = origLintIssues.filter((issue) => issue.source === 'strict-parser');
  const origHeuristicLintIssues = origLintIssues.filter((issue) => issue.source !== 'strict-parser');
  const origStrictIssues = data.original_strict_parse_clean === false ? Math.max(origParserLintIssues.length, 1) : 0;

  updateStatusBadge(origErrors.length, origHeuristicLintIssues.length, origStrictIssues);

  inputBadge.className = 'badge ' + badgeClass(origErrors.length, origLintIssues.length);
  inputBadge.textContent = summarizeInputState(origErrors.length, origLintIssues.length, origStrictIssues);

  showingOriginalTree = true;
  updateTreePanel(data);
  renderIssues(origErrors, origLintIssues);
  renderSanitizedOutput(data.sanitized);
  updateSanitizedBadge(data);
  renderFixes(data.fixes || []);
}

function updateStatusBadge(parseErrors, lintIssues, strictIssues) {
  if (currentFormat === 'json' && strictIssues > 0) {
    statusBadge.className = 'status-badge errors';
    statusBadge.textContent = strictIssues + ' strict parse issue' + (strictIssues !== 1 ? 's' : '');
    return;
  }
  if (parseErrors > 0) {
    statusBadge.className = 'status-badge errors';
    statusBadge.textContent = parseErrors + ' parse error' + (parseErrors !== 1 ? 's' : '');
    return;
  }
  if (lintIssues > 0) {
    statusBadge.className = 'status-badge lint';
    statusBadge.textContent = lintIssues + ' lint issue' + (lintIssues !== 1 ? 's' : '');
    return;
  }
  statusBadge.className = 'status-badge clean';
  statusBadge.textContent = currentFormat === 'json' ? 'strict clean' : 'clean';
}

function summarizeInputState(parseErrors, lintIssues, strictIssues) {
  if (currentFormat === 'json' && strictIssues > 0) {
    return strictIssues + ' strict';
  }
  if (parseErrors > 0) {
    return parseErrors + ' error' + (parseErrors !== 1 ? 's' : '');
  }
  if (lintIssues > 0) {
    return lintIssues + ' lint';
  }
  return 'OK';
}

function badgeClass(parseErrors, lintIssues) {
  if (parseErrors > 0) {
    return 'has-errors';
  }
  if (lintIssues > 0) {
    return 'has-lint';
  }
  return 'no-errors';
}

function updateTreePanel(data) {
  const active = data || lastResult;
  if (!active) {
    return;
  }

  const treeText = showingOriginalTree
    ? (active.original_tree_text || active.tree_text || '')
    : (active.tree_text || '');
  const errors = showingOriginalTree
    ? (active.original_errors || [])
    : (active.errors || []);

  renderTree(treeText);

  treeBadge.className = 'badge ' + (errors.length > 0 ? 'has-errors' : 'no-errors');
  treeBadge.textContent = errors.length > 0
    ? errors.length + ' ERROR' + (errors.length !== 1 ? 's' : '')
    : 'No errors';

  updateStrictBadge(active, showingOriginalTree);

  if (showingOriginalTree) {
    treeToggle.textContent = 'Show sanitized';
    treeToggle.classList.remove('active');
  } else {
    treeToggle.textContent = 'Show original';
    treeToggle.classList.add('active');
  }
}

function updateStrictBadge(data, originalView) {
  if (currentFormat !== 'json') {
    strictBadge.className = 'badge badge-quiet';
    strictBadge.textContent = 'strict n/a';
    return;
  }

  const strictClean = originalView
    ? data.original_strict_parse_clean !== false
    : Boolean(data.strict_parse_clean);

  strictBadge.className = 'badge ' + (strictClean ? 'no-errors' : 'has-errors');
  strictBadge.textContent = strictClean ? 'strict valid' : 'strict invalid';
}

treeToggle.addEventListener('click', () => {
  showingOriginalTree = !showingOriginalTree;
  updateTreePanel();
});

function prettyTree(sexp) {
  let out = '';
  let depth = 0;
  let i = 0;
  const ind = () => '  '.repeat(depth);
  const endsWithNewline = () => {
    for (let k = out.length - 1; k >= 0; k--) {
      if (out[k] === '\n') {
        return true;
      }
      if (out[k] !== ' ') {
        return false;
      }
    }
    return out.length === 0;
  };

  while (i < sexp.length) {
    const ch = sexp[i];
    if (ch === '(') {
      if (out.length > 0 && !endsWithNewline()) {
        out += '\n';
      }
      out += ind() + ch;
      depth++;
      i++;
      continue;
    }
    if (ch === ')') {
      depth--;
      out += ch;
      i++;
      continue;
    }
    if (ch === ' ') {
      if (i + 1 < sexp.length && sexp[i + 1] === '(') {
        out += '\n';
        i++;
        continue;
      }
      out += ch;
      i++;
      continue;
    }
    out += ch;
    i++;
  }
  return out;
}

function renderTree(treeText) {
  if (!treeText) {
    treeOutput.innerHTML = '<span class="empty">No tree available.</span>';
    return;
  }

  const pretty = prettyTree(treeText);
  const html = pretty
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/[()]/g, (m) => `<span class="tn-paren">${m}</span>`)
    .replace(/(<span class="tn-paren">\(<\/span>)(ERROR\b)/g, '$1<span class="tn-error">$2</span>')
    .replace(/(<span class="tn-paren">\(<\/span>)(MISSING\b)/g, '$1<span class="tn-missing">$2</span>')
    .replace(/(<span class="tn-paren">\(<\/span>)([a-z][a-z0-9_]*)/g, '$1<span class="tn-node">$2</span>')
    .replace(/\b([a-z][a-z0-9_]*):/g, '<span class="tn-field">$1</span>:')
    .replace(/&quot;([^&]*)&quot;/g, '<span class="tn-string">&quot;$1&quot;</span>');

  treeOutput.innerHTML = html;
}

function renderIssues(errors, lintIssues) {
  const total = errors.length + lintIssues.length;
  if (total === 0) {
    issuesSection.style.display = 'none';
    return;
  }

  issuesSection.style.display = '';
  issuesSection.className = 'issues-section' + (errors.length === 0 ? ' lint-only' : '');

  const sourceCounts = new Map();
  lintIssues.forEach((issue) => {
    const source = issue.source || 'heuristic';
    sourceCounts.set(source, (sourceCounts.get(source) || 0) + 1);
  });

  const parts = [];
  if (errors.length > 0) {
    parts.push(errors.length + ' parse node' + (errors.length !== 1 ? 's' : ''));
  }
  sourceCounts.forEach((count, source) => {
    parts.push(count + ' ' + source);
  });
  issuesHeading.textContent = parts.join(' · ');

  issueItems.innerHTML = '';
  errors.forEach((error) => {
    const div = document.createElement('div');
    div.className = 'issue-item';
    const loc = `L${error.start_row + 1}:${error.start_col + 1}`;
    const text = error.text ? ` <em>${escHtml(error.text.slice(0, 50))}</em>` : '';
    div.innerHTML =
      `<span class="loc">${loc}</span>` +
      `<span class="type-tag ${String(error.type).toLowerCase()}">${escHtml(error.type)}</span>` +
      text;
    issueItems.appendChild(div);
  });

  lintIssues.forEach((issue) => {
    const div = document.createElement('div');
    const source = issue.source || 'heuristic';
    div.className = 'issue-item';
    div.innerHTML =
      `<span class="loc">L${issue.start_row + 1}:${issue.start_col + 1}</span>` +
      `<span class="type-tag lint">${escHtml(source)}</span>` +
      `<span class="fix-rule lint-rule">${escHtml(issue.rule)}</span>` +
      escHtml(issue.description);
    issueItems.appendChild(div);
  });
}

function renderSanitizedOutput(output) {
  sanitizedOut.innerHTML = '';
  if (!output) {
    sanitizedOut.innerHTML = '<span class="empty">Sanitized output will appear here.</span>';
    return;
  }

  const pre = document.createElement('pre');
  pre.style.cssText = 'padding:12px;font-family:var(--mono);font-size:13px;line-height:1.65;';
  pre.textContent = output;
  sanitizedOut.appendChild(pre);
}

function updateSanitizedBadge(data) {
  const parseReady = Boolean(data.parse_clean);
  const lintReady = Boolean(data.lint_clean);
  const strictReady = currentFormat === 'json' ? Boolean(data.strict_parse_clean) : true;

  if (parseReady && lintReady && strictReady) {
    sanitizedBadge.className = 'badge no-errors';
    sanitizedBadge.textContent = currentFormat === 'json' ? 'strict clean' : 'clean';
    return;
  }

  if (!strictReady) {
    sanitizedBadge.className = 'badge has-errors';
    sanitizedBadge.textContent = 'strict invalid';
    return;
  }

  if (!parseReady) {
    sanitizedBadge.className = 'badge has-errors';
    sanitizedBadge.textContent = 'parse errors';
    return;
  }

  sanitizedBadge.className = 'badge has-lint';
  sanitizedBadge.textContent = 'lint only';
}

function renderFixes(fixes) {
  if (!fixes || fixes.length === 0) {
    fixesSection.style.display = 'none';
    return;
  }

  fixesSection.style.display = '';
  fixItems.innerHTML = '';
  fixes.forEach((fix) => {
    const div = document.createElement('div');
    div.className = 'fix-item';
    div.innerHTML =
      `<span class="fix-rule">${escHtml(fix.rule)}</span>` +
      `<span class="fix-desc">${escHtml(fix.description)}</span>` +
      `<div class="fix-diff">` +
      `<div class="fix-before">${escHtml(fix.before)}</div>` +
      `<div class="fix-after">${escHtml(fix.after)}</div>` +
      `</div>`;
    fixItems.appendChild(div);
  });
}

function escHtml(s) {
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function setLoading(on) {
  spinner.classList.toggle('hidden', !on);
}

function resetUI() {
  statusBadge.className = 'status-badge idle';
  statusBadge.textContent = 'idle';
  treeOutput.innerHTML = '<span class="empty">Tree will appear here.</span>';
  sanitizedOut.innerHTML = '<span class="empty">Sanitized output will appear here.</span>';
  issuesSection.style.display = 'none';
  fixesSection.style.display = 'none';
  treeBadge.className = 'badge';
  treeBadge.textContent = '\u2014';
  inputBadge.className = 'badge';
  inputBadge.textContent = '\u2014';
  sanitizedBadge.className = 'badge';
  sanitizedBadge.textContent = '\u2014';
  strictBadge.className = 'badge badge-quiet';
  strictBadge.textContent = currentFormat === 'json' ? 'strict pending' : 'strict n/a';
  treeToggle.textContent = 'Show sanitized';
  treeToggle.classList.remove('active');
  lastSanitized = '';
  lastResult = null;
  showingOriginalTree = true;
}

copyBtn.addEventListener('click', () => {
  if (!lastSanitized) {
    return;
  }
  navigator.clipboard.writeText(lastSanitized).then(() => {
    copyBtn.textContent = 'Copied!';
    copyBtn.classList.add('copied');
    setTimeout(() => {
      copyBtn.textContent = 'Copy';
      copyBtn.classList.remove('copied');
    }, 1500);
  });
});

setFormat('yaml');
loadExamples();
