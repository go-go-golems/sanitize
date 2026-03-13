/* ── DOM refs ─────────────────────────────────────────────────────────────── */
const editor         = document.getElementById('editor');
const exampleSelect  = document.getElementById('example-select');
const statusBadge    = document.getElementById('status-badge');
const spinner        = document.getElementById('spinner');
const treeOutput     = document.getElementById('tree-output');
const treeBadge      = document.getElementById('tree-badge');
const treeToggle     = document.getElementById('tree-toggle');
const inputBadge     = document.getElementById('input-badge');
const issuesSection  = document.getElementById('issues-section');
const issuesHeading  = document.getElementById('issues-heading');
const issueItems     = document.getElementById('issue-items');
const sanitizedOut   = document.getElementById('sanitized-output');
const sanitizedBadge = document.getElementById('sanitized-badge');
const fixesSection   = document.getElementById('fixes-section');
const fixItems       = document.getElementById('fix-items');
const copyBtn        = document.getElementById('copy-btn');

let debounceTimer      = null;
let lastSanitized      = '';
let lastResult         = null;
let showingOriginalTree = true;

/* ── Examples ─────────────────────────────────────────────────────────────── */
async function loadExamples() {
  try {
    const res = await fetch('/api/examples');
    const examples = await res.json();
    examples.forEach((ex, i) => {
      const opt = document.createElement('option');
      opt.value = i;
      opt.textContent = ex.name;
      opt.title = ex.description || '';
      exampleSelect.appendChild(opt);
    });
    window._examples = examples;
  } catch (e) { console.error('Failed to load examples', e); }
}

exampleSelect.addEventListener('change', () => {
  const idx = exampleSelect.value;
  if (idx === '') { resetUI(); return; }
  const ex = window._examples && window._examples[parseInt(idx)];
  if (ex) { editor.value = ex.yaml; triggerAnalysis(); }
});

/* ── Analysis ─────────────────────────────────────────────────────────────── */
editor.addEventListener('input', () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(triggerAnalysis, 320);
});

async function triggerAnalysis() {
  const src = editor.value;
  if (!src.trim()) { resetUI(); return; }
  setLoading(true);
  try {
    const res = await fetch('/api/sanitize', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ yaml: src }),
    });
    if (!res.ok) throw new Error(await res.text());
    renderResult(await res.json());
  } catch (e) {
    treeOutput.textContent = 'Error: ' + e.message;
  } finally {
    setLoading(false);
  }
}

/* ── Render ───────────────────────────────────────────────────────────────── */
function renderResult(data) {
  lastSanitized = data.sanitized || '';
  lastResult    = data;

  const origErrCount  = (data.original_errors  || []).length;
  const origLintCount = (data.original_lint_issues || []).length;

  // Status badge reflects original input state
  if (origErrCount > 0) {
    statusBadge.className = 'status-badge errors';
    statusBadge.textContent = origErrCount + ' parse error' + (origErrCount !== 1 ? 's' : '');
  } else if (origLintCount > 0) {
    statusBadge.className = 'status-badge lint';
    statusBadge.textContent = origLintCount + ' lint issue' + (origLintCount !== 1 ? 's' : '');
  } else {
    statusBadge.className = 'status-badge clean';
    statusBadge.textContent = 'clean';
  }

  // Input badge
  inputBadge.className = 'badge ' + (origErrCount > 0 ? 'has-errors' : origLintCount > 0 ? 'has-lint' : 'no-errors');
  inputBadge.textContent = origErrCount > 0
    ? origErrCount + ' error' + (origErrCount !== 1 ? 's' : '')
    : origLintCount > 0 ? origLintCount + ' lint' : 'OK';

  // Tree panel: default to original
  showingOriginalTree = true;
  updateTreePanel(data);

  // Issues: always show original
  renderIssues(data.original_errors || [], data.original_lint_issues || []);

  // Sanitized output
  sanitizedOut.innerHTML = '';
  if (data.sanitized) {
    const pre = document.createElement('pre');
    pre.style.cssText = 'padding:12px;font-family:var(--mono);font-size:13px;line-height:1.65;';
    pre.textContent = data.sanitized;
    sanitizedOut.appendChild(pre);
  }

  // Sanitized badge
  sanitizedBadge.className = 'badge ' + (data.parse_clean && data.lint_clean ? 'no-errors' : data.parse_clean ? 'has-lint' : 'has-errors');
  sanitizedBadge.textContent = data.parse_clean && data.lint_clean ? 'clean' : data.parse_clean ? 'lint only' : 'has errors';

  // Fixes
  renderFixes(data.fixes || []);
}

/* ── Tree panel ───────────────────────────────────────────────────────────── */
function updateTreePanel(data) {
  data = data || lastResult;
  if (!data) return;

  const treeText = showingOriginalTree
    ? (data.original_tree_text || data.tree_text || '')
    : (data.tree_text || '');
  const errors = showingOriginalTree
    ? (data.original_errors || [])
    : (data.errors || []);

  renderTree(treeText);

  const errCount = errors.length;
  treeBadge.className = 'badge ' + (errCount > 0 ? 'has-errors' : 'no-errors');
  treeBadge.textContent = errCount > 0
    ? errCount + ' ERROR' + (errCount !== 1 ? 's' : '')
    : 'No errors';

  if (showingOriginalTree) {
    treeToggle.textContent = 'Show sanitized';
    treeToggle.classList.remove('active');
  } else {
    treeToggle.textContent = 'Show original';
    treeToggle.classList.add('active');
  }
}

treeToggle.addEventListener('click', () => {
  showingOriginalTree = !showingOriginalTree;
  updateTreePanel();
});

