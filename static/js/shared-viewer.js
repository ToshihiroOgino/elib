// 共有メモ閲覧用JavaScript

let sharedNoteConfig = {};

document.addEventListener("DOMContentLoaded", function () {
  // HTMLのデータ属性から設定を取得
  const dataElement = document.getElementById("shared-note-data");
  sharedNoteConfig = {
    noteId: dataElement.dataset.noteId,
    editable: dataElement.dataset.editable === "true",
    isSharedView: dataElement.dataset.isSharedView === "true",
  };

  // 日時を現地時間に変換
  convertUtcToLocal();

  // 読み取り専用での追加機能があれば実装
  setupReadOnlyFeatures();
});

function setupReadOnlyFeatures() {
  const contentDiv = document.querySelector(".shared-note-content");
  if (contentDiv) {
    contentDiv.addEventListener("mouseup", function () {
      const selection = window.getSelection();
      if (selection.toString().length > 0) {
        console.log("Selected text:", selection.toString());
      }
    });
  }
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
