package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var once sync.Once
var rootDir string

func projectRootDir() string {
	var onceErr error
	once.Do(func() {
		currentDir, err := os.Getwd()
		if err != nil {
			onceErr = err
			return
		}
		currentDir = filepath.Clean(currentDir)
		const targetName = "go.mod"
		for {
			targetPath := filepath.Join(currentDir, targetName)

			if file, err := os.Stat(targetPath); err == nil && !file.IsDir() {
				rootDir = filepath.Dir(targetPath)
				onceErr = nil
				return
			}

			parentDir := filepath.Dir(currentDir)
			// ルートディレクトリまで到達したら（親と自分が同じになったら）ループを終了
			if parentDir == currentDir {
				break
			}
			currentDir = parentDir
		}
		onceErr = fmt.Errorf("'%s' not found in any parent directories", targetName)
	})
	if onceErr != nil {
		panic(onceErr)
	}
	return rootDir
}
