// 共通エディター機能
// shared-editor.js と editor.js で重複している関数をまとめたもの

/**
 * 保存状態を動的に更新する共通関数
 * @param {string} status - "saved", "unsaved", "saving", "error"
 */
function updateSaveStatus(status) {
  const saveStatusElement = document.getElementById("save-status");

  if (!saveStatusElement) return;

  switch (status) {
    case "saved":
      saveStatusElement.innerHTML = "保存済み";
      break;
    case "unsaved":
      saveStatusElement.innerHTML = "未保存";
      break;
    case "saving":
      saveStatusElement.innerHTML = "保存中...";
      // 保存中のアニメーション効果
      saveStatusElement.style.animation = "pulse 1s infinite";
      break;
    case "error":
      saveStatusElement.innerHTML = "保存エラー";
      break;
    default:
      saveStatusElement.innerHTML = "保存済み";
  }

  // アニメーションをリセット（保存中以外）
  if (status !== "saving") {
    saveStatusElement.style.animation = "";
  }
}

/**
 * カーソル位置を更新する共通関数
 */
function updateCursorPosition() {
  const textarea = document.getElementById("note-content");
  const cursorPositionElement = document.getElementById("cursor-position");

  if (!textarea || !cursorPositionElement) return;

  const content = textarea.value;
  const cursorPos = textarea.selectionStart;
  const beforeCursor = content.substring(0, cursorPos);
  const lineNum = beforeCursor.split("\n").length;
  const colNum = beforeCursor.split("\n").pop().length + 1;

  cursorPositionElement.textContent = lineNum + ":" + colNum;
}

/**
 * 選択範囲の情報を更新する共通関数
 */
function updateSelectionInfo() {
  const textarea = document.getElementById("note-content");
  const charCountElement = document.getElementById("char-count");

  if (!textarea || !charCountElement) return;

  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;

  if (start !== end) {
    const selectedText = textarea.value.substring(start, end);
    const selectedLength = selectedText.length;

    // 選択情報を表示
    charCountElement.innerHTML = `${textarea.value.length} (選択: ${selectedLength})`;
  } else {
    // 選択なしの場合は通常表示
    charCountElement.innerHTML = textarea.value.length;
  }
}

/**
 * 統計情報を更新する共通関数
 */
function updateStats() {
  const textarea = document.getElementById("note-content");

  if (!textarea) return;

  const content = textarea.value;

  // 文字数を更新（選択範囲考慮）
  updateSelectionInfo();

  // 行数を更新
  const lines = content.split("\n").length;
  const lineCountElement = document.getElementById("line-count");
  const totalLinesElement = document.getElementById("total-lines");

  if (lineCountElement) {
    lineCountElement.innerHTML = lines;
  }
  if (totalLinesElement) {
    totalLinesElement.innerHTML = lines;
  }

  // カーソル位置を更新
  updateCursorPosition();
}

/**
 * タイトル編集を開始する共通関数
 */
function editTitle() {
  const titleSpan = document.getElementById("note-title");
  const titleInput = document.getElementById("title-input");

  if (titleSpan && titleInput) {
    titleSpan.classList.add("d-none");
    titleInput.classList.remove("d-none");
    titleInput.focus();
    titleInput.select();
  }
}

/**
 * トースト通知を表示する共通関数
 * Bootstrapがある場合はBootstrapのToastを使用、ない場合は簡易版を使用
 * @param {string} message - 表示するメッセージ
 * @param {string} type - "success", "error", "info"
 */
