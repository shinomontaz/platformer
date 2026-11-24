package main

// Собирает game.dat из assets/ + config/ (можно добавить tileset/ и др.).
// Имена записей — lowercase, слэши прямые. Предупреждает о коллизиях:
// если два файла отличаются только регистром — в dat останется один!

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	seen := make(map[string]string) // lower(name) -> исходное имя

	out, err := os.Create("game.dat")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	zw := zip.NewWriter(out)

	for _, root := range []string{"assets", "config"} {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			name := strings.ToLower(filepath.ToSlash(path))
			if prev, ok := seen[name]; ok && prev != filepath.ToSlash(path) {
				fmt.Printf("[warn] case collision: %q vs %q -> в dat останется %q\n", prev, path, name)
			}
			seen[name] = filepath.ToSlash(path)

			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()
			dst, err := zw.Create(name)
			if err != nil {
				return err
			}
			_, err = io.Copy(dst, src)
			return err
		})

		if err != nil {
			panic(err)
		}
	}

	if err := zw.Close(); err != nil {
		panic(err)
	}
	fmt.Println("game.dat готов:", len(seen), "файлов")
}
