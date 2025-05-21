// В пакете safeexit реализован анализатор кода, который выполняет поиск прямых вызовов os.Exit в функции main пакета main.
//
// # Использование
//
// Данный мультичекер необходимо использовать с помощью мультичекера:
//
//		multichecker.Main(
//			safeexit.SafeExitAnalyzer,
//	 	другие анализаторы,
//		)
package safeexit

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var SafeExitAnalyzer = &analysis.Analyzer{
	Name: "safeexit",
	Doc:  "check for using a direct os.Exit call in the main function of the main package",
	Run:  run,
}

// Run  отвечает за анализ исходного кода
func run(pass *analysis.Pass) (any, error) {

	// проверка, что пакет называется main
	if pass.Pkg.Name() != "main" {
		// err := fmt.Errorf("имя пакета не main")
		return nil, nil
	}

	for _, file := range pass.Files {

		ast.Inspect(file, func(n ast.Node) bool {

			// Получение диапозона строк функции main. Если возвращается 0, значит функция main не найдена.
			start, end := getMainFunctionLineRange(file)
			if start == 0 {
				return false
			}

			// Поиск вызовов функций
			if call, ok := n.(*ast.CallExpr); ok {
				// Проверка, является ли вызываемая функция os.Exit
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {

					if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "os" && fun.Sel.Name == "Exit" && fun.Pos() >= start && fun.End() <= end {
						// Сообщение об использовании os.Exit
						pass.Reportf(call.Pos(), "Avoid using os.Exit; prefer returning errors")
					}
				}
			}
			return true
		})
	}
	return nil, nil

}

// getMainFunctionLineRange отвечает за получение диапозона строк функции main
func getMainFunctionLineRange(file *ast.File) (start token.Pos, end token.Pos) {
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			if fd.Name.Name == "main" {
				start = fd.Pos()
				end = fd.End()
				//pass.Reportf(f.Pos(), "Функция main начинается на %v и заканчивается на %v", start, end)
				return false
			}
			return false
		}
		return true
	})
	return start, end
}
