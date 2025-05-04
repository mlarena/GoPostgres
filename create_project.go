package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <project_name>")
        return
    }

    projectName := os.Args[1]
    createProjectStructure(projectName)
    fmt.Printf("Project structure created for %s\n", projectName)
}

func createProjectStructure(projectName string) {
    dirs := []string{
        projectName,
        filepath.Join(projectName, "cmd", projectName),
        filepath.Join(projectName, "internal", "handlers"),
        filepath.Join(projectName, "internal", "models"),
        filepath.Join(projectName, "internal", "services"),
        filepath.Join(projectName, "internal", "utils"),
        filepath.Join(projectName, "pkg", "mylib"),
    }

    files := map[string]string{
        filepath.Join(projectName, "cmd", projectName, "main.go"): mainGoContent,
        filepath.Join(projectName, "internal", "handlers", "handlers.go"): handlersGoContent,
        filepath.Join(projectName, "internal", "models", "models.go"): modelsGoContent,
        filepath.Join(projectName, "internal", "services", "services.go"): servicesGoContent,
        filepath.Join(projectName, "internal", "utils", "utils.go"): utilsGoContent,
        filepath.Join(projectName, "pkg", "mylib", "mylib.go"): mylibGoContent,
        filepath.Join(projectName, "go.mod"): goModContent(projectName),
        filepath.Join(projectName, "README.md"): readmeContent,
    }

    for _, dir := range dirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            fmt.Printf("Error creating directory %s: %v\n", dir, err)
        }
    }

    for file, content := range files {
        err := os.WriteFile(file, []byte(content), 0644)
        if err != nil {
            fmt.Printf("Error creating file %s: %v\n", file, err)
        }
    }
}

const mainGoContent = `package main

import (
    "fmt"
    "%s/internal/handlers"
    "%s/internal/services"
)

func main() {
    fmt.Println("Starting %s...")

    // Пример использования внутренних пакетов
    handlers.HandleRequest()
    services.ProcessData()

    fmt.Println("%s finished.")
}
`

const handlersGoContent = `package handlers

import "fmt"

func HandleRequest() {
    fmt.Println("Handling request...")
}
`

const modelsGoContent = `package models

// Define your data models here
`

const servicesGoContent = `package services

import "fmt"

func ProcessData() {
    fmt.Println("Processing data...")
}
`

const utilsGoContent = `package utils

import "fmt"

func HelperFunction() {
    fmt.Println("Helper function called.")
}
`

const mylibGoContent = `package mylib

import "fmt"

func PublicFunction() {
    fmt.Println("Public function called.")
}
`

func goModContent(projectName string) string {
    return fmt.Sprintf(`module %s

go 1.20
`, projectName)
}

const readmeContent = `# %s

This is a Go project.
`
