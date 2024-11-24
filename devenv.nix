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
      cd back && air && cd ..
    '';
    build-api.exec = ''
      docker build -t api -f ./docker/Dockerfile.back .
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
    refresh-env.exec = ''
      if [ -f ./back/.env ]; then
        export $(cat ./back/.env | xargs)
        printf "\033[0;32mEnvironment variables refreshed\033[0m\n"
      fi
      if [ -f .env ]; then
        export $(cat .env | xargs)
        printf "\033[0;32mEnvironment variables refreshed\033[0m\n"
      fi
    '';
    start-docker.exec = ''
      docker compose -f ./docker/docker-compose.dev.yml up -d --remove-orphans
    '';
  };

  # Enter the development shell
  enterShell = ''
    git --version
    go version
    go install github.com/air-verse/air@latest
    go install github.com/amacneil/dbmate@latest
    go install github.com/volatiletech/sqlboiler/v4@latest
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
    go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest

    alias refresh-env="source $(which refresh-env)"

    # .env
    if [ ! -f ./back/.env ]; then
      cp ./back/.env.example ./back/.env
      cp ./back/config.yaml.example ./back/config.yaml
      printf "\033[0;32mCopied .env.example to .env, please fill it yourself if necessary\033[0m\n"
    fi
    # load .env
    export $(cat ./back/.env | xargs)

    # check if docker is installed
    if ! command -v docker &> /dev/null
    then
        echo "Docker could not be found, please install it"
    fi

    # check if docker-compose is installed
    if ! command -v docker compose &> /dev/null
    then
        echo "Docker-compose could not be found, please install it"
    fi

    # start docker compose
    docker compose -f ./docker/docker-compose.dev.yml up -d --remove-orphans

    printf "\033[0;32mDocker containers are up and running\033[0m\n"
    echo -----------------------------------
    printf "Available commands:\n"
    printf "start-api: Start the API server\n"
    printf "run-tests: Run tests\n"
    printf "lint-code: Lint the code\n"
    printf "live-reload: Start live-reload server\n"
    printf "refresh-env: Refresh environment variables\n"
    printf "start-docker: Start docker compose\n"
    echo -----------------------------------
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
