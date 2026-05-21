import './style.css';
import './app.css';

import {
    GetDecisions, CreateDecision, DeleteDecision, UpdateDecisionTitle,
    GetDecision,
    AddOption, RemoveOption,
    AddCriterion, UpdateCriterion, RemoveCriterion,
    SetScore,
    GetResults,
} from '../wailsjs/go/main/App';

// ── State ──────────────────────────────────────────────────────────────────
let currentDecisionID = null;

// ── DOM refs ───────────────────────────────────────────────────────────────
const views = {
    list:    document.getElementById('view-list'),
    edit:    document.getElementById('view-edit'),
    score:   document.getElementById('view-score'),
    results: document.getElementById('view-results'),
};

function showView(name) {
    Object.entries(views).forEach(([k, el]) => {
        el.classList.toggle('hidden', k !== name);
    });
}

// ── Helpers ────────────────────────────────────────────────────────────────
function showError(msg) {
    console.error(msg);
}

async function withErrorHandling(fn) {
    try {
        await fn();
    } catch (e) {
        showError(String(e));
    }
}

// ── Decision List View ─────────────────────────────────────────────────────
const decisionList  = document.getElementById('decision-list');
const emptyState    = document.getElementById('empty-state');
const btnNewDecision = document.getElementById('btn-new-decision');

async function renderList() {
    const decisions = await GetDecisions();
    decisionList.innerHTML = '';

    if (!decisions || decisions.length === 0) {
        decisionList.classList.add('hidden');
        emptyState.classList.remove('hidden');
        return;
    }
    decisionList.classList.remove('hidden');
    emptyState.classList.add('hidden');

    decisions.slice().reverse().forEach(d => {
        const card = document.createElement('div');
        card.className = 'decision-card';
        card.innerHTML = `
            <div class="decision-card-info">
                <div class="decision-card-title">${escHtml(d.title)}</div>
                <div class="decision-card-meta">
                    ${d.options.length} 个选项 · ${d.criteria.length} 个标准 · ${d.created_at.slice(0,10)}
                </div>
            </div>
            <div class="decision-card-actions">
                <button class="btn btn-danger btn-del" data-id="${d.id}">删除</button>
            </div>
        `;
        card.querySelector('.decision-card-info').addEventListener('click', () => openEdit(d.id));
        let delTimer = null;
        card.querySelector('.btn-del').addEventListener('click', async (e) => {
            e.stopPropagation();
            const btn = e.currentTarget;
            if (!btn.dataset.confirming) {
                btn.dataset.confirming = '1';
                btn.textContent = '确认删除？';
                btn.style.background = 'rgba(239,68,68,.2)';
                delTimer = setTimeout(() => {
                    btn.dataset.confirming = '';
                    btn.textContent = '删除';
                    btn.style.background = '';
                }, 3000);
                return;
            }
            clearTimeout(delTimer);
            await withErrorHandling(async () => {
                await DeleteDecision(d.id);
                await renderList();
            });
        });
        decisionList.appendChild(card);
    });
}

// ── Modal ──────────────────────────────────────────────────────────────────
const modalOverlay = document.getElementById('modal-overlay');
const modalInput   = document.getElementById('modal-input');
const modalCancel  = document.getElementById('modal-cancel');
const modalConfirm = document.getElementById('modal-confirm');

btnNewDecision.addEventListener('click', () => {
    modalInput.value = '';
    modalOverlay.classList.remove('hidden');
    setTimeout(() => modalInput.focus(), 50);
});

modalCancel.addEventListener('click', () => modalOverlay.classList.add('hidden'));

modalInput.addEventListener('keydown', e => {
    if (e.key === 'Enter') modalConfirm.click();
    if (e.key === 'Escape') modalCancel.click();
});

modalConfirm.addEventListener('click', async () => {
    const title = modalInput.value.trim();
    if (!title) return;
    await withErrorHandling(async () => {
        const d = await CreateDecision(title);
        modalOverlay.classList.add('hidden');
        await openEdit(d.id);
    });
});

// ── Edit View ──────────────────────────────────────────────────────────────
const editTitleEl      = document.getElementById('edit-title');
const optionsList      = document.getElementById('options-list');
const criteriaList     = document.getElementById('criteria-list');
const inputOption      = document.getElementById('input-option');
const btnAddOption     = document.getElementById('btn-add-option');
const inputCriterionName   = document.getElementById('input-criterion-name');
const inputCriterionWeight = document.getElementById('input-criterion-weight');
const btnAddCriterion  = document.getElementById('btn-add-criterion');

