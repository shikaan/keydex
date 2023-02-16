package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func makeCobraCmd(name string, keyFlag string) *cobra.Command {
	cmd := &cobra.Command{
		Use: name,
	}
	cmd.Flags().StringP("key", "k", keyFlag, "")
	cmd.SetArgs([]string{"--key", keyFlag})

	return cmd
}

func TestReadDatabaseArguments(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name          string
		args          args
		env           map[string]string
		wantDatabase  string
		wantReference string
		wantKey       string
	}{
		{"args: nil, env: database", args{makeCobraCmd("copy", ""), []string{}}, map[string]string{ENV_DATABASE: "database"}, "database", "", ""},
		{"args: database, env: nil", args{makeCobraCmd("copy", ""), []string{"database"}}, map[string]string{}, "database", "", ""},
		{"args: ref, env: database", args{makeCobraCmd("copy", ""), []string{"/ref"}}, map[string]string{ENV_DATABASE: "database"}, "database", "/ref", ""},
		{"args: database, ref, env: nil", args{makeCobraCmd("copy", ""), []string{"database", "/ref"}}, map[string]string{}, "database", "/ref", ""},
		{"args: database, ref, env: otherdb", args{makeCobraCmd("copy", ""), []string{"database", "/ref"}}, map[string]string{ENV_DATABASE: "otherdb"}, "database", "/ref", ""},
		{"args: ref, env: nil (failure case)", args{makeCobraCmd("copy", ""), []string{"/ref"}}, map[string]string{}, "/ref", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			gotDatabase, gotReference, gotKey := ReadDatabaseArguments(tt.args.cmd, tt.args.args)
			if gotDatabase != tt.wantDatabase {
				t.Errorf("ReadDatabaseArguments() gotDatabase = %v, want %v", gotDatabase, tt.wantDatabase)
			}
			if gotReference != tt.wantReference {
				t.Errorf("ReadDatabaseArguments() gotReference = %v, want %v", gotReference, tt.wantReference)
			}
			if gotKey != tt.wantKey {
				t.Errorf("ReadDatabaseArguments() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
