package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fabiobrug/mako.git/internal/testutil"
)

func TestDetectGoProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	// Create go.mod
	testutil.TempFile(t, tmpDir, "go.mod", "module test\n\ngo 1.21\n")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "Go" {
		t.Errorf("Expected Go, got %s", pt.Language)
	}
	if pt.BuildTool != "go" {
		t.Errorf("Expected go build tool, got %s", pt.BuildTool)
	}
	if pt.TestCmd != "go test ./..." {
		t.Errorf("Expected 'go test ./...', got '%s'", pt.TestCmd)
	}
	if pt.RunCmd != "go run ." {
		t.Errorf("Expected 'go run .', got '%s'", pt.RunCmd)
	}
}

func TestDetectNodeProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	// Create package.json
	testutil.TempFile(t, tmpDir, "package.json", `{"name": "test"}`)
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "JavaScript/TypeScript" {
		t.Errorf("Expected JavaScript/TypeScript, got %s", pt.Language)
	}
	if pt.BuildTool != "npm" {
		t.Errorf("Expected npm, got %s", pt.BuildTool)
	}
}

func TestDetectNodeWithYarn(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "package.json", `{"name": "test"}`)
	testutil.TempFile(t, tmpDir, "yarn.lock", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.BuildTool != "yarn" {
		t.Errorf("Expected yarn, got %s", pt.BuildTool)
	}
	if pt.TestCmd != "yarn test" {
		t.Errorf("Expected 'yarn test', got '%s'", pt.TestCmd)
	}
}

func TestDetectNextJS(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "package.json", `{"name": "test"}`)
	testutil.TempFile(t, tmpDir, "next.config.js", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Framework != "Next.js" {
		t.Errorf("Expected Next.js framework, got %s", pt.Framework)
	}
}

func TestDetectPythonProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "requirements.txt", "flask==2.0.1\n")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "Python" {
		t.Errorf("Expected Python, got %s", pt.Language)
	}
	if pt.BuildTool != "pip" {
		t.Errorf("Expected pip, got %s", pt.BuildTool)
	}
}

func TestDetectDjangoProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "requirements.txt", "django==4.0\n")
	testutil.TempFile(t, tmpDir, "manage.py", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Framework != "Django" {
		t.Errorf("Expected Django framework, got %s", pt.Framework)
	}
	if pt.RunCmd != "python manage.py runserver" {
		t.Errorf("Expected Django run command, got '%s'", pt.RunCmd)
	}
}

func TestDetectPythonPoetry(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "pyproject.toml", "[tool.poetry]\n")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.BuildTool != "poetry" {
		t.Errorf("Expected poetry, got %s", pt.BuildTool)
	}
}

func TestDetectRustProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "Cargo.toml", "[package]\nname = \"test\"\n")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "Rust" {
		t.Errorf("Expected Rust, got %s", pt.Language)
	}
	if pt.BuildTool != "cargo" {
		t.Errorf("Expected cargo, got %s", pt.BuildTool)
	}
	if pt.TestCmd != "cargo test" {
		t.Errorf("Expected 'cargo test', got '%s'", pt.TestCmd)
	}
}

func TestDetectJavaWithMaven(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "pom.xml", "<project></project>")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "Java" {
		t.Errorf("Expected Java, got %s", pt.Language)
	}
	if pt.BuildTool != "maven" {
		t.Errorf("Expected maven, got %s", pt.BuildTool)
	}
}

func TestDetectJavaWithGradle(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "build.gradle", "plugins { }")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.BuildTool != "gradle" {
		t.Errorf("Expected gradle, got %s", pt.BuildTool)
	}
}

func TestDetectRubyProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "Gemfile", "source 'https://rubygems.org'")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "Ruby" {
		t.Errorf("Expected Ruby, got %s", pt.Language)
	}
	if pt.BuildTool != "bundler" {
		t.Errorf("Expected bundler, got %s", pt.BuildTool)
	}
}

func TestDetectRailsProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "Gemfile", "gem 'rails'")
	os.MkdirAll(filepath.Join(tmpDir, "config"), 0755)
	testutil.TempFile(t, tmpDir, "config/routes.rb", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Framework != "Rails" {
		t.Errorf("Expected Rails framework, got %s", pt.Framework)
	}
	if pt.RunCmd != "rails server" {
		t.Errorf("Expected 'rails server', got '%s'", pt.RunCmd)
	}
}

func TestDetectPHPProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "composer.json", `{"name": "test"}`)
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "PHP" {
		t.Errorf("Expected PHP, got %s", pt.Language)
	}
	if pt.BuildTool != "composer" {
		t.Errorf("Expected composer, got %s", pt.BuildTool)
	}
}

func TestDetectLaravelProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "composer.json", `{"name": "test"}`)
	testutil.TempFile(t, tmpDir, "artisan", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Framework != "Laravel" {
		t.Errorf("Expected Laravel framework, got %s", pt.Framework)
	}
}

func TestDetectElixirProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	testutil.TempFile(t, tmpDir, "mix.exs", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if pt.Language != "Elixir" {
		t.Errorf("Expected Elixir, got %s", pt.Language)
	}
	if pt.BuildTool != "mix" {
		t.Errorf("Expected mix, got %s", pt.BuildTool)
	}
}

func TestDetectNoProject(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	// Empty directory
	pt := DetectProjectType()
	if pt != nil {
		t.Error("Expected no project type for empty directory")
	}
}

func TestGetProjectHint(t *testing.T) {
	tests := []struct {
		name string
		pt   *ProjectType
		want string
	}{
		{
			name: "Go project",
			pt: &ProjectType{
				Language:  "Go",
				BuildTool: "go",
			},
			want: "Go project (go)",
		},
		{
			name: "Next.js project",
			pt: &ProjectType{
				Language:  "JavaScript/TypeScript",
				Framework: "Next.js",
				BuildTool: "npm",
			},
			want: "JavaScript/TypeScript Next.js project (npm)",
		},
		{
			name: "Nil project",
			pt:   nil,
			want: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pt.GetProjectHint()
			if got != tt.want {
				t.Errorf("GetProjectHint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindProjectRoot(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	
	// Create nested structure
	projectDir := filepath.Join(tmpDir, "project")
	subDir := filepath.Join(projectDir, "src", "components")
	os.MkdirAll(subDir, 0755)
	
	// Create project marker at root
	testutil.TempFile(t, projectDir, "go.mod", "module test")
	
	// Change to subdirectory
	os.Chdir(subDir)
	
	// Should find project root
	root := FindProjectRoot()
	if root != projectDir {
		t.Errorf("Expected root %s, got %s", projectDir, root)
	}
}

func TestFindProjectRootNoMarker(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	
	// Change to directory without project markers
	os.Chdir(tmpDir)
	
	// Should return current directory
	root := FindProjectRoot()
	if root != tmpDir {
		t.Errorf("Expected current dir %s, got %s", tmpDir, root)
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	
	// Create test file
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)
	
	if !fileExists(filePath) {
		t.Error("Expected file to exist")
	}
	
	if fileExists(filepath.Join(tmpDir, "nonexistent.txt")) {
		t.Error("Expected file not to exist")
	}
	
	// Directory should not count as file
	if fileExists(tmpDir) {
		t.Error("Directory should not be detected as file")
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	
	// Create test directory
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	
	if !dirExists(subDir) {
		t.Error("Expected directory to exist")
	}
	
	if dirExists(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("Expected directory not to exist")
	}
	
	// File should not count as directory
	filePath := filepath.Join(tmpDir, "file.txt")
	os.WriteFile(filePath, []byte("test"), 0644)
	if dirExists(filePath) {
		t.Error("File should not be detected as directory")
	}
}

func TestConfigFilesDetection(t *testing.T) {
	tmpDir := testutil.TempDir(t)
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	// Create go project with go.sum
	testutil.TempFile(t, tmpDir, "go.mod", "module test")
	testutil.TempFile(t, tmpDir, "go.sum", "")
	
	pt := DetectProjectType()
	if pt == nil {
		t.Fatal("Expected project type to be detected")
	}
	
	if len(pt.ConfigFiles) != 2 {
		t.Errorf("Expected 2 config files, got %d", len(pt.ConfigFiles))
	}
	
	hasGoMod := false
	hasGoSum := false
	for _, file := range pt.ConfigFiles {
		if file == "go.mod" {
			hasGoMod = true
		}
		if file == "go.sum" {
			hasGoSum = true
		}
	}
	
	if !hasGoMod || !hasGoSum {
		t.Errorf("Expected both go.mod and go.sum in config files, got %v", pt.ConfigFiles)
	}
}

func BenchmarkDetectProjectType(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpDir)
	
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)
	
	// Create test project files
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name": "test"}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "next.config.js"), []byte(""), 0644)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectProjectType()
	}
}

func BenchmarkFindProjectRoot(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "mako-bench-*")
	defer os.RemoveAll(tmpDir)
	
	// Create nested structure
	subDir := filepath.Join(tmpDir, "a", "b", "c", "d")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	
	oldDir, _ := os.Getwd()
	os.Chdir(subDir)
	defer os.Chdir(oldDir)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindProjectRoot()
	}
}
