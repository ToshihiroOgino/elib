package log

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func Init() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				if source, ok := a.Value.Any().(*slog.Source); ok {
					if source.File == "" {
						a.Value = slog.StringValue("unknown")
						break
					}
					path := projectRootDir()
					path, err := filepath.Rel(path, source.File)
					if err != nil {
						panic(err)
					}
					path = filepath.ToSlash(path)
					path = fmt.Sprintf("%s:%d %s", path, source.Line, source.Function)
					a.Value = slog.StringValue(path)
				}
			}

			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
