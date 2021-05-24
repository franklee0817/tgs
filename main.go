package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var mainFuncStr = `package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(impl.APIImpl)
	err := imp.Init()
	if err != nil {
		fmt.Printf("apiImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(tarsfile.Api)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".obj")

	// Run application
	tars.Run()
}`

func main() {
	var mod = flag.String("m", "", "specify go mod, eg: github.com/franklee0817/tgs")
	var file = flag.String("f", "", "specify tars file, eg: test.tars")
	flag.Parse()
	if len(*mod) == 0 || !strings.HasSuffix(*file, "tars") {
		flag.Usage()
		return
	}

	s := strings.Split(*mod, "/")
	baseDir := s[len(s)-1]
	err := createPrj(baseDir)
	if err != nil {
		fmt.Printf("failed to create project:%v\n", err)
		return
	}

	err = parseTars(*mod, *file, baseDir+"/protocol")
	if err != nil {
		fmt.Printf("failed to parse tars file:%v\n", err)
		return
	}
	err = createMainFile(baseDir)
	if err != nil {
		fmt.Printf("failed to create main.go. you can paste the following code manually\n%s", mainFuncStr)
	}

	fmt.Println("succeed to create", baseDir)
}

func createMainFile(outDir string) error {
	file, err := os.OpenFile(outDir+"/main.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	fp := bufio.NewWriter(file)
	defer fp.Flush()
	_, err = fp.WriteString(mainFuncStr)
	if err != nil {
		return err
	}
	return nil
}

func parseTars(mod, file, outDir string) error {
	cmd := exec.Command("/bin/bash", "-c", "tars2go -module "+mod+" -outdir "+outDir+" "+file)
	byt, err := cmd.Output()
	fmt.Println(string(byt))
	if err == nil {
		cmd := exec.Command("/bin/bash", "-c", "cp "+file+" "+outDir+"/")
		_, _ = cmd.Output()
	}

	return err
}

// DirSpec 目录自述结构
type DirSpec struct {
	Path string
	Spec string
}

func createPrj(baseDir string) error {
	dirs := []DirSpec{
		{
			Path: "client",
			Spec: "这里是rpc请求的客户端代码封装",
		}, {
			Path: "config/mysql",
			Spec: "这里是mysql相关配置",
		}, {
			Path: "config/es",
			Spec: "这里是es相关配置",
		}, {
			Path: "config/redis",
			Spec: "这里是redis相关配置",
		}, {
			Path: "constant",
			Spec: "这里是常量信息",
		}, {
			Path: "data",
			Spec: "这里是项目使用到的所有数据类型和结构的封装",
		}, {
			Path: "impl",
			Spec: "这里是tars服务的接口实现",
		}, {
			Path: "model/es",
			Spec: "这里是es的数据库模板封装",
		}, {
			Path: "model/mysql",
			Spec: "这里是mysql的数据库模板封装",
		}, {
			Path: "model/redis",
			Spec: "这里是redis的模板封装",
		}, {
			Path: "protocol",
			Spec: "这里是tars协议生成的相关协议代码",
		}, {
			Path: "service",
			Spec: "这里是核心业务逻辑，请将业务逻辑相关的代码实现写在这里，impl作为controller层进行请求分发和出错处理",
		}, {
			Path: "tool",
			Spec: "这里是项目用到的工具类代码封装，请不要在这里写业务逻辑",
		},
	}
	shell := fmt.Sprintf(`mkdir -p %s; echo "%s" > %s/readme.md;`, baseDir, baseDir, baseDir)
	shell += "cd " + baseDir + ";"
	for _, dir := range dirs {
		shell += fmt.Sprintf(`mkdir -p %s; echo "%s" > %s/readme.md;`, dir.Path, dir.Spec, dir.Path)
	}
	cmd := exec.Command("/bin/bash", "-c", shell)
	byt, err := cmd.Output()

	fmt.Println(string(byt))
	return err
}
