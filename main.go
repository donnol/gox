package main

import (
	"fmt"
	"os"
	"strings"
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

	// 检查 gopath 是否存在相关库
	var exist bool
	for _, gopath := range gopaths {
		fmt.Println(gopath)
	}

	if !exist {
		// 获取
	} else {
		// 更新
	}
}
