В пакете main реализован мультичекер для анализа кода с заранее выбранными
анализаторами.

# Обзор

Мультичекер состоит из: - стандартных статических анализаторов пакета
golang.org/x/tools/go/analysis/passes; - всех анализаторов класса SA пакета
staticcheck.io; - анализаторов классов S1 и ST1 пакета staticcheck.io; - двух
публичных анализаторов

# Использование

Для запуска мультичекера необходимо выполнить команду:

go run cmd/staticlint/multichecker.go [список необходимых пакетов]

# Конфигурация

Список анализаторов классов S1 и ST1 можно изменять с помощью файла config.json
Если при чтении данных из файла возникла ошибка, мультичекер не прекратит свою
работу, будут использованы остальные анализаторы

Пример данных из файла config.json:

    {
        "staticcheck": [
            "S1000",
            "ST1002"
        ]
    }
