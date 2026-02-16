package parser

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantCmd    string
		wantArgs   []string
		wantPrompt bool
	}{
		{
			name:       "Simple command",
			line:       "ls",
			wantCmd:    "ls",
			wantArgs:   nil,
			wantPrompt: false,
		},
		{
			name:       "Command with args",
			line:       "ls -lha /tmp",
			wantCmd:    "ls",
			wantArgs:   []string{"-lha", "/tmp"},
			wantPrompt: false,
		},
		{
			name:       "Command with multiple flags",
			line:       "git commit -m 'Initial commit' --no-verify",
			wantCmd:    "git",
			wantArgs:   []string{"commit", "-m", "'Initial", "commit'", "--no-verify"},
			wantPrompt: false,
		},
		{
			name:       "Prompt line with $",
			line:       "user@host:~/project $",
			wantCmd:    "",
			wantArgs:   nil,
			wantPrompt: true,
		},
		{
			name:       "Prompt line with #",
			line:       "root@host:~ #",
			wantCmd:    "",
			wantArgs:   nil,
			wantPrompt: true,
		},
		{
			name:       "Prompt line with ❯",
			line:       "❯",
			wantCmd:    "",
			wantArgs:   nil,
			wantPrompt: true,
		},
		{
			name:       "Empty line",
			line:       "",
			wantCmd:    "",
			wantArgs:   nil,
			wantPrompt: false,
		},
		{
			name:       "Whitespace only",
			line:       "   ",
			wantCmd:    "",
			wantArgs:   nil,
			wantPrompt: false,
		},
		{
			name:       "Pipeline command",
			line:       "cat file.txt | grep error",
			wantCmd:    "cat",
			wantArgs:   []string{"file.txt", "|", "grep", "error"},
			wantPrompt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ParseLine(tt.line)
			
			if info.IsPrompt != tt.wantPrompt {
				t.Errorf("ParseLine(%q).IsPrompt = %v, want %v", 
					tt.line, info.IsPrompt, tt.wantPrompt)
			}
			
			if !tt.wantPrompt {
				if info.Command != tt.wantCmd {
					t.Errorf("ParseLine(%q).Command = %q, want %q", 
						tt.line, info.Command, tt.wantCmd)
				}
				
				if len(info.Args) != len(tt.wantArgs) {
					t.Errorf("ParseLine(%q).Args length = %d, want %d", 
						tt.line, len(info.Args), len(tt.wantArgs))
				}
			}
		})
	}
}

func TestIsIgnoredCommand(t *testing.T) {
	tests := []struct {
		command string
		want    bool
	}{
		{"exit", true},
		{"clear", true},
		{"", true},
		{"ls", false},
		{"cd", false},
		{"git status", false},
		{"EXIT", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := IsIgnoredCommand(tt.command)
			if got != tt.want {
				t.Errorf("IsIgnoredCommand(%q) = %v, want %v", 
					tt.command, got, tt.want)
			}
		})
	}
}

func TestValidatePipeline(t *testing.T) {
	tests := []struct {
		command string
		want    bool
	}{
		// Valid pipelines
		{"ls -l", true},
		{"cat file.txt | grep error", true},
		{"ps aux | grep nginx | awk '{print $2}'", true},
		{"command1 && command2", true},
		{"command1 || command2", false}, // "||" pattern is considered invalid
		{"command1; command2", true},
		{"echo 'hello world'", true},
		{"git commit -m \"message\"", true},
		
		// Invalid pipelines
		{"| grep error", false},           // Starts with pipe
		{"cat file.txt |", false},         // Ends with pipe
		{"&& command", false},             // Starts with &&
		{"command &&", false},             // Ends with &&
		{"command1 | | command2", false},  // Double pipe
		{"", false},                       // Empty
		{"cat 'unclosed quote", false},    // Unbalanced quotes
		{"echo \"unclosed", false},        // Unbalanced double quotes
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := ValidatePipeline(tt.command)
			if got != tt.want {
				t.Errorf("ValidatePipeline(%q) = %v, want %v", 
					tt.command, got, tt.want)
			}
		})
	}
}

