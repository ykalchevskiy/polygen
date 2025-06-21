package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create .polygen.json config file
		configFile := filepath.Join(tempDir, ".polygen.json")
		createFile(t, configFile, `{
	"$schema": "https://raw.githubusercontent.com/ykalchevskiy/polygen/main/schema.json",
	"types": [
		{
			"type": "ItemValue",
			"interface": "IsItemValue",
			"package": "pkg",
			"subtypes": {
				"ItemValue1": {},
				"ItemValue2": {}
			}
		}
	]
}`)

		// Run the generator
		cmd := exec.Command("go", "run", ".", "-config", configFile)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("generator failed: %v\nOutput: %s", err, output)
		}

		// Verify generated file exists
		genFile := filepath.Join(tempDir, "item_value_polygen.go")
		if _, err := os.Stat(genFile); os.IsNotExist(err) {
			t.Log(tempDir)
			t.Fatalf("generated file does not exist: %v", err)
		}

		// Verify generated file content
		code, err := os.ReadFile(genFile)
		if err != nil {
			t.Fatalf("failed to read generated file: %v", err)
		}

		required := []string{
			"package pkg",
			"type ItemValue struct {",
			"IsItemValue",
			"func (v ItemValue) MarshalJSON() ([]byte, error)",
			"func (v *ItemValue) UnmarshalJSON(data []byte) error",
			`"type", typeName`,
			`case "item-value-1":`,
			`case "item-value-2":`,
			`reflect.TypeOf((*ItemValue1)(nil)).Elem():`,
			`reflect.TypeOf((*ItemValue2)(nil)).Elem():`,
		}

		for _, r := range required {
			if !bytes.Contains(code, []byte(r)) {
				t.Errorf("generated code missing required part: %q", r)
				t.Logf("Generated code:\n%s", string(code))
			}
		}
	})

	t.Run("execute", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create .polygen.json config file
		configFile := filepath.Join(tempDir, ".polygen.json")
		createFile(t, configFile, `{
	"$schema": "https://raw.githubusercontent.com/ykalchevskiy/polygen/main/schema.json",
	"strictByDefault": true,
	"defaultDescriptor": "kind",
	"types": [
		{
			"type": "ItemValue",
			"interface": "IsItemValue",
			"package": "main",
			"subtypes": {
				"ItemValue1": {
					"name": "item-value-1"
				},
				"ItemValue2": {
					"name": "item-value-2"
				}
			}
		}
	]
}`)

		// Create types.go
		createFile(t, filepath.Join(tempDir, "item_value.go"), `package main

type IsItemValue interface {
	isItemValue()
}

type ItemValue1 struct {
	Value string
}

func (ItemValue1) isItemValue() {}

type ItemValue2 struct {
	Amount int
}

func (ItemValue2) isItemValue() {}
`)

		// Create main.go
		createFile(t, filepath.Join(tempDir, "main.go"), `package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	// Test value type
	text := ItemValue{IsItemValue: ItemValue1{Value: "hello"}}
	data, err := json.Marshal(text)
	if err != nil {
		fmt.Printf("error marshaling text: %v\n", err)
		return
	}
	fmt.Printf("text: %s\n", data)

	// Test struct type with pointer
	number := ItemValue{IsItemValue: &ItemValue2{Amount: 42}}
	data, err = json.Marshal(number)
	if err != nil {
		fmt.Printf("error marshaling number: %v\n", err)
		return
	}
	fmt.Printf("number: %s\n", data)

	// Test strict unmarshaling with unknown field
	extraJSON := `+"`"+`{"kind":"item-value-1","Value":"hello","extra":"field"}`+"`"+`
	var strict ItemValue
	err = json.Unmarshal([]byte(extraJSON), &strict)
	if err == nil {
		fmt.Printf("error: strict unmarshal should fail with unknown field\n")
		return
	}
	if !strings.Contains(err.Error(), "unknown field") {
		fmt.Printf("error: wrong error type for strict unmarshal: %v\n", err)
		return
	}
	fmt.Printf("strict: correct error on unknown field\n")

	// Test normal unmarshaling
	var decoded ItemValue
	err = json.Unmarshal([]byte(`+"`"+`{"kind":"item-value-1","Value":"decoded"}`+"`"+`), &decoded)
	if err != nil {
		fmt.Printf("error unmarshaling: %v\n", err)
		return
	}
	fmt.Printf("unmarshalled: %#v\n", decoded)

	// Test null handling
	var empty ItemValue
	data, _ = json.Marshal(empty)
	fmt.Printf("null: %s\n", data)
}`)

		// Initialize go module
		cmd := exec.Command("go", "mod", "init", "test")
		cmd.Dir = tempDir
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to initialize module: %v\nOutput: %s", err, output)
		}

		// Run the generator
		cmd = exec.Command("go", "run", ".", "-config", configFile)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("generator failed: %v\nOutput: %s", err, output)
		}

		// Verify generated file exists
		genFile := filepath.Join(tempDir, "item_value_polygen.go")
		if _, err := os.Stat(genFile); os.IsNotExist(err) {
			t.Fatalf("generated file does not exist: %v", err)
		}

		// Run the test program
		cmd = exec.Command("go", "run", ".")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("test failed: %v\nOutput: %s", err, output)
		}

		// Define expected lines in order
		expectedLines := []string{
			`text: {"kind":"item-value-1","Value":"hello"}`,
			`number: {"kind":"item-value-2","Amount":42}`,
			`strict: correct error on unknown field`,
			`unmarshalled: main.ItemValue{IsItemValue:main.ItemValue1{Value:"decoded"}}`,
			`null: null`,
			"", // Empty line at the end
		}

		lines := strings.Split(string(output), "\n")

		// Check we have exactly the expected number of lines
		if len(lines) != len(expectedLines) {
			t.Errorf("Got %d lines, want %d lines\nGot output:\n%s\nWant output:\n%s",
				len(lines), len(expectedLines),
				string(output), strings.Join(expectedLines, "\n"))
			return
		}

		// Check each line matches exactly
		for i, want := range expectedLines {
			got := lines[i]
			if got != want {
				t.Errorf("Line %d mismatch:\ngot:  %q\nwant: %q", i+1, got, want)
			}
		}
	})
}

func createFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}
