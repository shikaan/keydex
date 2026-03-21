package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shikaan/keydex/cmd"
	"github.com/shikaan/keydex/pkg/cli"
	"github.com/shikaan/keydex/pkg/info"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

const binaryPath = "./" + info.NAME
const fixtureDB = "test.kdbx"
const fixturePassword = "test-password"
const fixtureDB2 = "test2.kdbx"
const fixturePassword2 = "test-password-2"

func createFixtureDB() error {
	file, err := os.Create(fixtureDB)
	if err != nil {
		return err
	}

	db, err := kdbx.NewFromFile(file)
	if err != nil {
		return err
	}
	if err := db.SetPasswordAndKey(fixturePassword, ""); err != nil {
		return err
	}

	rootGroup := db.NewGroup("TestDB")
	codingGroup := db.NewGroup("Coding")

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
	rootGroup.Groups = append(rootGroup.Groups, *codingGroup)
	db.Content.Root.Groups = []gokeepasslib.Group{*rootGroup}

	return db.Save()
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

	if err = createFixtureDB2(); err != nil {
		panic(err.Error())
	}

	code := m.Run()

	os.RemoveAll(fixtureDB)
	os.RemoveAll(fixtureDB2)
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
		t.Errorf("command error: %s", err.Error())
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

	t.Run("shows help for subcommand create", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, nil, "create", "--help")

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "create") {
			t.Errorf("expected stdout to mention create, got:\n%s", stdout)
		}
	})
}

func TestCommandCreate(t *testing.T) {
	aliases := []string{"create", "new"}

	for _, alias := range aliases {
		t.Run(alias+" errors without arguments", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, nil, alias)

			if exitCode != 1 {
				t.Fatalf("expected exit code 1, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stderr, "Usage:") && !strings.Contains(stderr, "USAGE:") {
				t.Errorf("expected stderr to contain Usage, got:\n%s", stderr)
			}
		})

		t.Run(alias+" errors with filepath only", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, nil, alias, "test.kdbx")

			if exitCode != 1 {
				t.Fatalf("expected exit code 1, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stderr, "accepts 2 arg(s), received 1") {
				t.Errorf("expected stderr to contain args error, got:\n%s", stderr)
			}
		})

		t.Run(alias+" errors with too many args", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, nil, alias, "a", "b", "c")

			if exitCode != 1 {
				t.Fatalf("expected exit code 1, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stderr, "accepts 2 arg(s), received 3") {
				t.Errorf("expected stderr to contain args error, got:\n%s", stderr)
			}
		})

		t.Run(alias+" errors with existing file", func(t *testing.T) {
			_, stderr, exitCode := runKeydex(t, nil, alias, "README.md", "test")

			if exitCode != 1 {
				t.Fatalf("expected exit code 1, got %d. stderr: %s", exitCode, stderr)
			}
			if !strings.Contains(stderr, "exists") {
				t.Errorf("expected stderr to contain existing file error, got:\n%s", stderr)
			}
		})
	}
}

