let currentNoteId = "";
let isModified = false;

// 保存状態を動的に更新
function updateSaveStatus(status) {
  const saveStatusElement = document.getElementById("save-status");
  const saveStatusContainer = saveStatusElement.parentElement;

  switch (status) {
    case "saved":
      saveStatusElement.textContent = "保存済み";
      break;
    case "unsaved":
      saveStatusElement.textContent = "未保存";
      break;
    case "saving":
      saveStatusElement.textContent = "保存中...";
      // 保存中のアニメーション効果
      saveStatusElement.style.animation = "pulse 1s infinite";
      break;
    case "error":
      saveStatusElement.textContent = "保存エラー";
      break;
    default:
      saveStatusElement.textContent = "保存済み";
  }

  // アニメーションをリセット（保存中以外）
  if (status !== "saving") {
    saveStatusElement.style.animation = "";
  }
}

function updateStats() {
  const textarea = document.getElementById("note-content");
  const content = textarea.value;

  // 文字数を更新（選択範囲考慮）
  updateSelectionInfo();

  // 行数を更新
  const lines = content.split("\n").length;
  document.getElementById("total-lines").innerHTML = lines;

  // カーソル位置を更新
  updateCursorPosition();

  // 保存状態を更新
  if (!isModified) {
    isModified = true;
    updateSaveStatus("unsaved");
  }
}

// カーソル位置を更新
function updateCursorPosition() {
  const textarea = document.getElementById("note-content");
  const content = textarea.value;
  const cursorPos = textarea.selectionStart;
  const beforeCursor = content.substring(0, cursorPos);
  const lineNum = beforeCursor.split("\n").length;
  const colNum = beforeCursor.split("\n").pop().length + 1;
  document.getElementById("cursor-position").innerHTML = lineNum + ":" + colNum;
}

function editTitle() {
  document.getElementById("note-title").classList.add("d-none");
  document.getElementById("title-input").classList.remove("d-none");
  document.getElementById("title-input").focus();
}

function saveTitle() {
  const input = document.getElementById("title-input");
  const span = document.getElementById("note-title");
  span.innerHTML = input.value;
  span.classList.remove("d-none");
  input.classList.add("d-none");
  saveNote();
}

function saveNote() {
  const noteId = document.getElementById("note-id").value;
  const title = document.getElementById("title-input").value;
  const content = document.getElementById("note-content").value;

  // 保存中状態を表示
  updateSaveStatus("saving");

  fetch("/note/save", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      id: noteId,
      title: title,
      content: content,
    }),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.status === "success") {
        updateSaveStatus("saved");
        isModified = false;
      } else {
        updateSaveStatus("error");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      updateSaveStatus("error");
    });
}

function createNewNote() {
  window.location.href = "/note/new";
}

function deleteNote() {
  if (confirm("このメモを削除しますか？")) {
    const noteId = document.getElementById("note-id").value;

    fetch("/note/delete/" + encodeURIComponent(noteId), {
      method: "DELETE",
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.status === "success") {
          window.location.href = "/note";
        } else {
          alert("削除に失敗しました");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("削除に失敗しました");
      });
  }
}

function selectNote(noteId) {
  // URL encode the noteId to prevent injection
  window.location.href = "/note/" + encodeURIComponent(noteId);
}

// 共有機能（閲覧のみ）
function shareReadonly() {
  shareNote(false);
}

// 共有機能（編集可）
function shareEditable() {
  shareNote(true);
}

