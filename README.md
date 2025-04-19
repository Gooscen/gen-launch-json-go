# 🧠 VS Code Go Debug Launcher Generator

A lightweight and intelligent Go script that automatically generates `launch.json` configurations for debugging Go applications in Visual Studio Code — supporting multiple entry points and compound debugging with ease.

---

## ✨ Features

- 🔍 Recursively scans all subdirectories
- ✅ Detects directories with `package main` and a `func main()`
- ⚙️ Generates individual debug configurations for each valid main package
- 🧩 Automatically builds a `compound` config for launching all at once
- ⛔ Excludes itself (`gen_launch_json.go`) from being scanned

---

## 📦 Installation

No dependencies required. Just download or copy the script into your project root:

```bash
wget https://Gooscen/gen-launch-json-go/main/gen_launch_json.go
```


## 🚀 Usage
From your project root, run:
```bash
go run gen_launch_json.go
```
It will generate a .vscode/launch.json with debug configurations for all main() entries it finds.


## 🧪 Example Project Structure
```
my-project/
├── cmd/
│   ├── server/      <- contains main()
│   └── client/      <- contains main()
├── pkg/
│   └── utils/       <- no main()
├── gen_launch_json.go
```
The script will generate:
	•	Debug ./cmd/server
	•	Debug ./cmd/client
	•	Debug All Main Packages (as a compound launch)