package vars

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RuntimeContext struct {
	File      string
	EventType string
	CWD       string
	Extension string
	Basename  string
	Dirname   string
	RelFile   string
	BuildFile string
	StepName  string
	Timestamp string
	UUID      string
	OS        string
	Arch      string
	GitBranch string
	GitCommit string

	UserVars map[string]string
}

func BuildContext(userVars map[string]string, stepName, buildFile string, changedFile string, eventType string) RuntimeContext {
	cwd, _ := os.Getwd()

	relPath, _ := filepath.Rel(cwd, changedFile)
	ext := filepath.Ext(changedFile)
	base := strings.TrimSuffix(filepath.Base(changedFile), ext)
	dir := filepath.Dir(changedFile)

	branch := runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	commit := runGitCommand("rev-parse", "HEAD")

	return RuntimeContext{
		File:      changedFile,
		EventType: eventType,
		CWD:       cwd,
		Extension: ext,
		Basename:  base,
		Dirname:   dir,
		RelFile:   relPath,

		BuildFile: buildFile,
		StepName:  stepName,

		Timestamp: time.Now().Format(time.RFC3339),
		UUID:      uuid.New().String(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,

		GitBranch: strings.TrimSpace(branch),
		GitCommit: strings.TrimSpace(commit),

		UserVars: userVars,
	}
}

func runGitCommand(args ...string) string {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}