function shareNote(editable) {
  const noteId = document.getElementById("note-id").value;
  const shareType = editable ? "編集可" : "閲覧のみ";

  if (!noteId) {
    showToast("共有するメモが選択されていません");
    return;
  }

  fetch("/share", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      noteId: noteId,
      editable: editable,
    }),
  })
    .then((response) => {
      if (response.status === 200) {
        return response.json();
      } else {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    })
    .then((data) => {
      if (data.shareId) {
        const shareUrl = window.location.origin + "/share/" + encodeURIComponent(data.shareId);

        // 共有リストを更新
        addShareToList(data.shareId, editable);

        // クリップボードにコピー
        if (navigator.clipboard && window.isSecureContext) {
          navigator.clipboard
            .writeText(shareUrl)
            .then(() => {
              showToast(`共有リンク（${shareType}）を作成し、クリップボードにコピーしました`);
            })
            .catch((err) => {
              console.error("Failed to copy: ", err);
              showToast(`共有リンク（${shareType}）を作成しました: ${shareUrl}`);
            });
        } else {
          showToast(`共有リンク（${shareType}）を作成しました: ${shareUrl}`);
        }
      } else {
        showToast("共有リンクの作成に失敗しました");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      showToast(`共有リンクの作成に失敗しました: ${error.message}`);
    });
}

function initializeEditor(noteId) {
  currentNoteId = noteId;

  updateStats();
  updateSaveStatus("saved");

  // テキストエリアの参照を取得
  const textarea = document.getElementById("note-content");

  // 様々なイベントでカーソル位置と統計情報を更新
  textarea.addEventListener("click", updateCursorPosition);
  textarea.addEventListener("keyup", updateCursorPosition);
  textarea.addEventListener("keydown", updateCursorPosition);
  textarea.addEventListener("mouseup", updateCursorPosition);
  textarea.addEventListener("focus", updateCursorPosition);
  textarea.addEventListener("select", updateCursorPosition);

  // 選択範囲変更時の更新
  textarea.addEventListener("selectionchange", updateSelectionInfo);
  textarea.addEventListener("mouseup", updateSelectionInfo);
  textarea.addEventListener("keyup", updateSelectionInfo);

  textarea.addEventListener("input", updateStats);
  textarea.addEventListener("change", updateStats);
  textarea.addEventListener("paste", () => {
    // paste後に統計情報を更新するため少し遅延
    setTimeout(updateStats, 10);
  });

  // 自動保存（5秒ごと）
  setInterval(() => {
    if (isModified) {
      saveNote();
    }
  }, 5000);

  // キーボードショートカット（Ctrl+S で保存）
  textarea.addEventListener("keydown", (event) => {
    if (event.ctrlKey && event.key === "s") {
      event.preventDefault();
      saveNote();
    }
  });
}

document.addEventListener("DOMContentLoaded", function () {
  // noteIdは外部から設定される想定
  if (typeof window.noteId !== "undefined") {
    initializeEditor(window.noteId);
  }

  // シェア番号を設定
  initializeShareNumbers();
});

// シェア番号を初期化
function initializeShareNumbers() {
  updateShareNumbers();
}

// 共有リンクをコピー
function copyShareLink(shareId) {
  const shareUrl = window.location.origin + "/share/" + shareId;

  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard
      .writeText(shareUrl)
      .then(() => {
        showToast("共有リンクをコピーしました");
      })
      .catch((err) => {
        console.error("Failed to copy: ", err);
        fallbackCopyToClipboard(shareUrl);
      });
  } else {
    fallbackCopyToClipboard(shareUrl);
  }
}

// フォールバック用のコピー関数
function fallbackCopyToClipboard(text) {
  const textArea = document.createElement("textarea");
  textArea.value = text;
  document.body.appendChild(textArea);
  textArea.select();
  try {
    window.clipboard.writeText(textArea.value);
    showToast("共有リンクをコピーしました");
  } catch (err) {
    console.error("Failed to copy: ", err);
    showToast("コピーに失敗しました");
  }
  document.body.removeChild(textArea);
}

