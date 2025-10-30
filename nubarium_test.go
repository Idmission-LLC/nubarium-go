package nubarium

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func Test_ComprobanteDomicilio_Fixtures_Parse(t *testing.T) {
	fixturesDir := "testdata/responses"

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	if len(entries) == 0 {
		t.Fatalf("no fixtures found in %s", fixturesDir)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			b, err := os.ReadFile(filepath.Join(fixturesDir, e.Name()))
			if err != nil {
				t.Fatalf("read file: %v", err)
			}
			var resp ComprobanteDomicilioResponse
			if err := json.Unmarshal(b, &resp); err != nil {
				t.Fatalf("unmarshal: %v\njson: %s", err, string(b))
			}
			// Validate optional nested object shape if present
			_ = resp.Validaciones
		})
	}
}