func TestCommandCreateWithPassphrase(t *testing.T) {
	originalReadSecret := cli.ReadSecret
	originalConfirm := cli.Confirm
	defer func() {
		cli.ReadSecret = originalReadSecret
		cli.Confirm = originalConfirm
	}()

	t.Run("creates database successfully", func(t *testing.T) {
		password := "test-create-password"
		cli.ReadSecret = func(prompt string) string { return password }
		cli.Confirm = func(prompt string) bool { return false }

		dbPath := filepath.Join(t.TempDir(), "new.kdbx")

		err := cmd.Create.RunE(cmd.Create, []string{dbPath, "TestVault"})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Fatal("expected database file to be created")
		}

		db, err := kdbx.OpenFromPath(dbPath, password, "")
		if err != nil {
			t.Fatalf("failed to open created database: %v", err)
		}

		rootGroup := db.GetRootGroup()
		if rootGroup == nil {
			t.Fatal("expected root group to exist")
		}
		if rootGroup.Name != "TestVault" {
			t.Errorf("expected root group name 'TestVault', got '%s'", rootGroup.Name)
		}
	})

	t.Run("errors on passphrase mismatch", func(t *testing.T) {
		callCount := 0
		cli.ReadSecret = func(prompt string) string {
			callCount++
			if callCount == 1 {
				return "password1"
			}
			return "password2"
		}
		cli.Confirm = func(prompt string) bool { return false }

		dbPath := filepath.Join(t.TempDir(), "mismatch.kdbx")

		err := cmd.Create.RunE(cmd.Create, []string{dbPath, "TestVault"})
		if err == nil {
			t.Fatal("expected error for passphrase mismatch, got nil")
		}
		if !strings.Contains(err.Error(), "mismatch") {
			t.Errorf("expected mismatch error, got: %v", err)
		}

		if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
			t.Error("expected no database file to be created on mismatch")
		}
	})

	t.Run("creates database with keyfile when confirmed", func(t *testing.T) {
		password := "test-create-password"
		cli.ReadSecret = func(prompt string) string { return password }
		cli.Confirm = func(prompt string) bool {
			return strings.Contains(prompt, "keyfile")
		}

		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "withkey.kdbx")
		keyPath := filepath.Join(tmpDir, "withkey-key.xml")

		err := cmd.Create.RunE(cmd.Create, []string{dbPath, "TestVault"})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Fatal("expected database file to be created")
		}

		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			t.Fatal("expected keyfile to be created")
		}

		db, err := kdbx.OpenFromPath(dbPath, password, keyPath)
		if err != nil {
			t.Fatalf("failed to open created database with keyfile: %v", err)
		}

		rootGroup := db.GetRootGroup()
		if rootGroup == nil {
			t.Fatal("expected root group to exist")
		}
		if rootGroup.Name != "TestVault" {
			t.Errorf("expected root group name 'TestVault', got '%s'", rootGroup.Name)
		}
	})

	t.Run("errors on empty passphrase", func(t *testing.T) {
		cli.ReadSecret = func(prompt string) string { return "" }
		cli.Confirm = func(prompt string) bool { return false }

		dbPath := filepath.Join(t.TempDir(), "empty.kdbx")

		err := cmd.Create.RunE(cmd.Create, []string{dbPath, "TestVault"})
		if err == nil {
			t.Fatal("expected error for empty passphrase, got nil")
		}
		if !strings.Contains(err.Error(), "empty") {
			t.Errorf("expected empty passphrase error, got: %v", err)
		}

		if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
			t.Error("expected no database file to be created on empty passphrase")
		}
	})

	t.Run("errors when keyfile already exists", func(t *testing.T) {
		password := "test-create-password"
		cli.ReadSecret = func(prompt string) string { return password }
		cli.Confirm = func(prompt string) bool {
			return strings.Contains(prompt, "keyfile")
		}

		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "with-existing-key.kdbx")
		keyPath := filepath.Join(tmpDir, "with-existing-key-key.xml")

		if err := os.WriteFile(keyPath, []byte("already here"), 0o600); err != nil {
			t.Fatalf("failed to create existing keyfile: %v", err)
		}

		err := cmd.Create.RunE(cmd.Create, []string{dbPath, "TestVault"})
		if err == nil {
			t.Fatal("expected error for existing keyfile, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("expected existing keyfile error, got: %v", err)
		}
		if !strings.Contains(err.Error(), keyPath) {
			t.Errorf("expected error to mention keyfile path, got: %v", err)
		}
		if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
			t.Error("expected no database file to be created when keyfile already exists")
		}
	})
}

func createFixtureDB2() error {
	file, err := os.Create(fixtureDB2)
	if err != nil {
		return err
	}

	db, err := kdbx.NewFromFile(file)
	if err != nil {
		return err
	}
	if err := db.SetPasswordAndKey(fixturePassword2, ""); err != nil {
		return err
	}

	rootGroup := db.NewGroup("TestDB2")
	socialGroup := db.NewGroup("Social")

	twitter := gokeepasslib.NewEntry()
	twitter.Values = append(twitter.Values,
		gokeepasslib.ValueData{
			Key:   "Title",
			Value: gokeepasslib.V{Content: "Twitter"}},
		gokeepasslib.ValueData{
			Key: "Password",
			Value: gokeepasslib.V{
				Content:   "twpass",
				Protected: wrappers.NewBoolWrapper(true)}},
	)

	socialGroup.Entries = append(socialGroup.Entries, twitter)
	rootGroup.Groups = append(rootGroup.Groups, *socialGroup)
	db.Content.Root.Groups = []gokeepasslib.Group{*rootGroup}

	return db.Save()
}

