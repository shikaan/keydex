package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shikaan/keydex/pkg/info"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

const binaryPath = "./" + info.NAME
const fixtureDB = "test.kdbx"
const fixturePassword = "test-password"

func createFixtureDB() error {
	file, err := os.Create(fixtureDB)
	if err != nil {
		return err
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(fixturePassword)

	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "TestDB"
	rootGroup.Entries = make([]gokeepasslib.Entry, 0)

	codingGroup := gokeepasslib.NewGroup()
	codingGroup.Name = "Coding"
	codingGroup.Entries = make([]gokeepasslib.Entry, 0)
	codingGroup.Groups = make([]gokeepasslib.Group, 0)

	github := gokeepasslib.NewEntry()
	github.Values = append(github.Values,
		gokeepasslib.ValueData{
			Key:   "Title",
			Value: gokeepasslib.V{Content: "GitHub"}},
		gokeepasslib.ValueData{
			Key:   "UserName",
			Value: gokeepasslib.V{Content: "ghuser"}},
		gokeepasslib.ValueData{
			Key: "Password",
			Value: gokeepasslib.V{
				Content:   "ghpass123",
				Protected: wrappers.NewBoolWrapper(true)}},
	)

	gitlab := gokeepasslib.NewEntry()
	gitlab.Values = append(gitlab.Values,
		gokeepasslib.ValueData{
			Key:   "Title",
			Value: gokeepasslib.V{Content: "GitLab"}},
		gokeepasslib.ValueData{
			Key:   "UserName",
			Value: gokeepasslib.V{Content: "gluser"}},
		gokeepasslib.ValueData{
			Key: "Password",
			Value: gokeepasslib.V{
				Content: "glpass456", Protected: wrappers.NewBoolWrapper(true)}},
	)

	codingGroup.Entries = append(codingGroup.Entries, github, gitlab)
	rootGroup.Groups = append(rootGroup.Groups, codingGroup)
	db.Content.Root.Groups = []gokeepasslib.Group{rootGroup}

	if err := db.LockProtectedEntries(); err != nil {
		return err
	}
	if err := gokeepasslib.NewEncoder(file).Encode(db); err != nil {
		return err
	}
	file.Close()
	return nil
}

func TestMain(m *testing.M) {
	originalWD, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory")
	}
	defer os.Chdir(originalWD)

	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get test file path")
	}

	if err = os.Chdir(filepath.Join(filepath.Dir(testFilePath), "..")); err != nil {
		panic("Cannot set working directory")
	}

	_, err = os.Stat(binaryPath)

	if err != nil {
		panic("keydex binary must be present; run 'make build' before this test")
	}

	if err = createFixtureDB(); err != nil {
		panic(err.Error())
	}

	code := m.Run()

	os.RemoveAll(fixtureDB)
	os.Exit(code)
}

func runKeydex(t *testing.T, env map[string]string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Start from a clean env to avoid inheriting KEYDEX_* vars from the
	// developer's shell, then layer in the caller's overrides.
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
	}
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		println(err.Error())
		exitCode = -1
	}

	return outBuf.String(), errBuf.String(), exitCode
}

func TestCommandList(t *testing.T) {
	aliases := []string{"list", "ls"}

	for _, alias := range aliases {
		t.Run(alias+" list all entries (password from env)", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, fixtureDB)

			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stdout, "GitHub") {
				t.Errorf("expected stdout to contain GitHub, got:\n%s", stdout)
			}
			if !strings.Contains(stdout, "GitLab") {
				t.Errorf("expected stdout to contain GitLab, got:\n%s", stdout)
			}
		})

		t.Run(alias+" list all entries (password and archive from env)", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
				"KEYDEX_DATABASE":   fixtureDB,
			}, alias)

			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stdout, "GitHub") {
				t.Errorf("expected stdout to contain GitHub, got:\n%s", stdout)
			}
			if !strings.Contains(stdout, "GitLab") {
				t.Errorf("expected stdout to contain GitLab, got:\n%s", stdout)
			}
		})

		t.Run(alias+" fails with non-existing archive", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, "non-existent.kdbx")

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "no such file or directory") {
				t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with wrong password", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": "wrong-password",
			}, alias, fixtureDB)

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "Wrong password?") {
				t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with non-existing key file", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, "open", "--key", "nonexistent.key", fixtureDB)

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "no such file or directory") {
				t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with invalid key file", func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-invalid-key-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.WriteString("not-a-valid-keyfile-content")
			tmpFile.Close()

			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, "--key", tmpFile.Name(), fixtureDB)

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "Wrong password?") {
				t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with invalid database", func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-invalid-key-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.WriteString("not-a-valid-database-content")
			tmpFile.Close()

			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, tmpFile.Name())

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "unexpected EOF") {
				t.Errorf("expected 'unexpected EOF' in stderr, got:\n%s", stderr)
			}
		})
	}
}

