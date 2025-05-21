package safeexit // import "metrics/cmd/staticlint/safeexit"

В пакете safeexit реализован анализатор кода, который выполняет поиск прямых
вызовов os.Exit в функции main пакета main.

# Использование

Данный мультичекер необходимо использовать с помощью мультичекера:

    	multichecker.Main(
    		safeexit.SafeExitAnalyzer,
     	другие анализаторы,
    	)

var SafeExitAnalyzer = &analysis.Analyzer{ ... }
