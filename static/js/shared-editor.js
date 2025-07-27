// 共有メモ編集用JavaScript

let isUnsaved = false;
let lastSavedContent = "";
let sharedNoteConfig = {};

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

    // テキストエリアのイベントリスナー
    noteContent.addEventListener("input", function () {
      updateStats();
      markUnsaved();
    });

    noteContent.addEventListener("keydown", function (e) {
      updateCursorPosition();
    });

    noteContent.addEventListener("click", function (e) {
      updateCursorPosition();
    });

    // 自動保存（5秒間隔）
    setInterval(function () {
      if (isUnsaved) {
        saveSharedNote();
      }
    }, 5000);

    // ページ離脱時の警告
    window.addEventListener("beforeunload", function (e) {
      if (isUnsaved) {
        e.preventDefault();
        e.returnValue =
          "保存されていない変更があります。本当にページを離れますか？";
        return e.returnValue;
      }
    });
  }

  // 日時を現地時間に変換
  convertUtcToLocal();
});

function updateStats() {
  const textarea = document.getElementById("note-content");
  const content = textarea.value;

  document.getElementById("char-count").textContent = content.length;

  const lines = content.split("\n").length;
  document.getElementById("line-count").textContent = lines;
  document.getElementById("total-lines").textContent = lines;

  updateCursorPosition();
}

function updateCursorPosition() {
  const textarea = document.getElementById("note-content");
  const content = textarea.value;
  const cursorPos = textarea.selectionStart;
  const beforeCursor = content.substring(0, cursorPos);
  const lineNum = beforeCursor.split("\n").length;
  const colNum = beforeCursor.split("\n").pop().length + 1;
  document.getElementById("cursor-position").textContent =
    lineNum + ":" + colNum;
}

function markUnsaved() {
  const content = document.getElementById("note-content").value;
  if (content !== lastSavedContent) {
    isUnsaved = true;
    document.getElementById("save-status").textContent = "未保存";
    document.getElementById("save-status").style.color = "#ffc107";
  }
}

function saveSharedNote() {
  const shareId = sharedNoteConfig.shareId;
  const content = document.getElementById("note-content").value;
  const title = document.getElementById("note-title").textContent;

  // 保存状態を表示
  const saveStatus = document.getElementById("save-status");
  saveStatus.textContent = "保存中...";
  saveStatus.style.color = "#6c757d";

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
      saveStatus.textContent = "保存済み";
      saveStatus.style.color = "#28a745";
      showToast("保存が完了しました", "success");
    })
    .catch((error) => {
      console.error("Error:", error);
      saveStatus.textContent = "保存失敗";
      saveStatus.style.color = "#dc3545";
      showToast("保存に失敗しました: " + error.message, "error");
    });
}

function editTitle() {
  const titleSpan = document.getElementById("note-title");
  const titleInput = document.getElementById("title-input");

  titleSpan.classList.add("d-none");
  titleInput.classList.remove("d-none");
  titleInput.focus();
  titleInput.select();
}

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
  navigator.clipboard
    .writeText(window.location.href)
    .then(() => {
      showToast("リンクをクリップボードにコピーしました", "success");
    })
    .catch((err) => {
      console.error("Failed to copy: ", err);
      showToast("リンクのコピーに失敗しました", "error");
    });
}

function showToast(message, type = "info") {
  const toast = document.getElementById("toast");
  const toastMessage = document.getElementById("toast-message");
  const toastHeader = toast.querySelector(".toast-header");

  // タイプに応じて色を変更
  toastHeader.className = "toast-header";
  if (type === "success") {
    toastHeader.classList.add("bg-success", "text-white");
  } else if (type === "error") {
    toastHeader.classList.add("bg-danger", "text-white");
  } else {
    toastHeader.classList.add("bg-info", "text-white");
  }

  toastMessage.textContent = message;

  const bsToast = new bootstrap.Toast(toast);
  bsToast.show();
}