func TestCommandCopy(t *testing.T) {
	aliases := []string{"cp", "password", "pwd", "copy-password"}

	for _, alias := range aliases {
		// Cannot test clipboard here

		t.Run(alias+" fails with missing reference", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, fixtureDB)

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "Missing reference") {
				t.Errorf("expected 'Missing reference' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with non-existing archive", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, "non-existent.kdbx", "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "no such file or directory") {
				t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with wrong password", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": "wrong-password",
			}, alias, fixtureDB, "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "Wrong password?") {
				t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with non-existing key file", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, "open", "--key", "nonexistent.key", fixtureDB, "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "no such file or directory") {
				t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with invalid key file", func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-invalid-key-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.WriteString("not-a-valid-keyfile-content")
			tmpFile.Close()

			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, "--key", tmpFile.Name(), fixtureDB, "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "Wrong password?") {
				t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with invalid database", func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-invalid-key-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.WriteString("not-a-valid-database-content")
			tmpFile.Close()

			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, tmpFile.Name(), "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "unexpected EOF") {
				t.Errorf("expected 'unexpected EOF' in stderr, got:\n%s", stderr)
			}
		})
	}
}

func TestCommandOpen(t *testing.T) {
	aliases := []string{"open", "edit"}

	for _, alias := range aliases {
		// TUI is tested in tui_test.go

		t.Run(alias+" fails with non-existing archive", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, "non-existent.kdbx", "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "no such file or directory") {
				t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with wrong password", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": "wrong-password",
			}, alias, fixtureDB, "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatalf("expected exit != 0, got %d. stdout: %s", exitCode, stdout)
			}

			if !strings.Contains(stderr, "Wrong password?") {
				t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with non-existing key file", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, "open", "--key", "nonexistent.key", fixtureDB, "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "no such file or directory") {
				t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with invalid key file", func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-invalid-key-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.WriteString("not-a-valid-keyfile-content")
			tmpFile.Close()

			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, "--key", tmpFile.Name(), fixtureDB, "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "Wrong password?") {
				t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
			}
		})

		t.Run(alias+" fails with invalid database", func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-invalid-key-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.WriteString("not-a-valid-database-content")
			tmpFile.Close()

			_, stderr, exitCode := runKeydex(t, map[string]string{
				"KEYDEX_PASSPHRASE": fixturePassword,
			}, alias, tmpFile.Name(), "/TestDB/Coding/GitHub")

			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, "unexpected EOF") {
				t.Errorf("expected 'unexpected EOF' in stderr, got:\n%s", stderr)
			}
		})
	}
}

func TestCommandHelp(t *testing.T) {
	aliases := []string{"help", "-h", "--help"}

	for _, alias := range aliases {
		t.Run(alias+" shows general help", func(t *testing.T) {
			stdout, stderr, exitCode := runKeydex(t, nil, alias)

			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stdout, "Usage:") && !strings.Contains(stdout, "USAGE:") {
				t.Errorf("expected stdout to contain Usage, got:\n%s", stdout)
			}
			if !strings.Contains(stdout, "Commands") && !strings.Contains(stdout, "COMMANDS") {
				t.Errorf("expected stdout to contain Commands, got:\n%s", stdout)
			}
		})
	}

	t.Run("shows help with --help flag", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, nil, "--help")

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "Usage:") && !strings.Contains(stdout, "USAGE:") {
			t.Errorf("expected stdout to contain Usage, got:\n%s", stdout)
		}
	})

	t.Run("shows help for subcommand list", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, nil, "list", "--help")

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "list") && !strings.Contains(stdout, "ls") {
			t.Errorf("expected stdout to mention list/ls, got:\n%s", stdout)
		}
	})

	t.Run("shows help for subcommand open", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, nil, "open", "--help")

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "open") && !strings.Contains(stdout, "edit") {
			t.Errorf("expected stdout to mention open/edit, got:\n%s", stdout)
		}
	})

	t.Run("shows help for subcommand copy", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, nil, "copy", "--help")

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "copy") && !strings.Contains(stdout, "password") {
			t.Errorf("expected stdout to mention copy/password, got:\n%s", stdout)
		}
	})
}
