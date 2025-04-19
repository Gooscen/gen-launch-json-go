package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Request    string `json:"request"`
	Mode       string `json:"mode"`
	Program    string `json:"program"`
	BuildFlags string `json:"buildFlags"`
	Console    string `json:"console"`
}

type Compound struct {
	Name           string   `json:"name"`
	Configurations []string `json:"configurations"`
}

type LaunchJson struct {
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
	Compounds      []Compound      `json:"compounds"`
}

func main() {
	var configs []Configuration
	var compoundList []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}

		files, _ := filepath.Glob(filepath.Join(path, "*.go"))
		mainFuncFile := make([]string, 0)

		for _, file := range files {
			/*if file == "gen_launch_json.go"  {
				continue
			}*/
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			content := string(data)
			if strings.Contains(content, "package main") && strings.Contains(content, "func main()") {
				mainFuncFile = append(mainFuncFile, file)

			}
		}

		if len(mainFuncFile) == 1 {
			name := "Debug " + path
			configs = append(configs, Configuration{
				Name:       name,
				Type:       "go",
				Request:    "launch",
				Mode:       "debug",
				Program:    "${workspaceFolder}/" + path,
				BuildFlags: "-gcflags=all='-N -l'",
				Console:    "integratedTerminal",
			})
			compoundList = append(compoundList, name)
		} else if len(mainFuncFile) > 1 {
			for _, file := range mainFuncFile {
				name := "Debug " + path + "/" + file
				configs = append(configs, Configuration{
					Name:       name,
					Type:       "go",
					Request:    "launch",
					Mode:       "debug",
					Program:    "${workspaceFolder}/" + path + file,
					BuildFlags: "-gcflags=all='-N -l'",
					Console:    "integratedTerminal",
				})
				compoundList = append(compoundList, name)
			}

		}

		return nil
	})

	if err != nil {
		fmt.Println("❌ 遍历目录出错:", err)
		return
	}

	if len(configs) == 0 {
		fmt.Println("⚠️ 没有找到包含 main() 的 package main 目录。")
		return
	}

	launch := LaunchJson{
		Version:        "0.2.0",
		Configurations: configs,
		Compounds: []Compound{
			{
				Name:           "Debug All Main Packages",
				Configurations: compoundList,
			},
		},
	}

	output, err := json.MarshalIndent(launch, "", "  ")
	if err != nil {
		fmt.Println("❌ JSON 序列化失败:", err)
		return
	}

	if _, err := os.Stat(".vscode"); os.IsNotExist(err) {
		os.Mkdir(".vscode", 0755)
	}

	err = os.WriteFile(".vscode/launch.json", output, 0644)
	if err != nil {
		fmt.Println("❌ 写入 launch.json 失败:", err)
		return
	}

	fmt.Println("✅ launch.json 已生成，路径：.vscode/launch.json")
}
