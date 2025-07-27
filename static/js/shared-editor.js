// 共有メモ編集用JavaScript

let isUnsaved = false;
let lastSavedContent = "";
let sharedNoteConfig = {};

// updateSaveStatus関数はcommon-editor.jsに移動

document.addEventListener("DOMContentLoaded", function () {
  const dataElement = document.getElementById("shared-note-data");
  sharedNoteConfig = {
    shareId: dataElement.dataset.shareId,
    editable: dataElement.dataset.editable === "true",
    isSharedView: dataElement.dataset.isSharedView === "true",
  };

  const noteContent = document.getElementById("note-content");

  if (noteContent) {
    lastSavedContent = noteContent.value;
    updateStats();
    // 初期状態は保存済み
    updateSaveStatus("saved");

    // 共通のエディターイベントを初期化
    initializeCommonEditorEvents(noteContent, {
      enableAutoSave: true,
      autoSaveInterval: 5000,
      enableKeyboardShortcuts: true,
      saveCallback: () => {
        if (isUnsaved) {
          saveSharedNote();
        }
      },
      markUnsavedCallback: markUnsaved
    });

    // ページ離脱時の警告
    window.addEventListener("beforeunload", function (e) {
      if (isUnsaved) {
        e.preventDefault();
        e.returnValue = "保存されていない変更があります。本当にページを離れますか？";
        return e.returnValue;
      }
    });
  }

  // 日時を現地時間に変換
  convertUtcToLocal();
});

// updateStats, updateCursorPosition, updateSelectionInfo関数はcommon-editor.jsに移動

function markUnsaved() {
  const content = document.getElementById("note-content").value;
  if (content !== lastSavedContent) {
    isUnsaved = true;
    updateSaveStatus("unsaved");
  }
}

function saveSharedNote() {
  const shareId = sharedNoteConfig.shareId;
  const content = document.getElementById("note-content").value;
  const title = document.getElementById("note-title").textContent;

  // 保存中状態を表示
  updateSaveStatus("saving");

  fetch(`/share/${shareId}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      title: title,
      content: content,
    }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("保存に失敗しました");
      }
      return response.json();
    })
    .then((data) => {
      lastSavedContent = content;
      isUnsaved = false;
      updateSaveStatus("saved");
      showToast("保存が完了しました", "success");
    })
    .catch((error) => {
      console.error("Error:", error);
      updateSaveStatus("error");
      showToast("保存に失敗しました: " + error.message, "error");
    });
}

// editTitle関数はcommon-editor.jsに移動

function saveTitle() {
  const titleSpan = document.getElementById("note-title");
  const titleInput = document.getElementById("title-input");
  const newTitle = titleInput.value.trim() || "Untitled";

  titleSpan.textContent = newTitle;
  titleSpan.classList.remove("d-none");
  titleInput.classList.add("d-none");

  // タイトルの変更を保存
  markUnsaved();
  saveSharedNote();
}

function copyCurrentUrl() {
  copyToClipboard(window.location.href, "リンクをクリップボードにコピーしました", "リンクのコピーに失敗しました");
}

// showToast関数はcommon-editor.jsに移動
