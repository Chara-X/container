package container

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chara-X/util/archive/tar"
)

func Install(from string) {
	var to = strings.TrimSuffix(from, ".tar")
	tar.Extract(from, to)
	var r, _ = os.Open(to + "/manifest.json")
	var w, _ = os.Create(to + "/lowerdir")
	defer r.Close()
	defer w.Close()
	var manifest []map[string]any
	json.NewDecoder(r).Decode(&manifest)
	var layers []string
	for _, layer := range manifest[0]["Layers"].([]any) {
		var from, to = filepath.Join(to, layer.(string)), filepath.Join(to, "layers", filepath.Base(layer.(string)))
		layers = append(layers, to)
		tar.Extract(from, to)
	}
	w.WriteString(strings.Join(layers, ":"))
}
