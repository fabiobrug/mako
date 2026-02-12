package context

import (
	"os"
	"path/filepath"
	"strings"
)

// ProjectType represents the detected project type
type ProjectType struct {
	Language   string   // go, python, javascript, rust, etc.
	Framework  string   // node, django, flask, react, etc.
	BuildTool  string   // npm, pip, cargo, go, etc.
	TestCmd    string   // Command to run tests
	RunCmd     string   // Command to run the project
	BuildCmd   string   // Command to build
	ConfigFiles []string // Detected config files
}

// DetectProjectType analyzes the current directory for project indicators
func DetectProjectType() *ProjectType {
	// Check for specific project files
	if fileExists("go.mod") {
		return detectGoProject()
	}
	if fileExists("package.json") {
		return detectNodeProject()
	}
	if fileExists("requirements.txt") || fileExists("pyproject.toml") || fileExists("setup.py") {
		return detectPythonProject()
	}
	if fileExists("Cargo.toml") {
		return detectRustProject()
	}
	if fileExists("pom.xml") || fileExists("build.gradle") {
		return detectJavaProject()
	}
	if fileExists("Gemfile") {
		return detectRubyProject()
	}
	if fileExists("composer.json") {
		return detectPHPProject()
	}
	if fileExists("mix.exs") {
		return detectElixirProject()
	}

	return nil
}

func detectGoProject() *ProjectType {
	pt := &ProjectType{
		Language:  "Go",
		BuildTool: "go",
		TestCmd:   "go test ./...",
		RunCmd:    "go run .",
		BuildCmd:  "go build",
		ConfigFiles: []string{"go.mod"},
	}

	if fileExists("go.sum") {
		pt.ConfigFiles = append(pt.ConfigFiles, "go.sum")
	}

	return pt
}

func detectNodeProject() *ProjectType {
	pt := &ProjectType{
		Language:    "JavaScript/TypeScript",
		BuildTool:   "npm",
		TestCmd:     "npm test",
		RunCmd:      "npm start",
		BuildCmd:    "npm run build",
		ConfigFiles: []string{"package.json"},
	}

	// Check for yarn
	if fileExists("yarn.lock") {
		pt.BuildTool = "yarn"
		pt.TestCmd = "yarn test"
		pt.RunCmd = "yarn start"
		pt.BuildCmd = "yarn build"
		pt.ConfigFiles = append(pt.ConfigFiles, "yarn.lock")
	}

	// Check for pnpm
	if fileExists("pnpm-lock.yaml") {
		pt.BuildTool = "pnpm"
		pt.TestCmd = "pnpm test"
		pt.RunCmd = "pnpm start"
		pt.BuildCmd = "pnpm build"
		pt.ConfigFiles = append(pt.ConfigFiles, "pnpm-lock.yaml")
	}

	// Detect framework
	if fileExists("next.config.js") || fileExists("next.config.ts") {
		pt.Framework = "Next.js"
	} else if dirExists("src") && fileExists("src/App.tsx") {
		pt.Framework = "React"
	} else if fileExists("angular.json") {
		pt.Framework = "Angular"
	} else if fileExists("vue.config.js") {
		pt.Framework = "Vue"
	}

	return pt
}

func detectPythonProject() *ProjectType {
	pt := &ProjectType{
		Language:  "Python",
		BuildTool: "pip",
		TestCmd:   "pytest",
		RunCmd:    "python main.py",
		BuildCmd:  "python setup.py build",
		ConfigFiles: []string{},
	}

	if fileExists("requirements.txt") {
		pt.ConfigFiles = append(pt.ConfigFiles, "requirements.txt")
	}
	if fileExists("pyproject.toml") {
		pt.ConfigFiles = append(pt.ConfigFiles, "pyproject.toml")
		pt.BuildTool = "poetry"
		pt.TestCmd = "poetry run pytest"
		pt.RunCmd = "poetry run python main.py"
	}
	if fileExists("Pipfile") {
		pt.ConfigFiles = append(pt.ConfigFiles, "Pipfile")
		pt.BuildTool = "pipenv"
		pt.TestCmd = "pipenv run pytest"
	}

	// Detect framework
	if fileExists("manage.py") {
		pt.Framework = "Django"
		pt.RunCmd = "python manage.py runserver"
		pt.TestCmd = "python manage.py test"
	} else if fileExists("app.py") || fileExists("application.py") {
		pt.Framework = "Flask"
		pt.RunCmd = "flask run"
	}

	return pt
}

