# Contributing to 4CTF

We welcome contributions to the 4CTF project! To ensure a smooth process for everyone, please follow the guidelines below.

## Table of Contents

- [Contributing to 4CTF](#contributing-to-4ctf)
  - [Table of Contents](#table-of-contents)
  - [Code of Conduct](#code-of-conduct)
  - [How Can I Contribute?](#how-can-i-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Features](#suggesting-features)
    - [Submitting Code](#submitting-code)
  - [Pull Request Guidelines](#pull-request-guidelines)
  - [Commit Messages](#commit-messages)
    - [Examples](#examples)
    - [Common Types](#common-types)

---

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Be respectful and collaborative!

---

## How Can I Contribute?

### Reporting Bugs

1. Search for existing issues before submitting a new one.
2. Provide detailed steps to reproduce the bug, your environment, and any relevant logs.

### Suggesting Features

1. Open an issue with a clear and descriptive title.
2. Explain the problem and how the feature would address it.

### Submitting Code

1. Fork the repository and clone your fork.
2. Create a new branch from `develop`.
3. Make your changes.
4. Open a pull request (PR) targeting the `develop` branch.

---

## Pull Request Guidelines

1. **Target Branch**: All PRs must point to the `develop` branch. Changes are merged into `main` only for releases.
2. **Branch Naming**: Use descriptive branch names, e.g., `feature/add-something` or `bugfix/fix-issue`.
3. **Description**: Include a clear description of your changes and link to related issues if applicable.
4. **Tests**: Ensure your code is covered by tests if possible.
5. **Review Process**: Wait for at least one approval from a maintainer before merging.

---

## Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification. Please structure your commit messages as follows:

```a
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Examples

- `feat(api): add endpoint for user registration`
- `fix(ui): resolve button alignment issue on mobile`
- `docs(readme): update contributing section`

### Common Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no code logic changes)
- `refactor`: Code refactoring (neither fixes a bug nor adds a feature)
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (e.g., updating dependencies)

---

Thank you for contributing to 4CTF! 🎉
