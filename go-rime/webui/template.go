package webui

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Rime + LLM Test</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; padding: 20px; background: #f5f5f5; }
.container { max-width: 800px; margin: 0 auto; }
h1 { font-size: 20px; margin-bottom: 20px; font-weight: 600; }
.input-area { background: white; padding: 15px; border-radius: 8px; margin-bottom: 15px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
input { width: 100%; padding: 10px; font-size: 16px; border: 1px solid #ddd; border-radius: 4px; }
input:focus { outline: none; border-color: #4CAF50; }
.preedit { padding: 10px; background: #f9f9f9; border-radius: 4px; min-height: 40px; font-size: 18px; margin-bottom: 10px; }
.candidates { background: white; padding: 10px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.candidate { padding: 8px 12px; margin: 4px 0; border-radius: 4px; cursor: pointer; display: flex; align-items: center; font-size: 16px; }
.candidate:hover { background: #f0f0f0; }
.candidate.ai { background: #e8f5e9; }
.candidate .num { color: #999; margin-right: 10px; min-width: 20px; }
.candidate .text { flex: 1; }
.candidate .comment { color: #666; font-size: 14px; margin-left: 10px; }
.output { background: white; padding: 15px; border-radius: 8px; margin-bottom: 15px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); min-height: 60px; font-size: 18px; }
.controls { display: flex; gap: 10px; margin-bottom: 15px; }
button { padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; background: #4CAF50; color: white; }
button:hover { background: #45a049; }
button.secondary { background: #757575; }
button.secondary:hover { background: #616161; }
.status { font-size: 12px; color: #666; margin-top: 5px; }
.llm-status { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; }
.llm-on { background: #4CAF50; color: white; }
.llm-off { background: #999; color: white; }
</style>
</head>
<body>
<div class="container">
<h1>🎹 Rime + AI Input Method Test</h1>

<div class="controls">
<button onclick="clearInput()">Clear</button>
<button class="secondary" onclick="toggleLLM()">Toggle AI</button>
<span class="llm-status llm-on" id="llmStatus">AI ON</span>
</div>

<div class="input-area">
<input type="text" id="input" placeholder="Type pinyin here (e.g., nihao)" autocomplete="off">
<div class="status">Press Enter to get suggestions, click candidate to select</div>
</div>

<div class="preedit" id="preedit"></div>

<div class="candidates" id="candidates"></div>

<div class="output" id="output">Output will appear here...</div>

</div>

<script>
let output = [];

document.getElementById('input').addEventListener('keypress', async (e) => {
if (e.key === 'Enter') {
const text = e.target.value.trim();
if (!text) return;
await sendInput(text);
e.target.value = '';
}
});

async function sendInput(text) {
try {
const resp = await fetch('/api/input', {
method: 'POST',
headers: {'Content-Type': 'application/json'},
body: JSON.stringify({text})
});
const data = await resp.json();
if (data.success) {
displayCandidates(data);
} else {
alert('Error: ' + (data.error || 'Unknown error'));
}
} catch (err) {
alert('Request failed: ' + err.message);
}
}

function displayCandidates(data) {
document.getElementById('preedit').textContent = data.preedit || '(empty)';
const container = document.getElementById('candidates');
container.innerHTML = '';
if (data.candidates && data.candidates.length > 0) {
data.candidates.forEach((cand, i) => {
const div = document.createElement('div');
div.className = 'candidate' + (cand.isAI ? ' ai' : '');
div.onclick = () => selectCandidate(i);
div.innerHTML = ` + "`" + `
<span class="num">${i+1}.</span>
<span class="text">${cand.text}</span>
${cand.comment ? ` + "`" + `<span class="comment">${cand.comment}</span>` + "`" + ` : ''}
` + "`" + `;
container.appendChild(div);
});
} else {
container.innerHTML = '<div style="padding:10px;color:#999;">No candidates</div>';
}
}

async function selectCandidate(index) {
try {
const resp = await fetch('/api/select', {
method: 'POST',
headers: {'Content-Type': 'application/json'},
body: JSON.stringify({index})
});
const data = await resp.json();
if (data.success && data.text) {
output.push(data.text);
document.getElementById('output').textContent = output.join('');
}
// Clear display
document.getElementById('preedit').textContent = '';
document.getElementById('candidates').innerHTML = '';
} catch (err) {
alert('Select failed: ' + err.message);
}
}

async function clearInput() {
try {
await fetch('/api/clear', {method: 'POST'});
document.getElementById('preedit').textContent = '';
document.getElementById('candidates').innerHTML = '';
output = [];
document.getElementById('output').textContent = 'Output will appear here...';
} catch (err) {
alert('Clear failed: ' + err.message);
}
}

async function toggleLLM() {
try {
const resp = await fetch('/api/toggle-llm', {method: 'POST'});
const data = await resp.json();
if (data.success) {
const status = document.getElementById('llmStatus');
if (data.enabled) {
status.textContent = 'AI ON';
status.className = 'llm-status llm-on';
} else {
status.textContent = 'AI OFF';
status.className = 'llm-status llm-off';
}
}
} catch (err) {
alert('Toggle failed: ' + err.message);
}
}

// Focus input on load
document.getElementById('input').focus();
</script>
</body>
</html>
`
