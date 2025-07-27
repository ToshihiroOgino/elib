// 日時変換ユーティリティ関数

/**
 * UTC時刻をJST（日本標準時）に変換して表示する
 * @param {string} utcDateString - UTC時刻の文字列（YYYY-MM-DDTHH:mm:ss形式）
 * @param {string} prefix - 表示する際のプレフィックス（例：「作成日: 」）
 * @returns {string} JST形式の日時文字列
 */
function formatToJST(utcDateString, prefix = '') {
    try {
        // Go言語のtime.Timeから来るUTC時刻をパース
        const utcDate = new Date(utcDateString + 'Z'); // Zを追加してUTCとして認識させる
        
        // JST（UTC+9）に変換
        const jstOptions = {
            timeZone: 'Asia/Tokyo',
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            hour12: false
        };
        
        const jstTime = utcDate.toLocaleString('ja-JP', jstOptions).replace(/\//g, '/').replace(',', '');
        return prefix ? prefix + jstTime : jstTime;
    } catch (error) {
        console.error('日時変換エラー:', error);
        return utcDateString; // エラー時は元の文字列を返す
    }
}

/**
 * ページ読み込み時に全ての日時要素をJSTに変換
 * data-utc-time属性を持つすべての要素を検索して変換する
 */
function convertAllDatesToJST() {
    // data-utc-time属性を持つすべての要素を取得
    const dateElements = document.querySelectorAll('[data-utc-time]');
    
    dateElements.forEach(element => {
        const utcTime = element.getAttribute('data-utc-time');
        if (utcTime) {
            const currentText = element.textContent;
            let prefix = '';
            
            // 現在のテキストから「作成日:」や「更新日:」のプレフィックスを抽出
            if (currentText.includes('作成日:')) {
                prefix = '作成日: ';
            } else if (currentText.includes('更新日:')) {
                prefix = '更新日: ';
            }
            
            element.textContent = formatToJST(utcTime, prefix);
        }
    });
}

// DOM読み込み完了時に日時をJSTに変換（自動実行）
document.addEventListener('DOMContentLoaded', function () {
    convertAllDatesToJST();
});
