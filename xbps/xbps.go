package xbps

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type Pkg struct {
	Architecture  string
	BuildDate     string
	Sha256        string
	FileSize      string
	Homepage      string
	InstalledSize string
	License       string
	Maintainer    string
	PkgName       string
	PkgVersion    string
	Repository    string
	RunDepends    []string
	ShlibReqs     []string
	Desc          string
	SourceRev     string
}

var PkgRegex = [...]string{
	`architecture: (\S*)\n`,
	`build-date: (\d{4}-\d{2}-\d{2} \d{2}:\d{2} [A-Za-z]{3})\n`,
	`filename-sha256: ([a-z\d]{0,64})\n`,
	`filename-size: (\S*)\n`,
	`homepage: (.*)\n`,
	`installed_size: (\S*)\n`,
	`license: (.*)\n`,
	`maintainer: (.*)\n`,
	`pkgname: (\S*)\n`,
	`pkgver: \S*-(\d*\S*)\n`,
	`repository: (.*)\n`,
	`run_depends:\n((?:\t.*\d\n)*)`,
	`shlib-requires:\n((?:\t.*\d\n)*)`,
	`short_desc: (.*)\n`,
	`source-revisions: (.*)`,
}

var cmpPkgRegex []*regexp.Regexp

func init() {
	for _, reg := range PkgRegex {
		r, _ := regexp.Compile(reg)
		cmpPkgRegex = append(cmpPkgRegex, r)
	}
}

// Install package with xbps
func (p *Pkg) Install() error {
	fmt.Printf("Installing pkg: %+v\n", p.PkgName)
	return nil
}

// Query to find suitable packages with xbps
func Query(name string) ([]string, error) {
	out, err := exec.Command("xrs", name).Output()
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}
	var lines []string
	lines = append(lines, strings.Split(string(out), "\n")...)

	var pkgNames []string

	for _, pkgLine := range lines {
		if pkgLine == "" {
			break
		}
		r, _ := regexp.Compile(`\[([*-])\] ([A-Za-z-\d_\+]+[-32bitA-Za-z\d._]*)\s*(.*)`)
		matches := r.FindStringSubmatch(pkgLine)

		pkgNames = append(pkgNames, matches[2])
	}
	sort.Strings(pkgNames)
	return pkgNames, nil

}

// Info gets package information from xbps
func Info(name string) (Pkg, error) {
	out, err := exec.Command("xbps-query", "-RS", name).Output()
	if err != nil {
		return Pkg{}, err
	}

	var p Pkg

	values := reflect.ValueOf(&p).Elem()

	for i, r := range cmpPkgRegex {
		m := r.FindStringSubmatch(string(out))
		var v string

		if len(m) == 0 {
			v = "n/a"
			continue
		} else {
			v = m[1]
		}
		value := values.Field(i)

		switch value.Kind() {
		case reflect.String:
			v := strings.ReplaceAll(v, "   ", "")
			value.SetString(v)
		case reflect.Slice:
			p := parseToList(v)
			slice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(p), cap(p))
			for i := 0; i < len(p); i++ {
				val := slice.Index(i)
				val.SetString(p[i])
			}
			value.Set(slice)
		default:
			err := errors.New("Type not known")
			log.Fatal(err)
			return Pkg{}, err
		}
	}

	return p, nil
}

// Parse list of deps and shlibs to list of strings
func parseToList(s string) []string {
	n := strings.ReplaceAll(s, "\t", "")
	split := strings.Split(n, "\n")
	return split

}

const Tmpl = `=============================
    Package {{ .PkgName }}
=============================
       Version: {{ .PkgVersion}}		
    Build date: {{ .BuildDate }}
    Repository: {{ .Repository }}
   Description: {{ .Desc }}
       License: {{ .License }}
---------------
  Architecture: {{ .Architecture }}
  Dependencies: {{range $val := .RunDepends}}
                {{$val}}{{end}}
   Shared libs: {{range $val := .ShlibReqs}}
                {{$val}}{{end}}
---------------
        Sha256: {{ .Sha256 }}
---------------
    Maintainer: {{ .Maintainer }}
      Homepage: {{ .Homepage }}
---------------
 Download size: {{ .FileSize }}
Installed size:	{{ .InstalledSize }}
---------------

`
