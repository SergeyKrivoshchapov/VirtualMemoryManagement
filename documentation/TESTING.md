# Testing Documentation - Virtual Memory Management System

## Введение в тестирование

### Типы тестов в проекте

1. **Unit Tests** - тестирование отдельных компонентов
2. **Integration Tests** - тестирование взаимодействия компонентов
3. **API Tests** - тестирование API функций

---

## Структура тестов

```
VirtualMemoryManagement/
├── api/
│   └── vmm_api_test.go           # Тесты API функций
├── cache/
│   └── lru_test.go               # Тесты LRU кэша
├── errors/
│   └── errors_test.go            # Тесты обработки ошибок
├── storage/
│   ├── pagefile_test.go          # Тесты работы со страницами
│   ├── varcharfile_test.go       # Тесты работы со строками
│   └── binary_test.go            # Тесты бинарных операций
├── types/
│   ├── array/
│   │   └── array_test.go         # Тесты типов массивов
│   ├── bitmap/
│   │   └── bitmap_test.go        # Тесты битовых карт
│   ├── page/
│   │   └── page_test.go          # Тесты структуры страницы
│   └── result/
│       └── result_test.go        # Тесты структуры результата
└── tests/
    └── testutils/
        └── helpers.go            # Вспомогательные функции
```

---

## Запуск тестов

### Команды Go

```bash
# Запустить все тесты
go test ./...

# Запустить тесты конкретного пакета
go test ./api
go test ./storage
go test ./cache
go test ./types/...

# Запустить с подробным выводом
go test -v ./...

# Запустить с фильтром по имени теста
go test -v -run TestVMCreate ./api

# Запустить с таймаутом
go test -timeout 30s ./...

# Запустить с информацией о покрытии
go test -cover ./...

# Подробное покрытие
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Написание тестов

### Базовая структура теста

```go
package api

import (
    "testing"
    "os"
)

func TestVMCreate(t *testing.T) {
    // Setup - подготовка
    filename := "test_create.vm"
    defer os.Remove(filename)
    
    // Execute - выполнение
    result := VMCreate(filename, 100, "int", 0)
    
    // Assert - проверка
    if result.Success != 1 {
        t.Fatalf("Expected success, got error: %s", result.GetErrorMessage())
    }
}
```

### Вспомогательные функции (helpers.go)

```go
package testutils

import (
    "testing"
    "os"
    "VirtualMemoryManagement/types/result"
)

// AssertSuccess проверяет успешность результата
func AssertSuccess(t *testing.T, res result.Result) {
    t.Helper()
    if res.Success != 1 {
        t.Fatalf("Expected success, got error: %s", res.GetErrorMessage())
    }
}

// AssertError проверяет наличие ошибки с кодом
func AssertError(t *testing.T, res result.Result, expectedCode int) {
    t.Helper()
    if res.Success != 0 {
        t.Fatalf("Expected error, got success")
    }
    if res.ErrorCode != expectedCode {
        t.Fatalf("Expected error code %d, got %d", expectedCode, res.ErrorCode)
    }
}

// CleanupFile удаляет файл и связанные файлы
func CleanupFile(filename string) {
    os.Remove(filename)
    os.Remove(filename + ".varchar")
}

// CreateTempFile создает временный файл и возвращает его имя
func CreateTempFile(t *testing.T, prefix string) string {
    f, err := os.CreateTemp("", prefix)
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    f.Close()
    return f.Name()
}
```

---

## Примеры тестов

### Test 1: Создание и открытие файла

```go
func TestVMCreateAndOpen(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    // Создать
    createRes := VMCreate(filename, 100, "int", 0)
    AssertSuccess(t, createRes)
    
    // Открыть
    openRes := VMOpen(filename)
    AssertSuccess(t, openRes)
    
    // Получить handle
    handle, err := strconv.Atoi(openRes.String())
    if err != nil {
        t.Fatalf("Failed to parse handle: %v", err)
    }
    
    // Закрыть
    closeRes := VMClose(handle)
    AssertSuccess(t, closeRes)
}
```

### Test 2: Чтение и запись

```go
func TestVMReadWrite(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    // Setup
    VMCreate(filename, 100, "int", 0)
    openRes := VMOpen(filename)
    handle, _ := strconv.Atoi(openRes.String())
    
    // Write
    writeRes := VMWrite(handle, 0, "42")
    AssertSuccess(t, writeRes)
    
    // Read
    readRes := VMRead(handle, 0)
    AssertSuccess(t, readRes)
    
    // Verify
    if readRes.String() != "42" {
        t.Fatalf("Expected 42, got %s", readRes.String())
    }
    
    VMClose(handle)
}
```

### Test 3: Тестирование ошибок

```go
func TestVMOpenNonexistent(t *testing.T) {
    // Попытка открыть несуществующий файл
    res := VMOpen("nonexistent_file.vm")
    
    AssertError(t, res, errors.ErrCodeFileNotFound)
}

