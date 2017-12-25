package errors

import (
	"runtime"
	"strings"
)

var (
	gopath    string
	gopathlen int
)

func init() {
	pc, file, _, ok := runtime.Caller(0)
	if !ok || file == "?" {
		return
	}
	fn := runtime.FuncForPC(pc)
	fnstart := strings.LastIndex(fn.Name(), ".")
	if fnstart < 0 {
		return
	}
	fnpkg := fn.Name()[:strings.LastIndex(fn.Name(), "errors.init")]
	fnpkgstart := strings.Index(file, fnpkg)
	if fnpkgstart < 0 {
		return
	}
	gopathlen = fnpkgstart
	gopath = file[:gopathlen]
}

func trimGOPATH(filename string) string {
	if strings.HasPrefix(filename, gopath) {
		return filename[gopathlen:]
	}
	return filename
}
