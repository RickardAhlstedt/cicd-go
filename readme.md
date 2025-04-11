# CICD-Go 🚀
A lightweight, fast, and configurable CLI-based CI/CD tool for running builds, tests, and deployments.

## 📌 Features
- **Runs Build Pipelines by Default**: Just execute the command in your project directory.
- **Directory Watching**: Use `-watch` to automatically trigger builds on file changes.
- **Custom Build Steps**: Define commands in `build.yaml`.
- **Parallel Execution**: Run multiple steps simultaneously.
- **Conditional Steps**: Only execute steps if conditions are met (with support for advanced operators).
- **Dynamic Variables**: Use environment variables, git information, and custom variables in your steps.

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
# Inherit: path/to/base-file.yaml to inherit from and override by the local build-file
version: 1.1

vars:
  GLOBAL: "true"

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
    if: "$GLOBAL == true"

ignore:
  - "output/**"
  - ".git/**"
  - "*.log"
  - "build.yaml"
  - "readme.md"
```

Example with parallel steps:
```yaml
version: 1.0

parallel:
  - name: "Run php-cs-fixer"
    command: "php vendor/bin/php-cs-fixer fix $FILE"
  - name: "Run psalm"
    command: "php vendor/bin/psalm"
```

## Variables and Operators

In build.yaml, you can define custom variables and use them in your steps.
Built-in Variables:
```
$FILE: The file being processed.
$CWD: The current working directory.
$EVENT_TYPE: The type of file change event (e.g., "write", "create").
$BASENAME: The base name of the file (without extension).
$EXT: The file extension.
$DIRNAME: The directory name of the file.
$RELFILE: The relative path of the file.
$BUILD_FILE: The build file being used.
$BUILD_STEP: The name of the current step.
$TIMESTAMP: The current timestamp.
$UUID: A unique identifier for the build.
$OS: The operating system (e.g., "linux", "windows").
$ARCH: The system architecture (e.g., "amd64", "arm").
$GIT_BRANCH: The current git branch.
$GIT_COMMIT: The current git commit.
```

## Custom Variables:

You can define custom variables under vars: in your build.yaml:
```yaml
vars:
  GLOBAL: "true"
```

## Conditional operators:

You can use operators in the if: field to control when a step is executed. Supported operators:
- `==`: Equal to.
- `!=`: Not equal to.
- `^=`: String starts with
- `$=`: String ends with
- `*=`: String contains
- `~=`: String matches regex

## 🆕 Custom ignore-rules
The binary now loads a local .gitignore as additional ignores and merges it to the local build-file

## 🆕 Inheritance
By including `inherit` in your local build-file, you can add additional variables, steps, and ignores from a template-file.  
They way this builds the final config is:
1. If the local build-file has an `inherit`-property and the program can read that file, we read that and begin with merging the template and our local build-file.
2. If the program finds a .gitignore in your current working directory, it reads that file and adds that content to its own ignore-list.

For a complete working example, please see [.test-files](/.test-files).