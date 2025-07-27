let currentNoteId = "";
let isModified = false;

// updateSaveStatus, updateStats, updateCursorPosition関数はcommon-editor.jsに移動

// editTitle関数はcommon-editor.jsに移動

// editor.js専用のupdateStats関数（isModifiedロジック含む）
function updateStatsWithModified() {
  updateStats();
  
  // 保存状態を更新
  if (!isModified) {
    isModified = true;
    updateSaveStatus("unsaved");
  }
}

function saveTitle() {
  const input = document.getElementById("title-input");
  const span = document.getElementById("note-title");
  span.textContent = input.value;
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
    showToast("共有するメモが選択されていません", "error");
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
              showToast(`共有リンク（${shareType}）を作成し、クリップボードにコピーしました`, "success");
            })
            .catch((err) => {
              console.error("Failed to copy: ", err);
              showToast(`共有リンク（${shareType}）を作成しました: ${shareUrl}`, "info");
            });
        } else {
          showToast(`共有リンク（${shareType}）を作成しました: ${shareUrl}`, "info");
        }
      } else {
        showToast("共有リンクの作成に失敗しました", "error");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      showToast(`共有リンクの作成に失敗しました: ${error.message}`, "error");
    });
}

function initializeEditor(noteId) {
  currentNoteId = noteId;

  updateStats();
  updateSaveStatus("saved");

  // テキストエリアの参照を取得
  const textarea = document.getElementById("note-content");

  // 共通のエディターイベントを初期化
  initializeCommonEditorEvents(textarea, {
    enableAutoSave: true,
    autoSaveInterval: 5000,
    enableKeyboardShortcuts: true,
    saveCallback: () => {
      if (isModified) {
        saveNote();
      }
    },
    markUnsavedCallback: updateStatsWithModified
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
  copyToClipboard(shareUrl, "共有リンクをコピーしました", "コピーに失敗しました");
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
          showToast("共有リンクを削除しました", "success");
          // リストから動的に削除
          removeShareFromList(shareId);
        } else {
          showToast("削除に失敗しました", "error");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        showToast(`削除に失敗しました: ${error.message}`, "error");
      });
  }
}

// showToast関数はcommon-editor.jsに移動

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

// updateSelectionInfo関数はcommon-editor.jsに移動