func TestCommandDiff(t *testing.T) {
	t.Run("diff identical archives outputs nothing", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": fixturePassword,
		}, "diff", fixtureDB, fixtureDB)

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if stdout != "" {
			t.Errorf("expected empty output, got:\n%s", stdout)
		}
	})

	t.Run("diff different archives shows added and removed entries", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": fixturePassword2,
		}, "diff", fixtureDB, fixtureDB2)

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "-/") {
			t.Errorf("expected removed entries in output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "+/") {
			t.Errorf("expected added entries in output, got:\n%s", stdout)
		}
	})

	t.Run("output includes file name headers", func(t *testing.T) {
		absA, _ := filepath.Abs(fixtureDB)
		absB, _ := filepath.Abs(fixtureDB2)

		stdout, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": fixturePassword2,
		}, "diff", absA, absB)

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "--- "+fixtureDB) {
			t.Errorf("expected --- header with base name in output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "+++ "+fixtureDB2) {
			t.Errorf("expected +++ header with base name in output, got:\n%s", stdout)
		}
		if strings.Contains(stdout, absA) {
			t.Errorf("expected full path to be absent from output, got:\n%s", stdout)
		}
		if strings.Contains(stdout, absB) {
			t.Errorf("expected full path to be absent from output, got:\n%s", stdout)
		}
	})

	t.Run("output includes entry count line", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": fixturePassword2,
		}, "diff", fixtureDB, fixtureDB2)

		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d. stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "@@ -") {
			t.Errorf("expected @@ line in output, got:\n%s", stdout)
		}
	})

	t.Run("fails with wrong password for first archive", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": "wrong-password",
			"KEYDEX_PASSPHRASE_B": fixturePassword2,
		}, "diff", fixtureDB, fixtureDB2)

		if exitCode == 0 {
			t.Fatalf("expected non-zero exit code, got 0. stdout: %s", stdout)
		}
		if !strings.Contains(stderr, "Wrong password?") {
			t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
		}
	})

	t.Run("fails with wrong password for second archive", func(t *testing.T) {
		stdout, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": "wrong-password",
		}, "diff", fixtureDB, fixtureDB2)

		if exitCode == 0 {
			t.Fatalf("expected non-zero exit code, got 0. stdout: %s", stdout)
		}
		if !strings.Contains(stderr, "Wrong password?") {
			t.Errorf("expected 'Wrong password?' in stderr, got:\n%s", stderr)
		}
	})

	t.Run("fails with non-existing first archive", func(t *testing.T) {
		_, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": fixturePassword2,
		}, "diff", "non-existent.kdbx", fixtureDB2)

		if exitCode == 0 {
			t.Fatal("expected non-zero exit code")
		}
		if !strings.Contains(stderr, "no such file or directory") {
			t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
		}
	})

	t.Run("fails with non-existing second archive", func(t *testing.T) {
		_, stderr, exitCode := runKeydex(t, map[string]string{
			"KEYDEX_PASSPHRASE_A": fixturePassword,
			"KEYDEX_PASSPHRASE_B": fixturePassword2,
		}, "diff", fixtureDB, "non-existent.kdbx")

		if exitCode == 0 {
			t.Fatal("expected non-zero exit code")
		}
		if !strings.Contains(stderr, "no such file or directory") {
			t.Errorf("expected 'no such file or directory' in stderr, got:\n%s", stderr)
		}
	})
}

func TestCommandOpenWithPassphrase(t *testing.T) {
	originalReadSecret := cli.ReadSecret
	defer func() {
		cli.ReadSecret = originalReadSecret
	}()

	// Ensure env var is unset so GetPassphrase falls through to ReadSecret
	t.Setenv(cmd.ENV_PASSPHRASE, "")

	t.Run("opens database with correct passphrase from prompt", func(t *testing.T) {
		cli.ReadSecret = func(prompt string) string { return fixturePassword }

		err := cmd.Open.RunE(cmd.Open, []string{fixtureDB, "/TestDB/Coding/NonExistent"})
		if err == nil {
			t.Fatal("expected 'Missing entry' error, got nil")
		}
		if !strings.Contains(err.Error(), "Missing entry") {
			t.Errorf("expected 'Missing entry' error, got: %v", err)
		}
	})

	t.Run("fails with wrong passphrase from prompt", func(t *testing.T) {
		cli.ReadSecret = func(prompt string) string { return "wrong-password" }

		err := cmd.Open.RunE(cmd.Open, []string{fixtureDB, "/TestDB/Coding/GitHub"})
		if err == nil {
			t.Fatal("expected error for wrong passphrase, got nil")
		}
		if !strings.Contains(err.Error(), "Wrong password?") {
			t.Errorf("expected 'Wrong password?' error, got: %v", err)
		}
	})
}
