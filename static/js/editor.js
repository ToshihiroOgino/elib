let currentNoteId = '';
let isModified = false;

// 統計情報を更新
function updateStats() {
    const textarea = document.getElementById('note-content');
    const content = textarea.value;

    document.getElementById('char-count').textContent = content.length;

    const lines = content.split('\n').length;
    document.getElementById('line-count').textContent = lines;
    document.getElementById('total-lines').textContent = lines;

    const cursorPos = textarea.selectionStart;
    const beforeCursor = content.substring(0, cursorPos);
    const lineNum = beforeCursor.split('\n').length;
    const colNum = beforeCursor.split('\n').pop().length + 1;
    document.getElementById('cursor-position').textContent = lineNum + ':' + colNum;

    isModified = true;
    document.getElementById('save-status').textContent = '未保存';
}

function editTitle() {
    document.getElementById('note-title').classList.add('d-none');
    document.getElementById('title-input').classList.remove('d-none');
    document.getElementById('title-input').focus();
}

function saveTitle() {
    const input = document.getElementById('title-input');
    const span = document.getElementById('note-title');
    span.textContent = input.value;
    span.classList.remove('d-none');
    input.classList.add('d-none');
    saveNote();
}

function saveNote() {
    const noteId = document.getElementById('note-id').value;
    const title = document.getElementById('title-input').value;
    const content = document.getElementById('note-content').value;

    fetch('/note/save', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            id: noteId,
            title: title,
            content: content
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                document.getElementById('save-status').textContent = '保存済み';
                isModified = false;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            document.getElementById('save-status').textContent = '保存エラー';
        });
}

function createNewNote() {
    window.location.href = '/note/new';
}

function deleteNote() {
    if (confirm('このメモを削除しますか？')) {
        const noteId = document.getElementById('note-id').value;
        fetch('/note/delete/' + noteId, {
            method: 'DELETE'
        })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    window.location.href = '/note/new';
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('削除に失敗しました');
            });
    }
}

function shareNote() {
    const noteId = document.getElementById('note-id').value;
    const shareUrl = window.location.origin + '/note/view/' + noteId;
    navigator.clipboard.writeText(shareUrl).then(() => {
        alert('共有URLをクリップボードにコピーしました: ' + shareUrl);
    }).catch(err => {
        alert('共有URL: ' + shareUrl);
    });
}

function selectNote(noteId) {
    // if (isModified && !confirm('未保存の変更があります。移動しますか？')) {
    //     return;
    // }
    window.location.href = '/note/' + noteId;
}

function initializeEditor(noteId) {
    currentNoteId = noteId;
    updateStats();

    // カーソル位置の追跡
    const textarea = document.getElementById('note-content');
    textarea.addEventListener('click', updateStats);
    textarea.addEventListener('keyup', updateStats);

    // 自動保存（5秒ごと）
    setInterval(() => {
        if (isModified) {
            saveNote();
        }
    }, 5000);
}

document.addEventListener('DOMContentLoaded', function () {
    // noteIdは外部から設定される想定
    if (typeof window.noteId !== 'undefined') {
        initializeEditor(window.noteId);
    }
});
