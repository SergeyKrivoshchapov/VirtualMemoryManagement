# Virtual Memory Management System

> Система управления виртуальной памятью на Go с компиляцией в DLL для C# приложений

## 📋 Содержание

- [Описание](#описание)
- [Характеристики](#характеристики)
- [Архитектура](#архитектура)
- [Быстрый старт](#быстрый-старт)
- [Документация](#документация)
- [Примеры](#примеры)
- [API](#api)
- [Конфигурация](#конфигурация)
- [Тестирование](#тестирование)
- [Сборка](#сборка)
- [FAQ](#faq)

---

## Описание

Virtual Memory Management System (VMMS) - это комплексная система для управления виртуальной памятью, разработанная на Go. Система позволяет работать с большими массивами данных, выходящими за рамки доступной памяти, используя механизм пагинации и LRU кэширования.

**Основное назначение:** Система может быть скомпилирована в DLL и использована в C# GUI приложении для управления большими массивами данных с оптимальной производительностью.

---

## Характеристики

### Основные возможности

✅ **Типы данных:**
- `int` - целые числа (int32)
- `char` - строки фиксированной длины
- `varchar` - строки переменной длины

✅ **Виртуальная память:**
- Пагинация данных на страницы по 512 байт
- Автоматическое управление памятью
- Прозрачное сохранение на диск

✅ **Кэширование:**
- LRU (Least Recently Used) кэш
- Настраиваемый размер кэша (3-100 страниц)
- Оптимизация для последовательного доступа

✅ **Надежность:**
- Потокобезопасные операции
- Автоматическое сохранение данных
- Проверка целостности файлов

✅ **Управление:**
- Открытие/закрытие файлов
- Чтение/запись элементов
- Получение статистики использования

### Характеристики производительности

| Операция | Время | Примечания |
|----------|-------|-----------|
| VMCreate (1000 элемент) | ~10 мс | Создание файла |
| VMOpen | ~5 мс | Открытие существующего |
| VMRead (из кэша) | ~0.1 мс | Попадание в кэш |
| VMRead (с диска) | ~5 мс | Промах кэша |
| VMWrite (в кэш) | ~0.2 мс | Быстрая запись |
| VMClose (без грязных) | ~1 мс | Без сохранения |
| VMClose (с flush) | ~10-50 мс | Сохранение данных |

---

## Архитектура

### Слои системы

```
┌─────────────────────────────────────────┐
│      C# GUI Application (Frontend)      │
└────────────────────┬────────────────────┘
                     │ DLL
┌────────────────────▼────────────────────┐
│         DLL Wrapper (dll_wrapper.go)    │
│      (Экспорт функций для C#)           │
└────────────────────┬────────────────────┘
                     │
┌────────────────────▼────────────────────┐
│       API Layer (api/vmm_api.go)        │
│  VMCreate, VMOpen, VMRead, VMWrite...   │
└────────────────────┬────────────────────┘
                     │
┌────────────────────▼────────────────────┐
│   Virtual Array Layer (virtual_array)   │
│  Управление страницами и кэшем         │
└──┬─────────────────┬──────────────┬──────┘
   │                 │              │
   ▼                 ▼              ▼
┌──────────┐  ┌──────────┐  ┌──────────────┐
│  LRU     │  │PageFile  │  │VarcharFile   │
│  Cache   │  │ Storage  │  │ Storage      │
└──────┬───┘  └────┬─────┘  └────┬─────────┘
       │           │             │
       └───────────┼─────────────┘
                   ▼
           ┌──────────────────┐
           │  Файловая система│
           │  *.vm, *.varchar │
           └──────────────────┘
```

### Компоненты

| Компонент | Назначение | Файлы |
|-----------|-----------|-------|
| **API Layer** | Экспортируемые функции | `api/vmm_api.go` |
| **Virtual Array** | Логика виртуальной памяти | `virtualmemory/virtual_array.go` |
| **LRU Cache** | Кэширование страниц | `cache/lru.go` |
| **Page Storage** | Хранение на диск | `storage/pagefile.go` |
| **Varchar Storage** | Строки переменной длины | `storage/varcharfile.go` |
| **Types** | Типы данных | `types/*/` |
| **Errors** | Обработка ошибок | `errors/errors.go` |

---

## Быстрый старт

### Требования

- **Go 1.26+**
- **CGO** (для компиляции DLL)
- **MinGW** (для cross-compilation на Linux)

### Сборка

```bash
# Клонировать репозиторий
git clone <repo>
cd VirtualMemoryManagement

# Linux
./build.sh

# Windows PowerShell
.\build.ps1
```

### Первая программа (C#)

```csharp
using System;
using System.Runtime.InteropServices;

class Program {
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMCreate(string filename, int size, string type, int stringLength);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMOpen(string filename);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMWrite(int handle, int index, string value);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMRead(int handle, int index);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMClose(int handle);
    
    static void Main() {
        // Создать
        VMCreate("data.vm", 1000, "int", 0);
        
        // Открыть
        Result res = VMOpen("data.vm");
        int handle = int.Parse(new string(res.Data).TrimEnd('\0'));
        
        // Записать
        VMWrite(handle, 0, "42");
        
        // Прочитать
        Result readRes = VMRead(handle, 0);
        Console.WriteLine(new string(readRes.Data).TrimEnd('\0'));  // Output: 42
        
        // Закрыть
        VMClose(handle);
    }
}
```

---

## Документация

### Основные документы

| Документ | Описание |
|----------|---------|
| **[DOCUMENTATION.md](DOCUMENTATION.md)** | Полная документация API и компонентов |
| **[ARCHITECTURE.md](ARCHITECTURE.md)** | Детальное описание архитектуры |
| **[QUICKSTART.md](QUICKSTART.md)** | Быстрый старт и примеры |
| **[TESTING.md](TESTING.md)** | Руководство по тестированию |

### Быстрые ссылки

- [API документация](DOCUMENTATION.md#api-документация)
- [Типы данных](DOCUMENTATION.md#типы-данных)
- [Коды ошибок](DOCUMENTATION.md#коды-ошибок)
- [Примеры использования](QUICKSTART.md#примеры-кода)
- [Конфигурация](DOCUMENTATION.md#конфигурация)

---

## Примеры

### Пример 1: Массив целых чисел

```csharp
// Создать
VMCreate("numbers.vm", 10000, "int", 0);

// Открыть
int handle = int.Parse(new string(VMOpen("numbers.vm").Data).TrimEnd('\0'));

// Заполнить
for (int i = 0; i < 1000; i++) {
    VMWrite(handle, i, i.ToString());
}

// Прочитать
for (int i = 0; i < 1000; i++) {
    string value = new string(VMRead(handle, i).Data).TrimEnd('\0');
    Console.WriteLine($"{i}: {value}");
}

// Закрыть
VMClose(handle);
```

### Пример 2: Строки переменной длины

```csharp
VMCreate("strings.vm", 1000, "varchar", 100);
int handle = int.Parse(new string(VMOpen("strings.vm").Data).TrimEnd('\0'));

// Записать строки
VMWrite(handle, 0, "Hello");
VMWrite(handle, 1, "World");
VMWrite(handle, 2, "Virtual Memory");

// Прочитать
string str0 = new string(VMRead(handle, 0).Data).TrimEnd('\0');  // "Hello"

VMClose(handle);
```

### Пример 3: Управление кэшем

```csharp
// Увеличить кэш для интенсивной работы
SetCacheSize(50);

// ... интенсивная работа ...

// Вернуть к норме
SetCacheSize(10);
```

---

## API

### Основные функции

#### VMCreate

```csharp
Result VMCreate(string filename, int size, string type, int stringLength)
```

Создает новый виртуальный массив.

**Параметры:**
- `filename` - путь к файлу
- `size` - количество элементов
- `type` - тип: "int", "char", "varchar"
- `stringLength` - длина строки (для char/varchar)

**Результат:**
- Success = 1: создан успешно
- Success = 0: ошибка (ErrorCode содержит код)

#### VMOpen

```csharp
Result VMOpen(string filename)
```

Открывает существующий виртуальный массив.

**Результат:**
- Success = 1: Data содержит handle (строка)
- Success = 0: ошибка (файл не существует или уже открыт)

#### VMRead

```csharp
Result VMRead(int handle, int index)
```

Читает значение по индексу.

**Результат:**
- Success = 1: Data содержит значение
- Success = 0: ошибка

#### VMWrite

```csharp
Result VMWrite(int handle, int index, string value)
```

Записывает значение по индексу.

**Результат:**
- Success = 1: записано успешно
- Success = 0: ошибка

#### VMClose

```csharp
Result VMClose(int handle)
```

Закрывает файл и сохраняет данные.

---

## Конфигурация

### Параметры кэша

```go
// config/config.go
const (
    MinCacheSize     = 3    // Минимум
    MaxCacheSize     = 100  // Максимум
    DefaultCacheSize = 10   // По умолчанию
)
```

### Рекомендуемые значения

```
3-5:    Минимальное использование памяти
10:     Баланс (по умолчанию)
20-50:  Интенсивная работа
100:    Максимальная производительность
```

---

## Тестирование

### Запуск тестов

```bash
# Все тесты
go test ./...

# Конкретный пакет
go test ./api
go test ./storage

# С подробным выводом
go test -v ./...

# С покрытием
go test -cover ./...
```

### Структура тестов

```
api/vmm_api_test.go         - Тесты API
storage/pagefile_test.go    - Тесты хранилища
cache/lru_test.go           - Тесты кэша
types/*/\*_test.go          - Тесты типов
```

---

## Сборка

### Требования для сборки

1. **Go 1.26+**
   ```bash
   go version
   ```

2. **CGO**
   ```bash
   # Linux: установить mingw-w64
   sudo apt-get install mingw-w64
   
   # Windows: встроено
   ```

### Команды сборки

**Linux:**
```bash
./build.sh
```

**Windows:**
```powershell
.\build.ps1
```

**Универсальная сборка:**
```bash
./build-universal.sh
```

---

## Структура файлов

```
VirtualMemoryManagement/
├── api/                          # API слой
│   ├── vmm_api.go               # Основной API
│   └── vmm_api_test.go          # Тесты API
├── cache/                        # Кэширование
│   ├── lru.go                   # LRU кэш
│   └── lru_test.go              # Тесты
├── config/                       # Конфигурация
│   └── config.go                # Параметры
├── errors/                       # Обработка ошибок
│   └── errors.go                # Типы ошибок
├── storage/                      # Хранилище
│   ├── pagefile.go              # Страницы
│   ├── varcharfile.go           # Строки
│   ├── header.go                # Заголовок
│   └── binary.go                # Бинарные операции
├── types/                        # Типы данных
│   ├── array/
│   ├── bitmap/
│   ├── page/
│   └── result/
├── virtualmemory/               # Основная логика
│   └── virtual_array.go         # Виртуальный массив
├── dll_wrapper.go               # DLL экспорт
├── main.go                      # Точка входа
├── vmm.h                        # C заголовок
├── go.mod                       # Модули
├── DOCUMENTATION.md             # Полная документация
├── ARCHITECTURE.md              # Архитектура
├── QUICKSTART.md                # Быстрый старт
├── TESTING.md                   # Тестирование
└── README.md                    # Этот файл
```

---

## Коды ошибок

| Код | Ошибка | Решение |
|-----|--------|---------|
| -1 | FileNotFound | Файл не существует, создайте с VMCreate |
| -2 | OutOfMemory | Недостаточно памяти, уменьшите кэш |
| -3 | IndexOutOfRange | Индекс вне массива (0 до size-1) |
| **-4** | **FileOperation** | **Файл уже открыт, существует или I/O ошибка** |
| -5 | InvalidType | Неверный тип данных |
| -6 | InsufficientDisk | Нет места на диске |
| -7 | InvalidHandle | Handle не существует или закрыт |
| -8 | PageNotFound | Ошибка при чтении страницы |

---

## FAQ

### Q: Как работает виртуальная память?

**A:** Данные разделены на страницы (512 байт). LRU кэш хранит часто используемые страницы в памяти, остальные на диске. При обращении к странице вне кэша - читается с диска.

### Q: Что произойдет если открыть один файл дважды?

**A:** Второе открытие вернет ошибку -4 (FileOperation) с сообщением "file already opened".

### Q: Потокобезопасна ли система?

**A:** Да, все операции API защищены мьютексом. Однако обращаться можно только через DLL из разных потоков.

### Q: Как улучшить производительность?

**A:** 
1. Используйте последовательный доступ к элементам
2. Увеличьте размер кэша (SetCacheSize)
3. Минимизируйте Open/Close операции

### Q: Какой максимальный размер массива?

**A:** Теоретически - неограниченный (до 2^31 элементов), практически - размер диска.

### Q: Можно ли менять тип данных существующего массива?

**A:** Нет. Нужно создать новый массив и скопировать данные.

---

## Проблемы и решения

### Ошибка -4 при VMCreate

**Причина:** Файл уже существует

**Решение:**
```csharp
File.Delete("file.vm");
File.Delete("file.vm.varchar");
VMCreate("file.vm", 1000, "int", 0);
```

### Медленная производительность

**Причина:** Маленький кэш или случайный доступ

**Решение:**
```csharp
SetCacheSize(50);  // Увеличить кэш
// Использовать последовательный доступ
for (int i = 0; i < 10000; i++) {
    VMRead(handle, i);  // Хорошо
}
```

### Ошибка -3 при VMRead

**Причина:** Индекс вне диапазона

**Решение:**
```csharp
VMCreate("file.vm", 100, "int", 0);
// Правильный диапазон: 0-99
VMRead(handle, 50);  // OK
VMRead(handle, 100);  // Error -3
```

---

## Лицензия

Проект разработан в рамках лабораторной работы.

---

## Автор

Сергей Кривощапов

**Дата создания:** March 2026

---

## Ссылки

- [Полная документация](DOCUMENTATION.md)
- [Архитектура](ARCHITECTURE.md)
- [Быстрый старт](QUICKSTART.md)
- [Тестирование](TESTING.md)

---

**Последнее обновление:** March 18, 2026
**Версия:** 1.0
**Язык:** Go 1.26+
**Целевая платформа:** Windows/Linux

