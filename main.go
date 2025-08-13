package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Copy a single file from src to dst
func copyFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0777)
	if err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}

// Copy a directory recursively
func copyDir(src, dst string) error {
	err := os.RemoveAll(dst)
	if err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0777)
		}
		return copyFile(path, target)
	})
}

// Run a command with arguments
func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func main() {
	srcDirs := []string{"EnterpriseS", "csvlk-pack"}
	dst := `C:\Windows\System32\spp\tokens\skus`

	for _, src := range srcDirs {
		target := filepath.Join(dst, filepath.Base(src))
		if err := copyDir(src, target); err != nil {
			os.Exit(1)
		}
	}

	file, err := os.Open("cmd.txt")
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(os.ExpandEnv(line))
		if len(parts) > 0 {
			run(parts[0], parts[1:]...)
		}
	}
}