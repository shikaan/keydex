package credentials

import "testing"

func TestGetPassphrase(t *testing.T) {
	type args struct {
		database   string
		passphrase string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"returns passed passphrase", args{database: "", passphrase: "phrase"}, "phrase"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPassphrase(tt.args.database, tt.args.passphrase); got != tt.want {
				t.Errorf("GetPassphrase() = %v, want %v", got, tt.want)
			}
		})
	}
}
