package plugin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover_FindsPlugins(t *testing.T) {
	tmpDir := t.TempDir()

	// Create fake unravel-system binary
	binPath := filepath.Join(tmpDir, "unravel-system")
	require.NoError(t, os.WriteFile(binPath, []byte("#!/bin/sh\n"), 0755))

	// Create non-plugin binary
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "other-tool"), []byte("#!/bin/sh\n"), 0755))

	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir
		}
		return ""
	}

	s := &service{}
	plugins, err := s.Discover()
	require.NoError(t, err)
	assert.Equal(t, []string{"unravel-system"}, plugins)
}

func TestDiscover_EmptyPATH(t *testing.T) {
	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string { return "" }

	s := &service{}
	plugins, err := s.Discover()
	require.NoError(t, err)
	assert.Nil(t, plugins)
}

func TestDiscover_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "unravel-fake"), 0755))

	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir
		}
		return ""
	}

	s := &service{}
	plugins, err := s.Discover()
	require.NoError(t, err)
	assert.Empty(t, plugins)
}

func TestDiscover_SkipsNonExecutable(t *testing.T) {
	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "unravel-noexec"), []byte("data"), 0644))

	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir
		}
		return ""
	}

	s := &service{}
	plugins, err := s.Discover()
	require.NoError(t, err)
	assert.Empty(t, plugins)
}

func TestGetInfo_Success(t *testing.T) {
	origLookPath := execLookPath
	origCommand := execCommand
	defer func() {
		execLookPath = origLookPath
		execCommand = origCommand
	}()

	execLookPath = func(file string) (string, error) {
		return "/usr/bin/" + file, nil
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", `{"name":"unravel-system","version":"1.0.0","backend":"osquery","description":"System discovery"}`)
	}

	s := &service{}
	info, err := s.GetInfo("unravel-system")
	require.NoError(t, err)
	assert.Equal(t, "unravel-system", info.Name)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, "osquery", info.Backend)
	assert.Equal(t, "System discovery", info.Description)
}

func TestGetInfo_PluginNotFound(t *testing.T) {
	origLookPath := execLookPath
	defer func() { execLookPath = origLookPath }()
	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	s := &service{}
	_, err := s.GetInfo("unravel-missing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in PATH")
}

func TestExec_Success(t *testing.T) {
	origLookPath := execLookPath
	origCommand := execCommand
	defer func() {
		execLookPath = origLookPath
		execCommand = origCommand
	}()

	execLookPath = func(file string) (string, error) {
		return "/usr/bin/" + file, nil
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "output")
	}

	var stdout bytes.Buffer
	s := &service{}
	code, err := s.Exec("unravel-system", []string{"discover"}, nil, &stdout, os.Stderr)
	require.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), "output")
}

func TestExec_PluginNotFound(t *testing.T) {
	origLookPath := execLookPath
	defer func() { execLookPath = origLookPath }()
	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	s := &service{}
	_, err := s.Exec("unravel-missing", []string{"discover"}, nil, os.Stdout, os.Stderr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in PATH")
}
