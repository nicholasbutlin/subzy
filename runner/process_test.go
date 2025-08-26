package runner

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// setupTestFingerprints copies testdata/fingerprint.json to the expected location and returns a cleanup function
func setupTestFingerprints(t *testing.T) func() {
	fpPath, err := GetFingerprintPath()
	if err != nil {
		t.Fatalf("GetFingerprintPath failed: %v", err)
	}
	src := filepath.Join("testdata", "fingerprint.json")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read test fingerprint.json: %v", err)
	}
	err = os.WriteFile(fpPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write fingerprint.json to expected path: %v", err)
	}
	return func() { os.Remove(fpPath) }
}

func TestFingerprints_LoadsTestData(t *testing.T) {
	cleanup := setupTestFingerprints(t)
	defer cleanup()

	fps, err := Fingerprints()
	if err != nil {
		t.Fatalf("Fingerprints() failed: %v", err)
	}
	if len(fps) != 1 {
		t.Errorf("Expected 1 fingerprint, got %d", len(fps))
	}
	if fps[0].Service != "TestService" {
		t.Errorf("Unexpected fingerprint service: %v", fps[0].Service)
	}
}

func TestProcess_NoTargets(t *testing.T) {
	cleanup := setupTestFingerprints(t)
	defer cleanup()

	config := &Config{
		Targets:     "test_subdomains.txt",
		Concurrency: 1,
		Output:      "",
	}

	// Create a test subdomains file
	tmpfile, err := os.CreateTemp("", "test_subdomains.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.Write([]byte("sub1.example.com\nsub2.example.com\n"))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()
	config.Targets = tmpfile.Name()

	// Should not panic or error
	err = Process(config)
	if err != nil {
		t.Errorf("Process failed: %v", err)
	}
}

func TestReadSubdomains(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "subdomains.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	content := "sub1.example.com\nsub2.example.com\n"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	subdomains, err := readSubdomains(tmpfile.Name())
	if err != nil {
		t.Fatalf("readSubdomains failed: %v", err)
	}
	if !reflect.DeepEqual(subdomains, []string{"sub1.example.com", "sub2.example.com"}) {
		t.Errorf("Unexpected subdomains: %v", subdomains)
	}
}

func TestIsValidUrl(t *testing.T) {
	valid := "http://example.com"
	invalid := "not_a_url"
	if !isValidUrl(valid) {
		t.Errorf("Expected valid URL: %s", valid)
	}
	if isValidUrl(invalid) {
		t.Errorf("Expected invalid URL: %s", invalid)
	}
}
