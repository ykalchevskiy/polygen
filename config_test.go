package main

import "testing"

func Test_toKebabCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{s: ""},
			want: "",
		},
		{
			name: "1 word",
			args: args{s: "Test"},
			want: "test",
		},
		{
			name: "2 words",
			args: args{s: "TestCase"},
			want: "test-case",
		},
		{
			name: "3 words",
			args: args{s: "TestCaseExample"},
			want: "test-case-example",
		},
		{
			name: "number",
			args: args{s: "TestCase2"},
			want: "test-case-2",
		},
		{
			name: "number inside",
			args: args{s: "Test2Case"},
			want: "test-2-case",
		},
		{
			name: "numbers",
			args: args{s: "TestCase123"},
			want: "test-case-123",
		},
		{
			name: "abbreviation",
			args: args{s: "HTTPResponse"},
			want: "http-response",
		},
		{
			name: "abbreviations",
			args: args{s: "HTTPResponseCode"},
			want: "http-response-code",
		},
		{
			name: "mixed case",
			args: args{s: "TestHTTPResponseCode"},
			want: "test-http-response-code",
		},
		{
			name: "lowercase",
			args: args{s: "testcase"},
			want: "testcase",
		},
		{
			name: "lowercase first letter",
			args: args{s: "testCase"},
			want: "test-case",
		},
		{
			name: "1 letter work",
			args: args{s: "AResponse"},
			want: "a-response",
		},
		{
			name: "1 letter word with number",
			args: args{s: "AResponse2HTTP"},
			want: "a-response-2-http",
		},
		{
			name: "utf8 characters",
			args: args{s: "ПриветМир"},
			want: "привет-мир",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toKebabCase(tt.args.s); got != tt.want {
				t.Errorf("toKebabCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
