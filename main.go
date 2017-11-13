package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var xList = map[string]string{
	"tools":      "https://github.com/golang/tools",
	"build":      "https://github.com/golang/build",
	"sys":        "https://github.com/golang/sys",
	"crypto":     "https://github.com/golang/crypto",
	"net":        "https://github.com/golang/net",
	"scratch":    "https://github.com/golang/scratch",
	"text":       "https://github.com/golang/text",
	"sync":       "https://github.com/golang/sync",
	"mobile":     "https://github.com/golang/mobile",
	"playground": "https://github.com/golang/playground",
	"arch":       "https://github.com/golang/arch",
	"tour":       "https://github.com/golang/tour",
	"image":      "https://github.com/golang/image",
	"review":     "https://github.com/golang/review",
	"exp":        "https://github.com/golang/exp",
	"time":       "https://github.com/golang/time",
	"perf":       "https://github.com/golang/perf",
	"debug":      "https://github.com/golang/debug",
}

func main() {
	// 获取 gopath
	var gopathEnv string = os.Getenv("GOPATH")
	gopaths := strings.Split(gopathEnv, ":")
	if len(gopaths) == 0 {
		panic("please set your gopath! like: export GOPATH=/your/gopath")
	}
	gopath := gopaths[0]
	fmt.Println(gopath)

	// 创建目录$GOPATH/src/golang.org/x/
	fileDir := filepath.Join(gopath, "src", "golang.org", "x")
	_, err := os.Stat(fileDir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("创建目录失败: %s, err: %v\n", fileDir, err))
		}
	}

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
						fmt.Printf("clone gox failed, link: %s, err: %v\n", link, err)
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
					fmt.Printf("clone gox failed, link: %s, err: %v\n", link, err)
				}
			}(fileDir, link)
		} else {
			// 目录不为空，更新
			wg.Add(1)
			go func(file, link string) {
				defer wg.Done()
				err = pullGox(file, link)
				if err != nil {
					fmt.Printf("pull gox failed, link: %s, err: %v\n", link, err)
				}
			}(file, link)
		}
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
	fmt.Printf("begin clone '%s'\n", link)
	cmd := exec.Command("git", "clone", link)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s\n", output)
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
	fmt.Printf("begin pull '%s'\n", link)
	cmd := exec.Command("git", "pull")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s\n", output)
		return err
	}
	fmt.Printf("finish pull '%s'\n", link)

	return nil
}
