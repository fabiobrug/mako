class Mako < Formula
  desc "AI-native shell orchestrator"
  homepage "https://github.com/fabiobrug/mako"
  url "https://github.com/fabiobrug/mako/archive/v1.0.0.tar.gz"
  sha256 "" # Will be calculated during release
  license "MIT"
  head "https://github.com/fabiobrug/mako.git", branch: "main"

  depends_on "go" => :build
  depends_on "sqlite" => :build

  def install
    # Set version
    version_flag = "-X github.com/fabiobrug/mako.git/internal/config.CurrentVersion=#{version}"
    
    # Build main binary
    system "go", "build",
           "-tags", "fts5",
           "-ldflags", version_flag,
           "-o", "mako",
           "cmd/mako/main.go"
    
    # Build menu binary
    system "go", "build",
           "-o", "mako-menu",
           "cmd/mako-menu/main.go"
    
    # Install binaries
    bin.install "mako"
    bin.install "mako-menu"
    
    # Install shell completions
    bash_completion.install "packaging/completions/mako.bash" => "mako"
    zsh_completion.install "packaging/completions/mako.zsh" => "_mako"
    fish_completion.install "packaging/completions/mako.fish"
    
    # Install man page
    man1.install "docs/man/mako.1"
  end

  def post_install
    # Create ~/.mako directory
    mako_dir = "#{Dir.home}/.mako"
    mkdir_p mako_dir unless File.exist?(mako_dir)
    
    # Show setup message
    ohai "Mako installed successfully!"
    puts ""
    puts "Setup your Gemini API key:"
    puts "  mako config set api_key YOUR_KEY"
    puts ""
    puts "Get your key at: https://ai.google.dev/"
    puts ""
    puts "Quick start:"
    puts "  mako                    # Start Mako shell"
    puts "  mako ask \"list files\"   # Generate commands"
    puts "  mako help               # See all commands"
  end

  test do
    # Test version command
    assert_match "Mako", shell_output("#{bin}/mako version 2>&1", 0)
    
    # Test help command
    assert_match "Commands", shell_output("#{bin}/mako help 2>&1", 0)
  end
end
