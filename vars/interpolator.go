package vars

import (
	"regexp"
	"strings"
)

func Interpolate(input string, ctx RuntimeContext) string {
	re := regexp.MustCompile(`\$\{?(\w+)\}?`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		key := strings.TrimPrefix(match, "$")
		switch key {
		case "FILE":
			return ctx.File
		case "CWD":
			return ctx.CWD
		case "EVENT_TYPE":
			return ctx.EventType
		case "EXT":
			return ctx.Extension
		case "BASENAME":
			return ctx.Basename
		case "DIRNAME":
			return ctx.Dirname
		case "RELFILE":
			return ctx.RelFile
		case "BUILD_FILE":
			return ctx.BuildFile
		case "BUILD_STEP":
			return ctx.StepName
		case "TIMESTAMP":
			return ctx.Timestamp
		case "UUID":
			return ctx.UUID
		case "OS":
			return ctx.OS
		case "ARCH":
			return ctx.Arch
		case "GIT_BRANCH":
			return ctx.GitBranch
		case "GIT_COMMIT":
			return ctx.GitCommit
		default:
			if val, ok := ctx.UserVars[key]; ok {
				return val
			}
		}
		return match
	})
}