async function openEdit(id) {
    currentDecisionID = id;
    await renderEdit();
    showView('edit');
}

async function renderEdit() {
    const d = await GetDecision(currentDecisionID);

    // title
    editTitleEl.textContent = d.title;

    // options
    optionsList.innerHTML = '';
    (d.options || []).forEach(opt => {
        const item = document.createElement('div');
        item.className = 'tag-item';
        item.innerHTML = `
            <span>${escHtml(opt)}</span>
            <button class="btn btn-danger">×</button>
        `;
        item.querySelector('button').addEventListener('click', async () => {
            await withErrorHandling(async () => {
                await RemoveOption(currentDecisionID, opt);
                await renderEdit();
            });
        });
        optionsList.appendChild(item);
    });

    // criteria
    criteriaList.innerHTML = '';
    (d.criteria || []).forEach(c => {
        const item = document.createElement('div');
        item.className = 'criterion-item';
        item.innerHTML = `
            <div class="criterion-item-left">
                <span class="criterion-item-weight">× ${c.weight}</span>
                <span class="criterion-item-name" contenteditable="true" spellcheck="false" data-cid="${c.id}">${escHtml(c.name)}</span>
            </div>
            <div style="display:flex;align-items:center;gap:6px">
                <select class="select weight-sel" data-cid="${c.id}">
                    ${[1,2,3,4,5].map(w => `<option value="${w}" ${w===c.weight?'selected':''}>${w}</option>`).join('')}
                </select>
                <button class="btn btn-danger" data-cid="${c.id}">×</button>
            </div>
        `;
        // edit name inline
        const nameEl = item.querySelector('.criterion-item-name');
        nameEl.addEventListener('blur', async () => {
            const newName = nameEl.textContent.trim();
            if (!newName || newName === c.name) return;
            await withErrorHandling(async () => {
                await UpdateCriterion(currentDecisionID, c.id, newName, c.weight);
                await renderEdit();
            });
        });
        nameEl.addEventListener('keydown', e => { if (e.key === 'Enter') { e.preventDefault(); nameEl.blur(); } });

        // edit weight
        item.querySelector('.weight-sel').addEventListener('change', async (e) => {
            const newWeight = parseInt(e.target.value);
            await withErrorHandling(async () => {
                await UpdateCriterion(currentDecisionID, c.id, c.name, newWeight);
                await renderEdit();
            });
        });

        // remove
        item.querySelector('.btn-danger').addEventListener('click', async () => {
            await withErrorHandling(async () => {
                await RemoveCriterion(currentDecisionID, c.id);
                await renderEdit();
            });
        });
        criteriaList.appendChild(item);
    });
}

// title inline edit
editTitleEl.addEventListener('blur', async () => {
    const newTitle = editTitleEl.textContent.trim();
    if (!newTitle || !currentDecisionID) return;
    await withErrorHandling(() => UpdateDecisionTitle(currentDecisionID, newTitle));
});
editTitleEl.addEventListener('keydown', e => { if (e.key === 'Enter') { e.preventDefault(); editTitleEl.blur(); } });

// add option
btnAddOption.addEventListener('click', async () => {
    const val = inputOption.value.trim();
    if (!val) return;
    await withErrorHandling(async () => {
        await AddOption(currentDecisionID, val);
        inputOption.value = '';
        await renderEdit();
    });
});
inputOption.addEventListener('keydown', e => { if (e.key === 'Enter') btnAddOption.click(); });

// add criterion
btnAddCriterion.addEventListener('click', async () => {
    const name   = inputCriterionName.value.trim();
    const weight = parseInt(inputCriterionWeight.value);
    if (!name) return;
    await withErrorHandling(async () => {
        await AddCriterion(currentDecisionID, name, weight);
        inputCriterionName.value = '';
        await renderEdit();
    });
});
inputCriterionName.addEventListener('keydown', e => { if (e.key === 'Enter') btnAddCriterion.click(); });

// nav buttons
document.getElementById('btn-back-list').addEventListener('click', async () => {
    await renderList();
    showView('list');
});
document.getElementById('btn-go-score').addEventListener('click', () => openScore());
document.getElementById('btn-go-results').addEventListener('click', () => openResults());