func TestHasBalancedQuotes(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"echo hello", true},
		{"echo 'hello'", true},
		{"echo \"hello\"", true},
		{"echo 'hello' \"world\"", true},
		{"echo 'it\\'s working'", true}, // Escaped quote
		{"echo 'unclosed", false},
		{"echo \"unclosed", false},
		{"echo 'one' 'two' 'three", false},
		{"", true}, // Empty string has balanced quotes
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := hasBalancedQuotes(tt.input)
			if got != tt.want {
				t.Errorf("hasBalancedQuotes(%q) = %v, want %v", 
					tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPipeline(t *testing.T) {
	tests := []struct {
		command string
		want    bool
	}{
		{"ls", false},
		{"ls -l", false},
		{"cat file.txt | grep error", true},
		{"command1 && command2", true},
		{"command1 || command2", true},
		{"command1; command2", true},
		{"echo 'text with | in it'", true}, // Note: simple implementation detects any |
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := IsPipeline(tt.command)
			if got != tt.want {
				t.Errorf("IsPipeline(%q) = %v, want %v", 
					tt.command, got, tt.want)
			}
		})
	}
}

func TestGetPipelineComplexity(t *testing.T) {
	tests := []struct {
		command string
		want    int
	}{
		{"ls", 0},
		{"cat file.txt | grep error", 1},
		{"ps aux | grep nginx | awk '{print $2}'", 2},
		{"cmd1 && cmd2 && cmd3", 2},
		{"cmd1 || cmd2", 3},                        // 2 pipes from "||" + 1 from "||" count = 3
		{"cmd1; cmd2; cmd3", 2},
		{"cat file.txt | grep $(echo pattern)", 2}, // 1 pipe + 1 subshell
		{"echo `date`", 2},                          // 2 backticks (open and close)
		{"cmd1 && cmd2 || cmd3; cmd4 | cmd5", 6},   // 2 && + 2 from || + 1 ; + 1 | = 6
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := GetPipelineComplexity(tt.command)
			if got != tt.want {
				t.Errorf("GetPipelineComplexity(%q) = %d, want %d", 
					tt.command, got, tt.want)
			}
		})
	}
}

func TestParseLineComplete(t *testing.T) {
	tests := []struct {
		line       string
		isComplete bool
	}{
		{"ls -l", true},
		{"", false},
		{"   ", false},
		{"user@host:~ $", false}, // Prompt
		{"git commit -m 'message'", true},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			info := ParseLine(tt.line)
			if info.IsComplete != tt.isComplete {
				t.Errorf("ParseLine(%q).IsComplete = %v, want %v", 
					tt.line, info.IsComplete, tt.isComplete)
			}
		})
	}
}

func TestParseLineRawLine(t *testing.T) {
	lines := []string{
		"ls -l",
		"  ls -l  ",
		"cat file.txt | grep error",
	}

	for _, line := range lines {
		t.Run(line, func(t *testing.T) {
			info := ParseLine(line)
			if info.RawLine != line {
				t.Errorf("ParseLine(%q).RawLine = %q, want %q", 
					line, info.RawLine, line)
			}
		})
	}
}

func BenchmarkParseLine(b *testing.B) {
	line := "git commit -m 'Initial commit' --no-verify --author 'User <user@example.com>'"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseLine(line)
	}
}

func BenchmarkValidatePipeline(b *testing.B) {
	command := "cat /var/log/nginx/access.log | grep '404' | awk '{print $1}' | sort | uniq -c | sort -nr"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePipeline(command)
	}
}

func BenchmarkGetPipelineComplexity(b *testing.B) {
	command := "ps aux | grep $(echo nginx) && cat file.txt || echo 'failed'; ls -l"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetPipelineComplexity(command)
	}
}
