# Backend Development

## Environment Setup

### Prerequisites

1. Install [Devenv](https://devenv.sh/) on your system.
2. Ensure you have a compatible package manager (e.g., `nix`).

### Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/4t-24/4ctf && cd 4ctf
    ```

2. Start the development environment:

    ```bash
    devenv up
    ```

3. Open a development shell:

    ```bash
    devenv shell
    ```

---

## Features and Usage

### Environment Variables

#### Installed Packages

- `git` for version control.
- `go` for backend development.
- `golangci-lint` for linting Go code.

#### Installed Tools

- `air`: A live-reloading tool for Go development, installed automatically during setup.

---

## Development Shell

You can enter the development shell by using  : `devenv shell`.  
This allows you to have all the necessary tools to develop properly.