function showToast(message, type = "info") {
  // Bootstrap Toastが利用可能な場合
  const toast = document.getElementById("toast");
  if (toast && typeof bootstrap !== "undefined") {
    const toastMessage = document.getElementById("toast-message");
    const toastHeader = toast.querySelector(".toast-header");

    if (toastMessage && toastHeader) {
      // タイプに応じて色を変更
      toastHeader.className = "toast-header";
      if (type === "success") {
        toastHeader.classList.add("bg-success", "text-white");
      } else if (type === "error") {
        toastHeader.classList.add("bg-danger", "text-white");
      } else {
        toastHeader.classList.add("bg-info", "text-white");
      }

      toastMessage.innerHTML = message;
      const bsToast = new bootstrap.Toast(toast);
      bsToast.show();
      return;
    }
  }

  // フォールバック: 簡易トースト
  const toastElement = document.createElement("div");
  let bgClass = "alert-info";
  if (type === "success") {
    bgClass = "alert-success";
  } else if (type === "error") {
    bgClass = "alert-danger";
  }

  toastElement.className = `alert ${bgClass} position-fixed`;
  toastElement.style.cssText = "top: 20px; right: 20px; z-index: 9999; max-width: 300px;";
  toastElement.textContent = message;

  document.body.appendChild(toastElement);

  setTimeout(() => {
    if (toastElement.parentNode) {
      toastElement.parentNode.removeChild(toastElement);
    }
  }, 3000);
}

/**
 * 共通のエディター初期化処理
 * @param {HTMLTextAreaElement} textarea - テキストエリア要素
 * @param {Object} options - オプション設定
 */
function initializeCommonEditorEvents(textarea, options = {}) {
  if (!textarea) return;

  // デフォルトオプション
  const defaultOptions = {
    enableAutoSave: false,
    autoSaveInterval: 5000,
    enableKeyboardShortcuts: true,
    saveCallback: null,
    markUnsavedCallback: null,
  };

  const config = { ...defaultOptions, ...options };

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

  // 入力時の統計更新
  textarea.addEventListener("input", function () {
    updateStats();
    if (config.markUnsavedCallback) {
      config.markUnsavedCallback();
    }
  });

  textarea.addEventListener("change", updateStats);
  textarea.addEventListener("paste", () => {
    // paste後に統計情報を更新するため少し遅延
    setTimeout(updateStats, 10);
  });

  // 自動保存
  if (config.enableAutoSave && config.saveCallback) {
    setInterval(() => {
      config.saveCallback();
    }, config.autoSaveInterval);
  }

  // キーボードショートカット（Ctrl+S で保存）
  if (config.enableKeyboardShortcuts && config.saveCallback) {
    textarea.addEventListener("keydown", (event) => {
      if (event.ctrlKey && event.key === "s") {
        event.preventDefault();
        config.saveCallback();
      }
    });
  }
}

/**
 * クリップボードにテキストをコピーする共通関数
 * @param {string} text - コピーするテキスト
 * @param {string} successMessage - 成功時のメッセージ
 * @param {string} errorMessage - エラー時のメッセージ
 */
function copyToClipboard(
  text,
  successMessage = "クリップボードにコピーしました",
  errorMessage = "コピーに失敗しました"
) {
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        showToast(successMessage, "success");
      })
      .catch((err) => {
        console.error("Failed to copy: ", err);
        fallbackCopyToClipboard(text, successMessage, errorMessage);
      });
  } else {
    fallbackCopyToClipboard(text, successMessage, errorMessage);
  }
}

/**
 * フォールバック用のコピー関数
 * @param {string} text - コピーするテキスト
 * @param {string} successMessage - 成功時のメッセージ
 * @param {string} errorMessage - エラー時のメッセージ
 */
function fallbackCopyToClipboard(text, successMessage, errorMessage) {
  const textArea = document.createElement("textarea");
  textArea.value = text;
  document.body.appendChild(textArea);
  textArea.select();
  try {
    document.execCommand("copy");
    showToast(successMessage, "success");
  } catch (err) {
    console.error("Failed to copy: ", err);
    showToast(errorMessage, "error");
  }
  document.body.removeChild(textArea);
}

/**
 * 日時変換の共通関数（datetime-utils.jsから）
 */
function convertUtcToLocal() {
  if (typeof convertAllDatesToJST === "function") {
    convertAllDatesToJST();
  }
}
