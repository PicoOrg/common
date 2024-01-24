package strings

import "testing"

func Test_implement_FString(t *testing.T) {
	type args struct {
		template string
		vars     map[string]any
	}
	tests := []struct {
		name    string
		m       *implement
		args    args
		wantRet string
		wantErr bool
	}{
		{
			name: "test error",
			m:    &implement{},
			args: args{
				template: "Hello, $world!",
				vars:     map[string]any{},
			},
			wantRet: "",
			wantErr: true,
		},
		{
			name: "test true",
			m:    &implement{},
			args: args{
				template: "$HELLO, $$world!",
				vars:     map[string]any{"HELLO": "Hello"},
			},
			wantRet: "Hello, $world!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &implement{}
			gotRet, err := m.FString(tt.args.template, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("implement.FString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRet != tt.wantRet {
				t.Errorf("implement.FString() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
