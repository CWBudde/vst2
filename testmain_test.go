//go:build !plugin
// +build !plugin

package vst2_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// demoPluginPath is set by TestMain after building the demoplugin as a shared library.
var demoPluginPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "vst2-test-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	switch runtime.GOOS {
	case "linux":
		demoPluginPath = buildDemoPlugin(dir, "demoplugin.so")
	case "windows":
		demoPluginPath = buildDemoPlugin(dir, "demoplugin.dll")
	case "darwin":
		demoPluginPath = buildDarwinBundle(dir)
	}
	os.Exit(m.Run())
}

// buildDemoPlugin compiles the demoplugin as c-shared and returns the output path.
// Returns empty string on failure (tests will skip).
func buildDemoPlugin(dir, name string) string {
	out := filepath.Join(dir, name)
	cmd := exec.Command("go", "build", "-buildmode", "c-shared", "-tags", "plugin", "-o", out, "./demoplugin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return ""
	}
	return out
}

// buildDarwinBundle creates a minimal .vst bundle that macOS CFBundle can load.
func buildDarwinBundle(dir string) string {
	bundleDir := filepath.Join(dir, "demoplugin.vst", "Contents", "MacOS")
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		return ""
	}

	// Build the shared library into the bundle's MacOS directory.
	libPath := filepath.Join(bundleDir, "demoplugin")
	cmd := exec.Command("go", "build", "-buildmode", "c-shared", "-tags", "plugin", "-o", libPath, "./demoplugin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return ""
	}

	// Write minimal Info.plist required by CFBundle.
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>demoplugin</string>
	<key>CFBundleName</key>
	<string>demoplugin</string>
	<key>CFBundlePackageType</key>
	<string>BNDL</string>
</dict>
</plist>
`
	plistPath := filepath.Join(dir, "demoplugin.vst", "Contents", "Info.plist")
	if err := os.WriteFile(plistPath, []byte(plist), 0o644); err != nil {
		return ""
	}

	return filepath.Join(dir, "demoplugin.vst")
}
