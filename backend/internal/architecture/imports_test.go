package architecture

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestLayerImports(t *testing.T) {
	const prefix = "example.com/ai-engineering-template/backend/internal/"
	root := filepath.Join("..", "..")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		layer := ""
		switch {
		case strings.HasPrefix(rel, "internal/services/"):
			layer = "services"
		case strings.HasPrefix(rel, "internal/handlers/"):
			layer = "handlers"
		case strings.HasPrefix(rel, "cmd/"):
			layer = "cmd"
		default:
			t.Errorf("%s: unregistered production layer; update architecture policy", rel)
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			dep, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			if layer == "services" && (dep == "net/http" || strings.HasPrefix(dep, "net/http/")) {
				t.Errorf("%s: services cannot import %s", rel, dep)
			}
			if strings.HasPrefix(dep, prefix) {
				target := strings.Split(strings.TrimPrefix(dep, prefix), "/")[0]
				allowed := (layer == "handlers" && target == "services") || (layer == "cmd" && target == "handlers")
				if !allowed {
					t.Errorf("%s: forbidden dependency %s -> %s", rel, layer, target)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
