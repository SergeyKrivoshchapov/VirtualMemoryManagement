# Справочник команд и операций

## Содержание
1. [Команды Go](#команды-go)
2. [Команды сборки](#команды-сборки)
3. [API функции](#api-функции)
4. [Типы ошибок](#типы-ошибок)
5. [Конфигурационные переменные](#конфигурационные-переменные)
6. [Файловая структура](#файловая-структура)

---

## Команды Go

### Тестирование

```bash
# Запустить все тесты
go test ./...

# Запустить с подробным выводом
go test -v ./...

# Запустить тесты конкретного пакета
go test ./api
go test ./storage
go test ./cache
go test ./types

# Запустить конкретный тест
go test -v -run TestVMCreate ./api
go test -v -run TestLRU ./cache

# Запустить с таймаутом
go test -timeout 30s ./...

# Запустить с race detector (поиск race conditions)
go test -race ./...

# Генерировать отчет о покрытии
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
go tool cover -func=coverage.out
```

### Сборка и компиляция

```bash
# Собрать DLL для Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o VirtualMemoryManagement.dll .

# Собрать для Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o VirtualMemoryManagement.so .

# Собрать с удалением debug информации
go build -ldflags="-s -w" -o VirtualMemoryManagement.dll .

# Показать информацию о среде
go env
```

### Управление модулями

```bash
# Получить зависимости
go mod download

# Проверить зависимости
go mod verify

# Очистить кэш модулей
go clean -modcache

# Выполнить audit безопасности
go list -json -m all | nancy sleuth
```

### Форматирование и анализ

```bash
# Отформатировать код
go fmt ./...
gofmt -w .

# Проверить код
go vet ./...

# Более строгая проверка (требует golangci-lint)
golangci-lint run ./...
```

### Информация о коде

```bash
# Показать документацию
go doc ./api
go doc ./cache
go doc ./storage

# Запустить встроенный веб-сервер документации
godoc -http=:6060
```

---

## Команды сборки

### Linux сборка (build.sh)

```bash
#!/bin/bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  go build -o VirtualMemoryManagement.dll \
  -ldflags="-s -w" \
  .
```

**Выполнить:**
```bash
chmod +x build.sh
./build.sh
```

### Windows сборка (build.ps1)

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "1"

go build -o VirtualMemoryManagement.dll `
  -ldflags="-s -w" `
  .
```

**Выполнить:**
```powershell
.\build.ps1
```

### Универсальная сборка (build-universal.sh)

```bash
#!/bin/bash
# Windows 64-bit
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o bin/vmm_windows_amd64.dll .

# Windows 32-bit
GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o bin/vmm_windows_386.dll .

# Linux 64-bit
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/vmm_linux_amd64.so .
```

---

## API функции

### Сигнатуры C#

```csharp
// Управление файлами
Result VMCreate(string filename, int size, string type, int stringLength);
Result VMOpen(string filename);
Result VMClose(int handle);

// Работа с данными
Result VMRead(int handle, int index);
Result VMWrite(int handle, int index, string value);

// Управление кэшем
void SetCacheSize(int size);
int GetCacheSize();

// Информация
Result VMStats(int handle);
Result VMHelp(string filename);

// Вспомогательные
int GetHandle();
int[] GetAllHandles();
```

### Структура Result

```csharp
struct Result {
    int Success;           // 1 = успех, 0 = ошибка
    char[] Data;          // Данные результата (256 байт)
    int ErrorCode;        // Код ошибки (если Success = 0)
}
```

### Примеры использования

```csharp
// Создать массив
Result res = VMCreate("data.vm", 1000, "int", 0);
if (res.Success == 1) {
    // Успешно создан
}

// Открыть файл
Result res = VMOpen("data.vm");
if (res.Success == 1) {
    int handle = int.Parse(new string(res.Data).TrimEnd('\0'));
    // Использовать handle
}

// Записать элемент
VMWrite(handle, 0, "42");

// Прочитать элемент
Result readRes = VMRead(handle, 0);
string value = new string(readRes.Data).TrimEnd('\0');

// Закрыть файл
VMClose(handle);
```

---

## Типы ошибок

### Таблица кодов ошибок

```
Код    | Константа                   | Описание
-------|-----------------------------|-----------------------------------------
-1     | ErrCodeFileNotFound         | Файл не найден
-2     | ErrCodeOutOfMemory          | Нехватка памяти
-3     | ErrCodeIndexOutOfRange      | Индекс вне диапазона
-4     | ErrCodeFileOperation        | Ошибка файловой операции
-5     | ErrCodeInvalidType          | Неверный тип данных
-6     | ErrCodeInsufficientDisk     | Недостаточно дискового пространства
-7     | ErrCodeInvalidHandle        | Неверный handle
-8     | ErrCodePageNotFound         | Страница не найдена
-999   | (неизвестная ошибка)        | Неклассифицированная ошибка
```

### Обработка ошибок (C#)

```csharp
Result res = VMOpen("file.vm");

if (res.Success == 0) {
    string errorMsg = new string(res.Data).TrimEnd('\0');
    
    switch (res.ErrorCode) {
        case -1:
            Console.WriteLine($"File not found: {errorMsg}");
            break;
        case -3:
            Console.WriteLine($"Index out of range: {errorMsg}");
            break;
        case -4:
            Console.WriteLine($"File operation error: {errorMsg}");
            // Проверьте:
            // - Файл существует?
            // - Файл уже открыт?
            // - Доступ к файлу?
            break;
        case -7:
            Console.WriteLine($"Invalid handle: {errorMsg}");
            break;
        default:
            Console.WriteLine($"Unknown error {res.ErrorCode}: {errorMsg}");
            break;
    }
}
```

### Диагностика ошибки -4

Ошибка -4 (FileOperation) может означать:

| Сценарий | Решение |
|----------|---------|
| File already exists при Create | Удалить файл перед созданием |
| File already opened | Закрыть первый handle перед открытием |
| File not found при Open | Сначала создать файл (VMCreate) |
| I/O error | Проверить доступ к диску, свободное место |

---

## Конфигурационные переменные

### config/config.go

```go
const (
    BitsPerPage      = 128      // Элементов на странице
    BytesPerBitmap   = 16       // Размер bitmap (128/8)
    PhysicalPageSize = 512      // Размер физической страницы

    MinCacheSize     = 3        // Минимум страниц в кэше
    MaxCacheSize     = 100      // Максимум страниц в кэше
    DefaultCacheSize = 10       // По умолчанию
)

// Вычисляемые значения
func PageDataSize(elemSize int) int
func TotalPageSize(elemSize int) int
```

### Типы данных

```go
type Type byte

const (
    TypeInt     Type = 'I'  // Целые числа (int32)
    TypeChar    Type = 'C'  // Строки фиксированной длины
    TypeVarchar Type = 'V'  // Строки переменной длины
)
```

### Примеры расчетов

```
Для int массива (ElementSize = 4):
  PageDataSize = 128 * 4 = 512 байт
  TotalPageSize = 512 байт (выровнено)
  
Для char(50) (ElementSize = 50):
  PageDataSize = 128 * 50 = 6400 байт
  TotalPageSize = 6912 байт (выровнено)
  
Для varchar (ElementSize = 4):
  PageDataSize = 128 * 4 = 512 байт
  TotalPageSize = 512 байт (выровнено)

Array с 1000 элементов:
  PageCount = ceil(1000 / 128) = 8 страниц
  FileSize = 15 (header) + 8 * 512 = 4111 байт
```

---

## Файловая структура

### Структура .vm файла

```
Offset  Size   Тип           Описание
──────────────────────────────────────────────
0       2      char[2]       Сигнатура "VM"
2       8      int64         Размер массива
10      1      byte          Тип ('I', 'C', 'V')
11      4      int32         Длина строки
15      ...    Page[N]       Страницы данных
```

### Структура страницы

```
Offset  Size   Описание
──────────────────────────────────────
0       16     Bitmap (128 бит)
16      496    Data (элементы)
512 б   ВСЕГО  Размер страницы
```

### Структура .vm.varchar файла

```
Содержимое: последовательность строк переменной длины
Индекс: отображение index → offset в файле
```

---

## Команды для работы с файлами

### Linux

```bash
# Просмотр структуры файла
hexdump -C data.vm | head -20

# Проверить сигнатуру
head -c 2 data.vm

# Размер файла
ls -lh data.vm
stat data.vm

# Удалить файлы
rm data.vm data.vm.varchar

# Поиск всех .vm файлов
find . -name "*.vm" -type f
```

### Windows PowerShell

```powershell
# Просмотр первых 32 байт
Get-Content data.vm -Encoding Byte -TotalCount 32 | Format-Hex

# Размер файла
Get-ChildItem data.vm | Select-Object Length

# Удалить
Remove-Item data.vm
Remove-Item data.vm.varchar

# Поиск файлов
Get-ChildItem -Filter *.vm -Recurse
```

---

## Полезные утилиты

### Go инструменты

```bash
# Анализ производительности
go tool pprof

# Просмотр AST (Abstract Syntax Tree)
go tool ast

# Генерация документации
godoc ./api

# Форматирование
goimports -w .
```

### Внешние инструменты

```bash
# Анализ кода
golangci-lint run ./...

# Проверка безопасности
nancy sleuth

# Бенчмарки
benchstat bench.txt

# Профилирование памяти
pprof --http=:8080 cpu.prof
```

---

## Переменные окружения

### Для сборки DLL

```bash
# Windows
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

# Linux
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
```

### Для разработки

```bash
# Отключить модули (старый стиль)
export GO111MODULE=off

# Использовать локальный GOPATH
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

---

## Скрипты Makefile

### Основные цели

```makefile
test:           # Запустить все тесты
build:          # Собрать DLL
clean:          # Очистить сборки
coverage:       # Отчет о покрытии
fmt:            # Форматировать код
lint:           # Анализ кода
help:           # Справка
```

**Использование:**
```bash
make test
make build
make coverage
```

---

## Примеры одностроковых команд

```bash
# Запустить все тесты и показать результат
go test -v ./... 2>&1 | tee test_results.txt

# Собрать и показать размер
go build -ldflags="-s -w" -o vmm.dll . && ls -lh vmm.dll

# Запустить тест с поиском race conditions
go test -race -count=1 ./...

# Генерировать и открыть отчет о покрытии
go test -coverprofile=c.out ./... && go tool cover -html=c.out

# Форматировать и проверить все файлы
go fmt ./... && go vet ./...

# Скопировать DLL в папку C#
cp VirtualMemoryManagement.dll ../CSApp/bin/Debug/

# Удалить все .vm файлы
find . -name "*.vm*" -type f -delete
```

---

## Быстрая диагностика

### Проверка окружения

```bash
# Версия Go
go version

# Информация о среде
go env | grep -E "GOOS|GOARCH|CGO"

# Проверка CGO
go env CGO_ENABLED

# Установленные компиляторы
which gcc
which cc
which mingw-w64-x86_64-gcc
```

### Проверка проекта

```bash
# Все ли пакеты компилируются
go build ./...

# Все ли тесты проходят
go test -short ./...

# Размер DLL
ls -lh VirtualMemoryManagement.dll

# Зависимости
go mod graph

# Проверка импортов
go build -n ./api 2>&1 | grep import
```

---

## Профилирование

### CPU Profiling

```bash
# Запустить тесты с CPU профилем
go test -cpuprofile=cpu.prof ./api

# Анализировать профиль
go tool pprof cpu.prof
(pprof) top10
(pprof) list functionName
```

### Memory Profiling

```bash
# Запустить тесты с memory профилем
go test -memprofile=mem.prof ./api

# Анализировать профиль
go tool pprof mem.prof
```

### Benchmark

```bash
# Запустить бенчмарки
go test -bench=. -benchmem ./cache

# Сохранить результаты
go test -bench=. -benchmem ./cache > bench.txt

# Сравнить результаты
benchstat bench1.txt bench2.txt
```

---

**Версия:** 1.0
**Дата:** March 18, 2026
**Язык:** Go 1.26+