func TestVMCreateDuplicate(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    // Создать первый раз - OK
    res1 := VMCreate(filename, 100, "int", 0)
    AssertSuccess(t, res1)
    
    // Создать второй раз - Error
    res2 := VMCreate(filename, 100, "int", 0)
    AssertError(t, res2, errors.ErrCodeFileOperation)
}

func TestVMDoubleOpen(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    VMCreate(filename, 100, "int", 0)
    
    // Открыть первый раз - OK
    res1 := VMOpen(filename)
    AssertSuccess(t, res1)
    h1, _ := strconv.Atoi(res1.String())
    
    // Открыть второй раз - Error (файл уже открыт)
    res2 := VMOpen(filename)
    AssertError(t, res2, errors.ErrCodeFileOperation)
    
    VMClose(h1)
}

func TestVMIndexOutOfRange(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    VMCreate(filename, 100, "int", 0)
    openRes := VMOpen(filename)
    handle, _ := strconv.Atoi(openRes.String())
    
    // Индекс вне диапазона
    res := VMRead(handle, 100)  // size = 100, индекс 0-99
    AssertError(t, res, errors.ErrCodeIndexOutOfRange)
    
    VMClose(handle)
}

func TestVMInvalidHandle(t *testing.T) {
    // Использовать несуществующий handle
    res := VMRead(999, 0)
    AssertError(t, res, errors.ErrCodeInvalidHandle)
}
```

### Test 4: Тестирование кэша

```go
func TestCacheHitRate(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    // Создать с маленьким кэшем
    SetCacheSize(3)
    
    VMCreate(filename, 1000, "int", 0)
    openRes := VMOpen(filename)
    handle, _ := strconv.Atoi(openRes.String())
    
    // Написать
    for i := 0; i < 100; i++ {
        VMWrite(handle, i, strconv.Itoa(i))
    }
    
    // Последовательное чтение - должно быть быстро (hit cache)
    for i := 0; i < 100; i++ {
        res := VMRead(handle, i)
        if !res.IsSuccess() {
            t.Fatalf("Read failed at index %d: %s", i, res.GetErrorMessage())
        }
    }
    
    // Получить статистику
    statsRes := VMStats(handle)
    AssertSuccess(t, statsRes)
    
    VMClose(handle)
    SetCacheSize(10)  // Вернуть к норме
}
```

### Test 5: Работа со строками

```go
func TestVarcharArrays(t *testing.T) {
    filename := "test.vm"
    defer CleanupFile(filename)
    
    // Создать varchar массив
    createRes := VMCreate(filename, 100, "varchar", 100)
    AssertSuccess(t, createRes)
    
    openRes := VMOpen(filename)
    handle, _ := strconv.Atoi(openRes.String())
    
    // Написать строки
    testStrings := []string{"Hello", "World", "Virtual Memory", "Go"}
    for i, s := range testStrings {
        writeRes := VMWrite(handle, i, s)
        AssertSuccess(t, writeRes)
    }
    
    // Прочитать и проверить
    for i, expectedStr := range testStrings {
        readRes := VMRead(handle, i)
        AssertSuccess(t, readRes)
        
        if readRes.String() != expectedStr {
            t.Fatalf("Expected %s, got %s", expectedStr, readRes.String())
        }
    }
    
    VMClose(handle)
}
```

### Test 6: Несколько файлов одновременно

```go
func TestMultipleFiles(t *testing.T) {
    file1 := "test1.vm"
    file2 := "test2.vm"
    defer CleanupFile(file1)
    defer CleanupFile(file2)
    
    // Создать оба файла
    VMCreate(file1, 100, "int", 0)
    VMCreate(file2, 100, "varchar", 50)
    
    // Открыть оба
    res1 := VMOpen(file1)
    res2 := VMOpen(file2)
    AssertSuccess(t, res1)
    AssertSuccess(t, res2)
    
    h1, _ := strconv.Atoi(res1.String())
    h2, _ := strconv.Atoi(res2.String())
    
    // Работать с обоими
    VMWrite(h1, 0, "100")
    VMWrite(h2, 0, "Hello")
    
    // Проверить
    res := VMRead(h1, 0)
    if res.String() != "100" {
        t.Fatalf("h1 read failed")
    }
    
    res = VMRead(h2, 0)
    if res.String() != "Hello" {
        t.Fatalf("h2 read failed")
    }
    
    // Закрыть оба
    VMClose(h1)
    VMClose(h2)
}
```

### Test 7: Сохранение данных на диск

```go
func TestDataPersistence(t *testing.T) {
    filename := "persistence_test.vm"
    defer CleanupFile(filename)
    
    // Создать и написать
    VMCreate(filename, 100, "int", 0)
    res := VMOpen(filename)
    h1, _ := strconv.Atoi(res.String())
    
    for i := 0; i < 50; i++ {
        VMWrite(h1, i, strconv.Itoa(i*10))
    }
    
    VMClose(h1)
    
    // Переоткрыть и проверить
    res = VMOpen(filename)
    h2, _ := strconv.Atoi(res.String())
    
    for i := 0; i < 50; i++ {
        readRes := VMRead(h2, i)
        AssertSuccess(t, readRes)
        
        expected := strconv.Itoa(i * 10)
        if readRes.String() != expected {
            t.Fatalf("Index %d: expected %s, got %s", i, expected, readRes.String())
        }
    }
    
    VMClose(h2)
}
```

---

## Тестирование компонентов

### Тестирование PageFile

```go
package storage