// ── Score View ─────────────────────────────────────────────────────────────
const scoreTitleEl = document.getElementById('score-title');
const scoreTable   = document.getElementById('score-table');

async function openScore() {
    await renderScore();
    showView('score');
}

async function renderScore() {
    const d = await GetDecision(currentDecisionID);
    scoreTitleEl.textContent = d.title;
    scoreTable.innerHTML = '';

    if (!d.options.length || !d.criteria.length) {
        scoreTable.innerHTML = '<tr><td style="color:var(--text-muted);padding:20px">请先在编辑页添加选项和标准</td></tr>';
        return;
    }

    // build score lookup
    const scoreMap = {};
    (d.scores || []).forEach(s => {
        scoreMap[`${s.option}::${s.criterion_id}`] = s.value;
    });

    // header row
    const thead = document.createElement('thead');
    const hRow = document.createElement('tr');
    const emptyTh = document.createElement('th');
    emptyTh.textContent = '选项 / 标准';
    hRow.appendChild(emptyTh);
    d.criteria.forEach(c => {
        const th = document.createElement('th');
        th.textContent = `${c.name} (×${c.weight})`;
        hRow.appendChild(th);
    });
    thead.appendChild(hRow);
    scoreTable.appendChild(thead);

    // body
    const tbody = document.createElement('tbody');
    d.options.forEach(opt => {
        const row = document.createElement('tr');
        const nameTd = document.createElement('td');
        nameTd.textContent = opt;
        row.appendChild(nameTd);

        d.criteria.forEach(c => {
            const td = document.createElement('td');
            const key = `${opt}::${c.id}`;
            const curVal = scoreMap[key] || 0;
            const sel = document.createElement('select');
            sel.className = 'score-select';
            sel.innerHTML = `<option value="0" ${!curVal ? 'selected' : ''}>—</option>` +
                [1,2,3,4,5].map(v => `<option value="${v}" ${v===curVal?'selected':''}>${v}</option>`).join('');
            sel.addEventListener('change', async () => {
                const val = parseInt(sel.value);
                if (val === 0) return;
                await withErrorHandling(() => SetScore(currentDecisionID, opt, c.id, val));
            });
            td.appendChild(sel);
            row.appendChild(td);
        });
        tbody.appendChild(row);
    });
    scoreTable.appendChild(tbody);
}

document.getElementById('btn-back-edit').addEventListener('click', async () => {
    await renderEdit();
    showView('edit');
});
document.getElementById('btn-score-results').addEventListener('click', () => openResults());

// ── Results View ───────────────────────────────────────────────────────────
const resultsTitleEl = document.getElementById('results-title');
const resultsList    = document.getElementById('results-list');

async function openResults() {
    await renderResults();
    showView('results');
}

async function renderResults() {
    const d       = await GetDecision(currentDecisionID);
    const results = await GetResults(currentDecisionID);
    resultsTitleEl.textContent = d.title;
    resultsList.innerHTML = '';

    if (!results || results.length === 0) {
        resultsList.innerHTML = '<p style="color:var(--text-muted);padding:20px">没有数据，请先添加选项、标准并评分</p>';
        return;
    }

    const maxScore = Math.max(...results.map(r => r.score), 1);

    results.forEach(r => {
        const card = document.createElement('div');
        card.className = 'result-card';
        const pct = maxScore > 0 ? (r.score / maxScore * 100).toFixed(1) : 0;
        const rankClass = r.rank <= 3 ? `rank-${r.rank}` : '';
        card.innerHTML = `
            <div class="result-rank ${rankClass}">${r.rank}</div>
            <div class="result-info">
                <div class="result-option">${escHtml(r.option)}</div>
                <div class="result-bar-bg">
                    <div class="result-bar-fill" style="width:${pct}%"></div>
                </div>
            </div>
            <div class="result-score">${r.score.toFixed(2)}</div>
        `;
        resultsList.appendChild(card);
    });
}

document.getElementById('btn-back-score').addEventListener('click', async () => {
    await renderScore();
    showView('score');
});

// ── Utility ────────────────────────────────────────────────────────────────
function escHtml(str) {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}

// ── Init ───────────────────────────────────────────────────────────────────
renderList();
