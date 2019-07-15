package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	xList = map[string]string{
		"tools":      "github.com/golang/tools",
		"build":      "github.com/golang/build",
		"sys":        "github.com/golang/sys",
		"crypto":     "github.com/golang/crypto",
		"net":        "github.com/golang/net",
		"scratch":    "github.com/golang/scratch",
		"text":       "github.com/golang/text",
		"sync":       "github.com/golang/sync",
		"mobile":     "github.com/golang/mobile",
		"playground": "github.com/golang/playground",
		"arch":       "github.com/golang/arch",
		"tour":       "github.com/golang/tour",
		"image":      "github.com/golang/image",
		"review":     "github.com/golang/review",
		"exp":        "github.com/golang/exp",
		"time":       "github.com/golang/time",
		"perf":       "github.com/golang/perf",
		"debug":      "github.com/golang/debug",
	}
	toolList = map[string]string{
		"gocode":     "github.com/nsf/gocode",
		"gogetdoc":   "github.com/zmb3/gogetdoc",
		"godoc":      "golang.org/x/tools/cmd/godoc",
		"godef":      "github.com/rogpeppe/godef",
		"guru":       "golang.org/x/tools/cmd/guru",
		"go-outline": "github.com/ramya-rao-a/go-outline",
		"go-symbols": "github.com/acroca/go-symbols",
		"gorename":   "golang.org/x/tools/cmd/gorename",
		"golint":     "github.com/golang/lint/golint",
		// "gometalinter":  "github.com/alecthomas/gometalinter",
		"gometalinter": "gopkg.in/alecthomas/gometalinter.v1",
		"goreturns":    "github.com/sqs/goreturns",
		"goimports":    "golang.org/x/tools/cmd/goimports",
		"gotests":      "github.com/cweill/gotests",
		"gopkgs":       "github.com/uudashr/gopkgs/cmd/gopkgs",
		"gomodifytags": "github.com/fatih/gomodifytags",
		"gotype-live":  "github.com/tylerb/gotype-live",
		"impl":         "github.com/josharian/impl",
		"dlv":          "github.com/derekparker/delve/cmd/dlv",
		"goplay":       "github.com/haya14busa/goplay/cmd/goplay",
		// "megacheck":     "honnef.co/go/tools/...",
		"go-langserver": "github.com/sourcegraph/go-langserver",
	}
)

func main() {
	// 获取 gopath
	var gopathEnv = os.Getenv("GOPATH")
	var gopaths []string
	switch runtime.GOOS {
	case "windows":
		gopaths = strings.Split(gopathEnv, ";")
	default:
		gopaths = strings.Split(gopathEnv, ":")
	}
	if len(gopaths) == 0 {
		panic("please set your gopath! like: export GOPATH=/your/gopath")
	}
	gopath := gopaths[0]
	fmt.Println(gopath)

	// 创建目录 $GOPATH/src/golang.org/x/
	fileDir := filepath.Join(gopath, "src", "golang.org", "x")
	_, err := os.Stat(fileDir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("创建目录失败: %s, err: %v\n", fileDir, err))
		}
	}

	// 下载 x 包
	wg := new(sync.WaitGroup)
	for key, link := range xList {
		file := filepath.Join(fileDir, key)

		// 是否存在相关库目录
		_, err = os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				// 目录不存在，添加
				wg.Add(1)
				go func(fileDir, link string) {
					defer wg.Done()
					err = cloneGox(fileDir, link)
					if err != nil {
						fmt.Printf("=== clone gox failed, link: %s, err: %v\n", link, err)
					}
				}(fileDir, link)
			}
		} else if fis, err := ioutil.ReadDir(file); err == nil && len(fis) == 0 {
			// 空目录，添加
			wg.Add(1)
			go func(fileDir, link string) {
				defer wg.Done()
				err = cloneGox(fileDir, link)
				if err != nil {
					fmt.Printf("=== clone gox failed, link: %s, err: %v\n", link, err)
				}
			}(fileDir, link)
		} else {
			// 目录不为空，更新
			wg.Add(1)
			go func(file, link string) {
				defer wg.Done()
				err = pullGox(file, link)
				if err != nil {
					fmt.Printf("=== pull gox failed, link: %s, err: %v\n", link, err)
				}
			}(file, link)
		}
	}

	wg.Wait()

	// 安装工具
	for _, tool := range toolList {
		wg.Add(1)
		go func(tool string) {
			defer wg.Done()
			err := downloadGox(tool)
			if err != nil {
				fmt.Printf("=== download %s failed, err: %v\n", tool, err)
				return
			}
		}(tool)
	}
	wg.Wait()
}

func cloneGox(file, link string) error {
	// cd [file]
	err := os.Chdir(file)
	if err != nil {
		return err
	}

	// git clone [link]
	link = "https://" + link
	fmt.Printf("begin clone '%s'\n", link)
	cmd := exec.Command("git", "clone", link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("finish clone '%s'\n", link)

	return nil
}

func pullGox(file, link string) error {
	// cd [file]
	err := os.Chdir(file)
	if err != nil {
		return err
	}

	// git pull
	link = "https://" + link
	fmt.Printf("begin pull '%s'\n", link)
	cmd := exec.Command("git", "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("finish pull '%s'\n", link)

	return nil
}

func downloadGox(link string) error {
	// go get -u [link]
	fmt.Printf("begin go get -u '%s'\n", link)
	cmd := exec.Command("go", "get", "-u", link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("finish go get -u '%s'\n", link)

	return nil
}