import (
    "testing"
    "os"
    "VirtualMemoryManagement/types/page"
    "VirtualMemoryManagement/types/array"
)

func TestPageFileCreateRead(t *testing.T) {
    filename := "test_page.vm"
    defer os.Remove(filename)
    
    // Create
    pf := NewPageFile(filename)
    if err := pf.Create(128, array.TypeInt, 0); err != nil {
        t.Fatalf("Create failed: %v", err)
    }
    
    // ReadPage
    p, err := pf.ReadPage(0)
    if err != nil {
        t.Fatalf("ReadPage failed: %v", err)
    }
    
    if p == nil {
        t.Fatal("Page is nil")
    }
    
    pf.Close()
}

func TestPageFileWritePage(t *testing.T) {
    filename := "test_page.vm"
    defer os.Remove(filename)
    
    pf := NewPageFile(filename)
    pf.Create(128, array.TypeInt, 0)
    
    // Create a page
    p := &page.Page{
        AbsoluteNumber: 0,
        Bitmap: [16]byte{0xFF, 0xFF, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        Data: make([]byte, 496),
        Dirty: true,
    }
    
    // Write
    if err := pf.WritePage(p); err != nil {
        t.Fatalf("WritePage failed: %v", err)
    }
    
    // Read and verify
    p2, _ := pf.ReadPage(0)
    if p2.AbsoluteNumber != 0 {
        t.Fatal("Page number mismatch")
    }
    
    pf.Close()
}
```

### Тестирование LRU Cache

```go
package cache

import (
    "testing"
    "VirtualMemoryManagement/types/page"
)

func TestLRUGet(t *testing.T) {
    lru := NewLRU(3)
    
    p := &page.Page{AbsoluteNumber: 0}
    lru.Put(p)
    
    retrieved := lru.Get(0)
    if retrieved == nil {
        t.Fatal("Page not found in cache")
    }
}

func TestLRUEviction(t *testing.T) {
    lru := NewLRU(2)
    
    // Add 3 pages (should evict first)
    p1 := &page.Page{AbsoluteNumber: 1}
    p2 := &page.Page{AbsoluteNumber: 2}
    p3 := &page.Page{AbsoluteNumber: 3}
    
    evicted := lru.Put(p1)  // Cache: [p1]
    if evicted != nil {
        t.Fatal("Unexpected eviction")
    }
    
    evicted = lru.Put(p2)   // Cache: [p2, p1]
    if evicted != nil {
        t.Fatal("Unexpected eviction")
    }
    
    evicted = lru.Put(p3)   // Cache: [p3, p2], p1 evicted
    if evicted == nil {
        t.Fatal("Expected eviction")
    }
    if evicted.AbsoluteNumber != 1 {
        t.Fatalf("Expected page 1 to be evicted, got %d", evicted.AbsoluteNumber)
    }
    
    // p1 should not be in cache
    if lru.Get(1) != nil {
        t.Fatal("Page 1 should be evicted")
    }
}

func TestLRUOrder(t *testing.T) {
    lru := NewLRU(3)
    
    // Add pages
    for i := 1; i <= 3; i++ {
        lru.Put(&page.Page{AbsoluteNumber: i})
    }
    
    // Access page 1 (makes it fresh)
    lru.Get(1)
    
    // Add page 4 (should evict page 2, not 1)
    evicted := lru.Put(&page.Page{AbsoluteNumber: 4})
    
    if evicted.AbsoluteNumber != 2 {
        t.Fatalf("Expected page 2 evicted, got %d", evicted.AbsoluteNumber)
    }
}
```

---

## Покрытие тестами

### Генерация отчета

```bash
# Создать файл покрытия
go test -coverprofile=coverage.out ./...

# Открыть в браузере
go tool cover -html=coverage.out

# Вывести в консоль
go tool cover -func=coverage.out
```

### Целевое покрытие

```
Компонент           | Целевое покрытие
────────────────────┼─────────────────
API (vmm_api)       | 95%+
Virtual Array       | 90%+
Storage             | 85%+
Cache               | 95%+
Types               | 95%+
Errors              | 90%+
────────────────────┼─────────────────
Общее               | 85%+
```

---

## Отладка тестов

### Вывод отладочной информации

```bash
# Запустить с verbose флагом
go test -v ./...

# Запустить с конкретным тестом
go test -v -run TestVMCreate ./api

# Запустить с логированием
go test -v -race ./...  # Race detector
```

### Отладка в IDE

**GoLand/IntelliJ IDEA:**
1. Правый клик на функцию теста
2. Select "Run TestName" или "Debug TestName"
3. Используйте breakpoints и пошаговое выполнение

---

## Best Practices

### Хорошие тесты

✅ **Делайте:**
```go
func TestSpecificBehavior(t *testing.T) {
    // 1. Четкое имя теста
    // 2. Setup
    defer cleanup()
    
    // 3. Одна концепция за тест
    result := VMCreate("file.vm", 100, "int", 0)
    
    // 4. Явные проверки
    if result.Success != 1 {
        t.Fatalf("Expected success, got error: %s", result.GetErrorMessage())
    }
}
```

❌ **Избегайте:**
```go
func TestEverything(t *testing.T) {
    // Плохие имена
    // Несколько концепций
    // Неявные проверки
    // Отсутствие cleanup
}
```

### Использование subtests

```go
func TestVMOperations(t *testing.T) {
    tests := []struct {
        name string
        fn   func(t *testing.T)
    }{
        {"Create", testCreate},
        {"Open", testOpen},
        {"Read", testRead},
        {"Write", testWrite},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, tt.fn)
    }
}

func testCreate(t *testing.T) {
    // ...
}
```

### Table-driven tests

```go
func TestVMCreateTypes(t *testing.T) {
    tests := []struct {
        name         string
        size         int
        typ          string
        stringLength int
        shouldFail   bool
    }{
        {"IntArray", 100, "int", 0, false},
        {"CharArray", 100, "char", 20, false},
        {"VarcharArray", 100, "varchar", 100, false},
        {"InvalidType", 100, "invalid", 0, true},
        {"ZeroSize", 0, "int", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            filename := tt.name + ".vm"
            defer CleanupFile(filename)
            
            result := VMCreate(filename, tt.size, tt.typ, tt.stringLength)
            
            if tt.shouldFail {
                if result.Success == 1 {
                    t.Fatalf("Expected failure for %s", tt.name)
                }
            } else {
                if result.Success != 1 {
                    t.Fatalf("Expected success for %s: %s", tt.name, result.GetErrorMessage())
                }
            }
        })
    }
}
```

---

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.26
    
    - name: Run tests
      run: go test -v -cover ./...
    
    - name: Run race detector
      run: go test -race ./...
```

---

**Версия:** 1.0
**Дата:** March 18, 2026