/* ── Pretty-print sexp ────────────────────────────────────────────────────── */
function prettyTree(sexp) {
  let out = '', depth = 0, i = 0;
  const ind = () => '  '.repeat(depth);
  while (i < sexp.length) {
    const ch = sexp[i];
    if (ch === '(') {
      if (out.length > 0 && out[out.length - 1] !== '\n') out += '\n' + ind();
      out += ch; depth++; i++;
    } else if (ch === ')') {
      depth--; out += ch; i++;
    } else if (ch === ' ') {
      if (i + 1 < sexp.length && sexp[i + 1] === '(') {
        out += '\n' + ind(); i++;
      } else { out += ch; i++; }
    } else { out += ch; i++; }
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
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    .replace(/[()]/g, m => `<span class="tn-paren">${m}</span>`)
    .replace(/(<span class="tn-paren">\(<\/span>)(ERROR\b)/g,
      '$1<span class="tn-error">$2</span>')
    .replace(/(<span class="tn-paren">\(<\/span>)(MISSING\b)/g,
      '$1<span class="tn-missing">$2</span>')
    .replace(/(<span class="tn-paren">\(<\/span>)([a-z][a-z0-9_]*)/g,
      '$1<span class="tn-node">$2</span>')
    .replace(/\b([a-z][a-z0-9_]*)(:<\/span>)/g,
      '<span class="tn-field">$1</span>$2')
    .replace(/\b([a-z][a-z0-9_]*):/g,
      '<span class="tn-field">$1</span>:')
    .replace(/&quot;([^&]*)&quot;/g,
      '<span style="color:#888;font-style:italic;">&quot;$1&quot;</span>');
  treeOutput.innerHTML = html;
}

/* ── Issues ───────────────────────────────────────────────────────────────── */
function renderIssues(errors, lintIssues) {
  const total = errors.length + lintIssues.length;
  if (total === 0) { issuesSection.style.display = 'none'; return; }

  issuesSection.style.display = '';
  issuesSection.className = 'issues-section' + (errors.length === 0 ? ' lint-only' : '');

  const parts = [];
  if (errors.length > 0)     parts.push(errors.length + ' ERROR node' + (errors.length !== 1 ? 's' : ''));
  if (lintIssues.length > 0) parts.push(lintIssues.length + ' lint issue' + (lintIssues.length !== 1 ? 's' : ''));
  issuesHeading.textContent = parts.join(' · ');

  issueItems.innerHTML = '';
  errors.forEach(e => {
    const div = document.createElement('div');
    div.className = 'issue-item';
    const loc  = `L${e.start_row + 1}:${e.start_col + 1}`;
    const text = e.text ? ` <em>${escHtml(e.text.slice(0, 50))}</em>` : '';
    div.innerHTML = `<span class="loc">${loc}</span><span class="type-tag ${e.type.toLowerCase()}">${e.type}</span>${text}`;
    issueItems.appendChild(div);
  });
  lintIssues.forEach(li => {
    const div = document.createElement('div');
    div.className = 'issue-item';
    div.innerHTML =
      `<span class="type-tag lint">lint</span>` +
      `<span class="fix-rule" style="background:#e8f4fd;color:#1a6ea8;">${escHtml(li.rule)}</span>` +
      escHtml(li.description);
    issueItems.appendChild(div);
  });
}

/* ── Fixes ────────────────────────────────────────────────────────────────── */
function renderFixes(fixes) {
  if (!fixes || fixes.length === 0) { fixesSection.style.display = 'none'; return; }
  fixesSection.style.display = '';
  fixItems.innerHTML = '';
  fixes.forEach(f => {
    const div = document.createElement('div');
    div.className = 'fix-item';
    div.innerHTML =
      `<span class="fix-rule">${escHtml(f.rule)}</span>` +
      `<span class="fix-desc">${escHtml(f.description)}</span>` +
      `<div class="fix-diff">` +
        `<div class="fix-before">${escHtml(f.before)}</div>` +
        `<div class="fix-after">${escHtml(f.after)}</div>` +
      `</div>`;
    fixItems.appendChild(div);
  });
}

/* ── Helpers ──────────────────────────────────────────────────────────────── */
function escHtml(s) {
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}
function setLoading(on) { spinner.classList.toggle('hidden', !on); }
function resetUI() {
  statusBadge.className = 'status-badge idle'; statusBadge.textContent = 'idle';
  treeOutput.innerHTML = '<span class="empty">Tree will appear here.</span>';
  sanitizedOut.innerHTML = '<span class="empty">Sanitized output will appear here.</span>';
  issuesSection.style.display = 'none'; fixesSection.style.display = 'none';
  treeBadge.className = 'badge'; treeBadge.textContent = '\u2014';
  inputBadge.className = 'badge'; inputBadge.textContent = '\u2014';
  sanitizedBadge.className = 'badge'; sanitizedBadge.textContent = '\u2014';
  treeToggle.textContent = 'Show sanitized'; treeToggle.classList.remove('active');
  lastSanitized = ''; lastResult = null;
}

/* ── Copy ─────────────────────────────────────────────────────────────────── */
copyBtn.addEventListener('click', () => {
  if (!lastSanitized) return;
  navigator.clipboard.writeText(lastSanitized).then(() => {
    copyBtn.textContent = 'Copied!'; copyBtn.classList.add('copied');
    setTimeout(() => { copyBtn.textContent = 'Copy'; copyBtn.classList.remove('copied'); }, 1500);
  });
});

/* ── Init ─────────────────────────────────────────────────────────────────── */
loadExamples();
