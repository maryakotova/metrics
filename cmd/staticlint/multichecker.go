package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/staticcheck"
)

const Config = `config.json`

type ConfigData struct {
	Staticcheck []string
}

func main() {

	mychecks := []*analysis.Analyzer{
		//  стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		nilness.Analyzer,
		unreachable.Analyzer,
		inspect.Analyzer,
	}

	// все анализаторы класса SA пакета staticcheck.io
	mychecks = addSAAnalyzers(mychecks)

	// анализаторы из файла config.json (анализаторы классов S1 и ST1 пакета staticcheck.io)
	mychecks, err := addAnalyzerFromConfig(mychecks)
	if err != nil {
		err = fmt.Errorf("ошибка при чтении настроек из файла: %w, анализаторы будут запущены без использования настроек", err)
		fmt.Println(err.Error())
	}

	multichecker.Main(
		mychecks...,
	)

}

// все анализаторы класса SA пакета staticcheck.io
func addSAAnalyzers(checks []*analysis.Analyzer) []*analysis.Analyzer {
	for _, a := range staticcheck.Analyzers {
		checks = append(checks, a.Analyzer)
	}

	return checks
}

// анализаторы из файла config.json (анализаторы классов S1 и ST1 пакета staticcheck.io)
func addAnalyzerFromConfig(checks []*analysis.Analyzer) ([]*analysis.Analyzer, error) {
	appfile, err := os.Executable()
	if err != nil {
		return checks, err
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		return checks, err
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		return checks, err
	}

	staticChecks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		staticChecks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if staticChecks[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}

	return checks, nil
}