func detectRustProject() *ProjectType {
	return &ProjectType{
		Language:    "Rust",
		BuildTool:   "cargo",
		TestCmd:     "cargo test",
		RunCmd:      "cargo run",
		BuildCmd:    "cargo build",
		ConfigFiles: []string{"Cargo.toml"},
	}
}

func detectJavaProject() *ProjectType {
	pt := &ProjectType{
		Language:  "Java",
		TestCmd:   "mvn test",
		BuildCmd:  "mvn package",
		RunCmd:    "java -jar target/*.jar",
		ConfigFiles: []string{},
	}

	if fileExists("pom.xml") {
		pt.BuildTool = "maven"
		pt.ConfigFiles = append(pt.ConfigFiles, "pom.xml")
	}
	if fileExists("build.gradle") {
		pt.BuildTool = "gradle"
		pt.TestCmd = "gradle test"
		pt.BuildCmd = "gradle build"
		pt.RunCmd = "gradle run"
		pt.ConfigFiles = append(pt.ConfigFiles, "build.gradle")
	}

	return pt
}

func detectRubyProject() *ProjectType {
	pt := &ProjectType{
		Language:    "Ruby",
		BuildTool:   "bundler",
		TestCmd:     "rake test",
		RunCmd:      "ruby main.rb",
		BuildCmd:    "bundle install",
		ConfigFiles: []string{"Gemfile"},
	}

	// Detect Rails
	if fileExists("config/routes.rb") {
		pt.Framework = "Rails"
		pt.RunCmd = "rails server"
		pt.TestCmd = "rails test"
	}

	return pt
}

func detectPHPProject() *ProjectType {
	pt := &ProjectType{
		Language:    "PHP",
		BuildTool:   "composer",
		TestCmd:     "phpunit",
		RunCmd:      "php -S localhost:8000",
		BuildCmd:    "composer install",
		ConfigFiles: []string{"composer.json"},
	}

	// Detect Laravel
	if fileExists("artisan") {
		pt.Framework = "Laravel"
		pt.RunCmd = "php artisan serve"
		pt.TestCmd = "php artisan test"
	}

	return pt
}

func detectElixirProject() *ProjectType {
	return &ProjectType{
		Language:    "Elixir",
		BuildTool:   "mix",
		TestCmd:     "mix test",
		RunCmd:      "mix phx.server",
		BuildCmd:    "mix compile",
		ConfigFiles: []string{"mix.exs"},
	}
}

// GetProjectHint returns a human-readable hint about the detected project
func (pt *ProjectType) GetProjectHint() string {
	if pt == nil {
		return ""
	}

	var parts []string
	
	if pt.Language != "" {
		parts = append(parts, pt.Language)
	}
	if pt.Framework != "" {
		parts = append(parts, pt.Framework)
	}
	if pt.BuildTool != "" {
		parts = append(parts, "project ("+pt.BuildTool+")")
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " ")
}

// Helper functions
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FindProjectRoot walks up the directory tree to find the project root
func FindProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	projectMarkers := []string{
		"go.mod", "package.json", "Cargo.toml", "requirements.txt",
		"pom.xml", "build.gradle", "Gemfile", "composer.json",
		".git",
	}

	for {
		for _, marker := range projectMarkers {
			if fileExists(filepath.Join(dir, marker)) {
				return dir
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	// Return current directory if no project root found
	currentDir, _ := os.Getwd()
	return currentDir
}
