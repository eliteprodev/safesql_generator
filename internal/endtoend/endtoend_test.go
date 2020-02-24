package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kyleconroy/sqlc/internal/cmd"
)

func TestExamples(t *testing.T) {
	t.Parallel()
	examples, err := filepath.Abs(filepath.Join("..", "..", "examples"))
	if err != nil {
		t.Fatal(err)
	}

	files, err := ioutil.ReadDir(examples)
	if err != nil {
		t.Fatal(err)
	}

	for _, replay := range files {
		if !replay.IsDir() {
			continue
		}
		tc := replay.Name()
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(examples, tc)
			var stderr bytes.Buffer
			output, err := cmd.Generate(path, &stderr)
			if err != nil {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, output)
		})
	}
}

func TestReplay(t *testing.T) {
	t.Parallel()

	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, replay := range files {
		if !replay.IsDir() {
			continue
		}
		tc := replay.Name()
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path, _ := filepath.Abs(filepath.Join("testdata", tc))
			var stderr bytes.Buffer
			expected := expectedStderr(t, path)
			output, err := cmd.Generate(path, &stderr)
			if len(expected) == 0 && err != nil {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, output)
			if diff := cmp.Diff(expected, stderr.String()); diff != "" {
				t.Errorf("stderr differed (-want +got):\n%s", diff)
			}
		})
	}
}

func cmpDirectory(t *testing.T, dir string, actual map[string]string) {
	expected := map[string]string{}
	var ff = func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if file.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".kt") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") || strings.Contains(path, "src/test/") {
			return nil
		}
		blob, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		expected[path] = string(blob)
		return nil
	}
	if err := filepath.Walk(dir, ff); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(expected, actual, cmpopts.EquateEmpty()) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			name := name
			tn := strings.Replace(name, dir+"/", "", -1)
			t.Run(tn, func(t *testing.T) {
				if actual[name] == "" {
					t.Errorf("%s is empty", name)
					return
				}
				if diff := cmp.Diff(contents, actual[name]); diff != "" {
					t.Errorf("%s differed (-want +got):\n%s", name, diff)
				}
			})
		}
	}
}

func expectedStderr(t *testing.T, dir string) string {
	t.Helper()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	stderr := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		rd, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			t.Fatalf("could not open %s: %v", file.Name(), err)
		}
		scanner := bufio.NewScanner(rd)
		capture := false
		for scanner.Scan() {
			text := scanner.Text()
			if text == "-- stderr" {
				capture = true
				continue
			}
			if capture == true && strings.HasPrefix(text, "--") {
				stderr += strings.TrimPrefix(text, "-- ") + "\n"
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
	}
	return stderr
}
