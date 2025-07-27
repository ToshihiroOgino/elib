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

function selectNote(noteId) {
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
    
    // シェア番号を設定
    initializeShareNumbers();
});

// シェア番号を初期化
function initializeShareNumbers() {
    const shareItems = document.querySelectorAll('.share-item');
    shareItems.forEach((item, index) => {
        const numberSpan = item.querySelector('.share-number');
        if (numberSpan) {
            numberSpan.textContent = index + 1;
        }
    });
}

// 共有リンクをコピー
function copyShareLink(shareId) {
    const shareUrl = window.location.origin + '/share/' + shareId;
    
    if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(shareUrl).then(() => {
            showToast('共有リンクをコピーしました');
        }).catch(err => {
            console.error('Failed to copy: ', err);
            fallbackCopyToClipboard(shareUrl);
        });
    } else {
        fallbackCopyToClipboard(shareUrl);
    }
}

// フォールバック用のコピー関数
function fallbackCopyToClipboard(text) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.select();
    try {
        document.execCommand('copy');
        showToast('共有リンクをコピーしました');
    } catch (err) {
        console.error('Failed to copy: ', err);
        showToast('コピーに失敗しました');
    }
    document.body.removeChild(textArea);
}

// 共有を削除
function deleteShare(shareId) {
    if (confirm('この共有リンクを削除しますか？')) {
        fetch('/share/' + shareId, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.message) {
                showToast('共有リンクを削除しました');
                // ページをリロードして共有リストを更新
                window.location.reload();
            } else {
                showToast('削除に失敗しました');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showToast('削除に失敗しました');
        });
    }
}

// トースト通知を表示
function showToast(message) {
    // 簡単なトースト表示（Bootstrap使用時はBootstrapのToast使用可能）
    const toast = document.createElement('div');
    toast.className = 'alert alert-info position-fixed';
    toast.style.cssText = 'top: 20px; right: 20px; z-index: 9999; max-width: 300px;';
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 3000);
}