// 共有を削除
function deleteShare(shareId) {
  if (confirm("この共有リンクを削除しますか？")) {
    fetch("/share/" + encodeURIComponent(shareId), {
      method: "DELETE",
    })
      .then((response) => {
        if (response.status === 200) {
          return response.json();
        } else {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
      })
      .then((data) => {
        if (data.message) {
          showToast("共有リンクを削除しました");
          // リストから動的に削除
          removeShareFromList(shareId);
        } else {
          showToast("削除に失敗しました");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        showToast(`削除に失敗しました: ${error.message}`);
      });
  }
}

// トースト通知を表示
function showToast(message) {
  // 簡単なトースト表示（Bootstrap使用時はBootstrapのToast使用可能）
  const toast = document.createElement("div");
  toast.className = "alert alert-info position-fixed";
  toast.style.cssText = "top: 20px; right: 20px; z-index: 9999; max-width: 300px;";
  toast.innerHTML = message;

  document.body.appendChild(toast);

  setTimeout(() => {
    if (toast.parentNode) {
      toast.parentNode.removeChild(toast);
    }
  }, 3000);
}

// 共有リストに新しいアイテムを追加
function addShareToList(shareId, editable) {
  const sharesList = document.getElementById("shares-list");

  // "共有リンクはありません" のメッセージがあれば削除
  const noSharesMessage = sharesList.querySelector(".text-muted");
  if (noSharesMessage) {
    noSharesMessage.remove();
  }

  // 新しい共有アイテムを作成
  const shareItem = document.createElement("div");
  shareItem.className = "card mb-2 share-item";

  const shareBody = document.createElement("div");
  shareBody.className = "card-body p-2";

  const shareContent = document.createElement("div");
  shareContent.className = "d-flex justify-content-between align-items-center";

  const shareInfo = document.createElement("div");
  const numberSpan = document.createElement("span");
  numberSpan.className = "badge bg-secondary me-2 share-number";

  const typeSpan = document.createElement("small");
  if (editable) {
    typeSpan.className = "text-success";
    typeSpan.innerHTML = "編集可";
  } else {
    typeSpan.className = "text-info";
    typeSpan.innerHTML = "閲覧のみ";
  }

  shareInfo.appendChild(numberSpan);
  shareInfo.appendChild(typeSpan);

  const shareActions = document.createElement("div");

  const copyButton = document.createElement("button");
  copyButton.className = "btn btn-sm btn-outline-primary me-1";
  copyButton.title = "リンクをコピー";
  copyButton.textContent = "📋";
  copyButton.onclick = () => copyShareLink(shareId);

  const deleteButton = document.createElement("button");
  deleteButton.className = "btn btn-sm btn-outline-danger";
  deleteButton.title = "削除";
  deleteButton.textContent = "🗑️";
  deleteButton.onclick = () => deleteShare(shareId);

  shareActions.appendChild(copyButton);
  shareActions.appendChild(deleteButton);

  shareContent.appendChild(shareInfo);
  shareContent.appendChild(shareActions);
  shareBody.appendChild(shareContent);
  shareItem.appendChild(shareBody);

  // リストに追加
  sharesList.appendChild(shareItem);

  // 番号を更新
  updateShareNumbers();
}

// 共有番号を更新
function updateShareNumbers() {
  const shareItems = document.querySelectorAll(".share-item");
  shareItems.forEach((item, index) => {
    const numberSpan = item.querySelector(".share-number");
    if (numberSpan) {
      numberSpan.innerHTML = index + 1;
    }
  });
}

// 共有アイテムをリストから削除
function removeShareFromList(shareId) {
  const shareItems = document.querySelectorAll(".share-item");
  shareItems.forEach((item) => {
    const deleteButton = item.querySelector(`button[onclick*="${shareId}"]`);
    if (deleteButton) {
      item.remove();
    }
  });

  // 番号を更新
  updateShareNumbers();

  // 共有がない場合はメッセージを表示
  const sharesList = document.getElementById("shares-list");
  const remainingItems = sharesList.querySelectorAll(".share-item");
  if (remainingItems.length === 0) {
    const noSharesMessage = document.createElement("div");
    noSharesMessage.className = "text-muted small";
    noSharesMessage.textContent = "共有リンクはありません";
    sharesList.appendChild(noSharesMessage);
  }
}

// フッター情報の追加ユーティリティ関数
function getDetailedStats() {
  const textarea = document.getElementById("note-content");
  const content = textarea.value;

  const stats = {
    characters: content.length,
    charactersNoSpaces: content.replace(/\s/g, "").length,
    words: content.trim() ? content.trim().split(/\s+/).length : 0,
    lines: content.split("\n").length,
    paragraphs: content.split(/\n\s*\n/).filter((p) => p.trim()).length,
  };

  return stats;
}

// 選択範囲の情報を更新
function updateSelectionInfo() {
  const textarea = document.getElementById("note-content");
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;

  if (start !== end) {
    const selectedText = textarea.value.substring(start, end);
    const selectedLength = selectedText.length;
    const selectedWords = selectedText.trim() ? selectedText.trim().split(/\s+/).length : 0;

    // 選択情報を表示（例：フッターの文字数部分に追加表示）
    const charCountElement = document.getElementById("char-count");
    charCountElement.textContent = `${textarea.value.length} (選択: ${selectedLength})`;
  } else {
    // 選択なしの場合は通常表示
    const charCountElement = document.getElementById("char-count");
    charCountElement.textContent = textarea.value.length;
  }
}
