{ pkgs, lib, config, inputs, ... }:

{
  # Set environment variables
  env = {
    GO111MODULE = "on"; # Enable Go modules
    PATH = "${pkgs.go}/bin:$HOME/go/bin:$PATH"; # Add Go and Go tools to PATH
  };

  # Install required packages
  packages = [
    pkgs.git           # Git for version control
    pkgs.go            # Go programming language
    pkgs.golangci-lint # Go linter
  ];

  # Define reusable scripts
  scripts = {
    start-api.exec = ''
      go run back/main.go serve
    '';
    run-tests.exec = ''
      go test ./...
    '';
    lint-code.exec = ''
      golangci-lint run
    '';
    live-reload.exec = ''
      air
    '';
  };

  # Enter the development shell
  enterShell = ''
    git --version
    go version
    go install github.com/air-verse/air@latest
  '';

  # Define tasks for easy command execution
  tasks = {
    "serve".exec = "go run back/main.go serve";
    "test".exec = "go test ./...";
    "lint".exec = "golangci-lint run";
    "live-reload".exec = "air";
  };

  # Test configuration
  enterTest = ''
    echo "Running tests..."
    git --version | grep --color=auto "${pkgs.git.version}"
    go version | grep "go"
  '';

  # Pre-commit hooks (optional)
  pre-commit.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
