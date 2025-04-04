# CICD-Go 🚀
A lightweight, fast, and configurable CLI-based CI/CD tool for running builds, tests, and deployments.

## 📌 Features
- **Runs Build Pipelines by Default**: Just execute the command in your project directory.
- **Directory Watching**: Use `-watch` to automatically trigger builds on file changes.
- **Custom Build Steps**: Define commands in `build.yaml`.
- **Parallel Execution**: Run multiple steps simultaneously.
- **Conditional Steps**: Only execute steps if conditions are met.

## 🔧 Installation
### Linux/macOS
```sh
go build -o cicd-go ./cmd
mv cicd-go /usr/local/bin/
```

### Windows
1. Run:
```sh
go build -o cicd-go.exe ./cmd
```
2. Move cicd-go.exe to a folder in your system's $PATH.

## 🚀 Usage

### Running the pipeline (default behavior)
```sh
cicd-go
```

### Watching for file changes and re-running the pipeline
```sh
cicd-go -watch
```
## 🏁 Flags
### -file=<path>
Used to specify a custom build-file.
### -watch
Watches the current directory for file-changes and runs the pipeline.

## ⚙️ Example build.yaml
```yaml
# This is used to build cicd-go
version: 1.0
setup:
  - name: "Install dependencies"
    command: "go mod tidy"

steps:
  - name: "Build binary"
    command: "go build -o cicd-go ./cmd"

post_build:
  - name: "Moving binary"
    command: "mv cicd-go output/cicd-go"
  - name: "Copying binary to global path"
    command: "cp output/cicd-go /usr/local/bin"
    if: "env.GLOBAL == 'true'"

ignore:
  - "output/"
  - ".git"
  - "*.log"
```

Example with parallel steps:
```yaml
version: 1.0

parallel:
  - name: "Run php-cs-fixer"
    command: "php vendor/bin/php-cs-fixer fix --dry-run"
  - name: "Run psalm"
    command: "php vendor/bin/psalm"
```

