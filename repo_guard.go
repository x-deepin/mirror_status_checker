package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func DebugOut(format string, args ...interface{}) {
	if *Debug {
		fmt.Printf(format, args...)
	}
}

// CalacGuardPath 根据dirs以及count的数量计算出合适放置
// guard的路径. dirs应该是已经排序好的顺序．
// 返回路径的数量 = Min(len(dirs),count)
func CalacGuardPath(dirs []string, count int) []string {
	n := len(dirs)
	if count >= n {
		return dirs
	}

	sep := n / count
	var r []string
	for i, p := range dirs {
		if i%sep != 0 {
			DebugOut("%d filter %q  (%d)\n", i, p, sep)
			continue
		}
		DebugOut("%d hit %q  (%d)\n", i, p, sep)
		r = append(r, p)
	}

	// feed the last one if need
	if n >= count && len(r) != count {
		r = append(r, dirs[n-1])
	}

	DebugOut("Total directory: %d; Wish Guard : %d (Actual: %d, Sep:%d)\n", n, count, len(r), sep)
	return r
}

func ListSubDir(root string) []string {
	fs, err := ioutil.ReadDir(root)
	if err != nil {

		return nil
	}

	var r = []string{root}
	for _, finfo := range fs {
		if !finfo.IsDir() {
			continue
		}
		self := path.Join(root, finfo.Name())
		r = append(r, self)
		r = append(r, ListSubDir(self)...)
	}
	return r
}

func MakeIndexFileName(base string) string { return path.Join(base, *GuardIndexName) }
func MakeGuardFileName(base string) string { return path.Join(base, *GuardFileName) }

func WriteGuards(baseDir string, guardPaths []string) {
	var indexContent []string

	for _, p := range guardPaths {
		gpath := MakeGuardFileName(p)
		err := ioutil.WriteFile(gpath, ([]byte)(TimeNow), 0644)
		if err != nil {
			DebugOut("WriteGuard %q failed: %v\n", p, err)
			continue
		}

		indexContent = append(indexContent, strings.TrimPrefix(gpath, baseDir))
	}

	d, err := json.Marshal(indexContent)
	if err != nil {
		fmt.Printf("WriteGuardIndex marshal %q failed: %v\n", indexContent, err)
	}
	err = ioutil.WriteFile(MakeIndexFileName(baseDir), d, 0644)
	if err != nil {
		fmt.Printf("WriteGuardIndex  failed: %v\n", err)
	}
}

func ParseBaseDir(subs []string) (string, error) {
	var base string = path.Dir(strings.TrimRight(subs[0], string(os.PathSeparator)))
	for _, s := range subs {
		s = strings.TrimRight(s, string(os.PathSeparator))
		if base != path.Dir(s) {
			return base, fmt.Errorf("Not same base directory %q != Dir(%q)", base, s)
		}
	}
	return base, nil
}

var (
	TimeNow        = fmt.Sprintf("%v", time.Now().Unix())
	GuardFileName  = flag.String("guard-file-name", "__GUARD__"+TimeNow, "the guard file name")
	GuardIndexName = flag.String("index-file-name", "__GUARD__INDEX__", "the guard index name under root directory")
	GuardCount     = flag.Int("count", 100, "number of guard file you wish to set. The realy guard files is not exactly same as the number")
	Debug          = flag.Bool("debug", false, "")

	CleanGuard = flag.Bool("clean-guard", false, "clean the guard files and exit")
)

func main() {
	flag.Usage = func() {
		fmt.Println(`
根据实际的文件目录，采用深度优先搜索指定的目录．并在合适的位置
安置$guard-file文件．
注意:所有-开头的选项必须在指定目前之前设置.
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(-1)
		return
	}

	subs := args[:]
	base, err := ParseBaseDir(subs)
	if err != nil {
		fmt.Println("目前目录需要有相同的上级目录", err)
		os.Exit(-1)
		return
	}

	if *CleanGuard {
		CleanGuards(base)
	} else {
		var targets []string
		for _, d := range subs {
			targets = append(targets, ListSubDir(d)...)
		}
		GenerateGuards(base, targets)
	}
}

func DecodeIndex(r io.Reader) ([]string, error) {
	d := json.NewDecoder(r)
	var lines []string
	err := d.Decode(&lines)
	return lines, err
}

func CleanGuards(base string) {
	p := MakeIndexFileName(base)
	r, err := os.Open(p)
	if err != nil {
		panic("Can't find index file")
	}
	defer r.Close()

	lines, err := DecodeIndex(r)
	if err != nil {
		panic("Invalid Index File " + err.Error())
	}

	err = os.Remove(p)
	if err == nil {
		DebugOut("Removed %q\n", p)
	}
	for _, l := range lines {
		p = path.Join(base, l)
		err = os.Remove(p)
		if err == nil {
			DebugOut("Removed %q\n", p)
		}
	}
}

func GenerateGuards(base string, targets []string) {
	guards := CalacGuardPath(targets, *GuardCount)
	WriteGuards(base, guards)
}
