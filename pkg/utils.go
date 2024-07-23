package pkg

import (
	"runtime"
	"strings"
)

func GetPackageName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get caller information")
	}

	// Получаем информацию о функции
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		panic("Could not get function information")
	}

	// Получаем полный путь к файлу
	fullPath := fn.Name()
	// Извлекаем название пакета из пути к файлу
	packageName := extractPackageName(fullPath)
	return packageName
}

func extractPackageName(fullPath string) string {
	parts := strings.Split(fullPath, "/")
	name := strings.Split(parts[len(parts)-1], ".")[0]
	return name
}
