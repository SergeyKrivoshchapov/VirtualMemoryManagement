# Virtual Memory Management System - Полная документация

## Содержание
1. [Введение](#введение)
2. [Архитектура проекта](#архитектура-проекта)
3. [Структура файлов](#структура-файлов)
4. [Основные компоненты](#основные-компоненты)
5. [API документация](#api-документация)
6. [Типы данных](#типы-данных)
7. [Обработка ошибок](#обработка-ошибок)
8. [Кэширование](#кэширование)
9. [Работа с файлами](#работа-с-файлами)
10. [Конфигурация](#конфигурация)
11. [Примеры использования](#примеры-использования)
12. [Тестирование](#тестирование)
13. [Сборка проекта](#сборка-проекта)
14. [Коды ошибок](#коды-ошибок)

---

## Введение

**Virtual Memory Management System** - это система управления виртуальной памятью, написанная на Go и предназначенная для компиляции в DLL для использования в C# GUI приложении.

### Основные возможности:
- Создание и управление виртуальными массивами трех типов: int, char, varchar
- Механизм виртуальной памяти с разбиением на страницы
- LRU кэширование для оптимизации доступа к данным
- Сохранение состояния на диск
- Обработка множественных файлов одновременно
- Предотвращение одновременного открытия одного файла

### Назначение:
Система позволяет работать с большими массивами данных, выходящими за рамки доступной памяти, сохраняя часть данных на диск и используя кэш для быстрого доступа к часто используемым элементам.

---

## Архитектура проекта

```
┌─────────────────────────────────────────────────────┐
│         C# GUI Application (Frontend)               │
└─────────────────────────────────────────────────────┘
                         ↓ (DLL)
┌─────────────────────────────────────────────────────┐
│             DLL Wrapper (dll_wrapper.go)             │
│         (Экспорт функций для C#)                    │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│      API Layer (api/vmm_api.go)                     │
│  VMCreate, VMOpen, VMClose, VMRead, VMWrite, Stats  │
└─────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────┐
│        Virtual Array Layer                          │
│    (virtualmemory/virtual_array.go)                 │
│  - Управление страницами                           │
│  - Взаимодействие с кэшем                          │
│  - Обработка чтения/записи                         │
└──────────────────────────────────────────────────────┘
         ↙              ↓              ↘
    ┌──────────┐  ┌──────────┐  ┌──────────────┐
    │Cache     │  │Page      │  │Varchar       │
    │(LRU)     │  │Storage   │  │Storage       │
    └──────────┘  └──────────┘  └──────────────┘
         ↓              ↓              ↓
    ┌─────────────────────────────────────────────┐
    │    Storage Layer                            │
    │  - pagefile.go (хранение страниц)          │
    │  - varcharfile.go (хранение строк)         │
    │  - binary.go (бинарные операции)           │
    └─────────────────────────────────────────────┘
         ↓
    ┌─────────────────────────────────────────────┐
    │      Файловая система                       │
    │  - *.vm (файл виртуальной памяти)          │
    │  - *.vm.varchar (хранилище строк)          │
    └─────────────────────────────────────────────┘
```

---

## Структура файлов

```
VirtualMemoryManagement/
├── api/                          # Слой API
│   ├── vmm_api.go               # Основной API экспортируемые функции
│   └── vmm_api_test.go          # Тесты API
│
├── cache/                        # Система кэширования
│   ├── cacheinterface.go        # Интерфейс кэша
│   ├── lru.go                   # LRU кэш реализация
│   └── lru_test.go              # Тесты LRU кэша
│
├── config/                       # Конфигурация
│   └── config.go                # Константы конфигурации
│
├── errors/                       # Обработка ошибок
│   ├── errors.go                # Определение типов ошибок
│   └── errors_test.go           # Тесты ошибок
│
├── storage/                      # Слой хранилища
│   ├── interfaces.go            # Интерфейсы хранилища
│   ├── pagefile.go              # Работа со страницами
│   ├── pagefile_test.go         # Тесты страниц
│   ├── header.go                # Структура заголовка файла
│   ├── binary.go                # Бинарные операции
│   ├── binary_test.go           # Тесты бинарных операций
│   ├── varcharfile.go           # Работа со строками переменной длины
│   └── varcharfile_test.go      # Тесты работы со строками
│
├── types/                        # Типы данных
│   ├── array/
│   │   ├── array.go             # Тип и информация о массиве
│   │   └── array_test.go        # Тесты
│   ├── page/
│   │   ├── page.go              # Структура страницы
│   │   └── page_test.go         # Тесты
│   ├── bitmap/
│   │   ├── bitmap.go            # Битовая карта страницы
│   │   └── bitmap_test.go       # Тесты
│   └── result/
│       └── result.go            # Структура результата API
│
├── virtualmemory/               # Основная логика
│   └── virtual_array.go         # Реализация виртуального массива
│
├── CLI/                         # C# приложение (не часть Go проекта)
│   └── TimpLaba2_VirtualMemory/
│
├── tests/                       # Вспомогательные тесты
│   └── testutils/
│       └── helpers.go           # Вспомогательные функции для тестов
│
├── main.go                      # Точка входа (для DLL)
├── dll_wrapper.go               # Экспорт функций для DLL
├── cgo_types.go                 # Типы для CGO
├── vmm.h                        # Заголовок для C#
├── go.mod                       # Модуль Go
├── Makefile                     # Сборка проекта
├── build.sh                     # Linux сборка
├── build.ps1                    # Windows PowerShell сборка
├── build-universal.sh           # Универсальная сборка
└── DOCUMENTATION.md             # Эта документация
```

---

## Основные компоненты

### 1. API Layer (`api/vmm_api.go`)

Основной слой взаимодействия с системой. Экспортирует следующие функции:

#### Управление файлами
- `VMCreate(filename string, size int, typ string, stringLength int) Result` - Создает новый виртуальный массив
- `VMOpen(filename string) Result` - Открывает существующий виртуальный массив
- `VMClose(handle int) Result` - Закрывает виртуальный массив и сохраняет данные

#### Операции с данными
- `VMRead(handle int, index int) Result` - Чтение элемента по индексу
- `VMWrite(handle int, index int, value string) Result` - Запись элемента по индексу

#### Статистика и управление
- `VMStats(handle int) Result` - Получить статистику использования
- `SetCacheSize(size int)` - Установить размер кэша
- `GetCacheSize() int` - Получить текущий размер кэша
- `GetHandle() int` - Получить первый доступный handle
- `GetAllHandles() []int` - Получить все активные handles

#### Управление кэшем
- `VMHelp(filename string) Result` - Получить справку по командам

**Особенности:**
- Все функции защищены мьютексом для потокобезопасности
- Поддерживается одновременное открытие нескольких файлов
- Предотвращается открытие одного файла дважды одновременно
- Возвращаемый результат содержит успех/ошибку и данные

### 2. Virtual Array (`virtualmemory/virtual_array.go`)

Основной компонент, реализующий логику виртуальной памяти.

**Ключевые методы:**
- `Create/CreateWithCacheSize()` - Создание нового виртуального массива
- `Open/OpenWithCacheSize()` - Открытие существующего
- `Read(index int) (interface{}, error)` - Чтение элемента
- `Write(index int, value interface{}) error` - Запись элемента
- `Close() error` - Закрытие файла
- `FlushDirtyPages() error` - Сохранение грязных страниц
- `GetStats() string` - Получить статистику

**Работает с:**
- PageStorage - для хранения страниц
- VarcharStorage - для хранения строк переменной длины
- Cache - для кэширования часто используемых страниц

### 3. Кэш (`cache/lru.go`)

LRU (Least Recently Used) кэш реализует стратегию вытеснения.

**Методы:**
- `Get(pageNumber int) *Page` - Получить страницу из кэша
- `Put(p *Page) *Page` - Добавить страницу в кэш (возвращает вытесненную)
- `Contains(pageNumber int) bool` - Проверить наличие в кэше
- `Size() int` - Текущий размер кэша
- `All() []*Page` - Получить все страницы в кэше

**Стратегия:**
- Когда кэш переполнен, вытесняется наименее недавно использованная страница
- Вытесненная страница должна быть сохранена на диск, если она помечена как грязная

### 4. Хранилище (`storage/`)

#### PageFile
Управляет сохранением страниц на диск.

**Структура файла:**
```
┌──────────┬────────────────────────────────────┬─────────────┬─────────────┬───────────────┐
│Signature │ Header (Size, Type, StringLength)  │  Page 0     │  Page 1     │   Page N      │
│  "VM"    │ 8 bytes + type + stringLength      │  512 bytes  │  512 bytes  │   512 bytes   │
└──────────┴────────────────────────────────────┴─────────────┴─────────────┴───────────────┘
```

**Методы:**
- `Create(size, typ, stringLength)` - Создать новый файл
- `Open(filename)` - Открыть существующий файл
- `ReadPage(pageNum) (*Page, error)` - Читать страницу
- `WritePage(page *Page) error` - Писать страницу
- `Close() error` - Закрыть файл

#### VarcharFile
Хранилище для строк переменной длины.

**Структура:**
- Отдельный файл (имя.vm.varchar)
- Хранит строки переменной длины
- Для int и char массивов не используется

#### BinaryIO
Вспомогательные функции для бинарного ввода-вывода.

### 5. Типы данных (`types/`)

#### Array Type
```go
type Type byte
const (
    TypeInt     Type = 'I'  // Целые числа (4 байта)
    TypeChar    Type = 'C'  // Строки фиксированной длины
    TypeVarchar Type = 'V'  // Строки переменной длины
)
```

#### Page
Единица виртуальной памяти.
- AbsoluteNumber - номер страницы в массиве
- Bitmap - битовая карта (какие элементы заполнены)
- Data - данные элементов
- Dirty - флаг изменения

#### Bitmap
Битовая карта отслеживает, какие элементы в странице заполнены.
- 128 бит на страницу = 128 элементов максимум
- По 1 биту на элемент

#### Result
Стандартный результат всех операций API:
```go
type Result struct {
    Success   int         // 1 - успех, 0 - ошибка
    Data      [256]byte   // Данные результата
    ErrorCode int         // Код ошибки
}
```

---

## API документация

### VMCreate

**Создает новый виртуальный массив на диске**

```go
func VMCreate(filename string, size int, typ string, stringLength int) Result
```

**Параметры:**
- `filename` (string) - путь к файлу для создания
- `size` (int) - количество элементов в массиве
- `typ` (string) - тип массива: "int"/"I", "char"/"C", "varchar"/"V"
- `stringLength` (int) - длина строки (для char, обязательно; для varchar, макс длина)

**Возвращает:**
- Success = 1 и Data = "Created" при успехе
- Success = 0 и ErrorCode при ошибке

**Код ошибки 49 (ErrCodeFileOperation):**
- Возвращается если файл уже существует

**Пример:**
```
VMCreate("data.vm", 1000, "int", 0)  // Создать массив целых чисел
VMCreate("names.vm", 100, "varchar", 50)  // Создать массив строк переменной длины
```

---

### VMOpen

**Открывает существующий виртуальный массив**

```go
func VMOpen(filename string) Result
```

**Параметры:**
- `filename` (string) - путь к файлу для открытия

**Возвращает:**
- Success = 1 и Data = "handle_id" (строка с числом) при успехе
- Success = 0 и ErrorCode при ошибке

**Особенности:**
- Один файл не может быть открыт дважды одновременно
- Возвращаемый handle используется для всех последующих операций
- Файл загружается в памяти с использованием кэша

**Код ошибки 49:**
- Файл не существует
- Файл уже открыт другим handle'ом

**Пример:**
```
result = VMOpen("data.vm")
// result.Data содержит строку "1" (handle)
```

---

### VMClose

**Закрывает виртуальный массив и сохраняет все данные**

```go
func VMClose(handle int) Result
```

**Параметры:**
- `handle` (int) - идентификатор открытого файла

**Возвращает:**
- Success = 1 и Data = "Closed" при успехе
- Success = 0 и ErrorCode при ошибке

**Действия:**
1. Проверка валидности handle'а
2. Сохранение всех грязных страниц на диск (FlushDirtyPages)
3. Закрытие файла
4. Удаление handle'а из таблицы

**Пример:**
```
VMClose(1)
```

---

### VMRead

**Читает значение элемента по индексу**

```go
func VMRead(handle int, index int) Result
```

**Параметры:**
- `handle` (int) - идентификатор открытого файла
- `index` (int) - индекс элемента в массиве

**Возвращает:**
- Success = 1 и Data = значение элемента при успехе
- Success = 0 и ErrorCode при ошибке

**Процесс:**
1. Определить номер страницы: pageNum = index / elementsPerPage
2. Проверить кэш (Get из LRU)
3. Если не в кэше:
   - Читать со диска (ReadPage)
   - Добавить в кэш (Put)
4. Получить элемент из страницы по смещению

**Коды ошибок:**
- -7 (ErrCodeInvalidHandle) - неверный handle
- -3 (ErrCodeIndexOutOfRange) - индекс вне массива

**Пример:**
```
VMRead(1, 10)  // Прочитать 11-й элемент (индекс 0-based)
```

---

### VMWrite

**Записывает значение элемента по индексу**

```go
func VMWrite(handle int, index int, value string) Result
```

**Параметры:**
- `handle` (int) - идентификатор открытого файла
- `index` (int) - индекс элемента в массиве
- `value` (string) - значение для записи

**Возвращает:**
- Success = 1 и Data = "Written" при успехе
- Success = 0 и ErrorCode при ошибке

**Процесс:**
1. Определить номер страницы
2. Получить/загрузить страницу из кэша (аналогично Read)
3. Преобразовать значение согласно типу:
   - int: ParseInt(value) → int32
   - char/varchar: строка как есть
4. Записать в страницу (отмечает страницу как Dirty)
5. При нехватке места вытесняется старая страница

**Коды ошибок:**
- -7 (ErrCodeInvalidHandle) - неверный handle
- -3 (ErrCodeIndexOutOfRange) - индекс вне массива
- -5 (ErrCodeInvalidType) - неверный тип данных

**Пример:**
```
VMWrite(1, 10, "42")        // Записать целое число
VMWrite(2, 5, "Hello")      // Записать строку
```

---

### VMStats

**Получает статистику использования виртуального массива**

```go
func VMStats(handle int) Result
```

**Параметры:**
- `handle` (int) - идентификатор открытого файла

**Возвращает:**
- Success = 1 и Data со статистикой при успехе

**Информация в статистике:**
- Размер массива
- Тип данных
- Количество страниц
- Размер кэша
- Статистика обращений
- Статистика попаданий/промахов в кэш

**Пример:**
```
VMStats(1)  // Получить полную статистику
```

---

### SetCacheSize / GetCacheSize

**Управление размером кэша**

```go
func SetCacheSize(size int)
func GetCacheSize() int
```

**Ограничения:**
- Минимум: 3 страницы
- Максимум: 100 страниц
- По умолчанию: 10 страниц

**Пример:**
```
SetCacheSize(20)  // Установить размер кэша на 20 страниц
size := GetCacheSize()
```

---

## Типы данных

### Тип массива (Array Type)

| Тип | Код | Размер элемента | Описание |
|-----|-----|-----------------|---------|
| int | 'I' | 4 байта | Целые числа (int32) |
| char | 'C' | stringLength | Строки фиксированной длины |
| varchar | 'V' | 4 байта (указатель) | Строки переменной длины |

### Структура страницы

```
Page (512 байт)
├── Bitmap (16 байт = 128 бит)        // Какие элементы заполнены
│   └── 1 бит на элемент
├── Data (496 байт)                   // Данные элементов
│   └── Каждый элемент по ElementSize байт
└── Metadata
    ├── AbsoluteNumber (номер в массиве)
    ├── Dirty (флаг изменения)
    └── AccessCount (количество обращений)
```

**Bitmap:**
- 128 бит = 128 элементов максимум на странице
- 1 бит = 1 элемент (0 = не заполнен, 1 = заполнен)

**Dirty флаг:**
- Устанавливается при записи
- Используется при вытеснении из кэша для определения необходимости сохранения

**AccessCount:**
- Отслеживает количество обращений к странице
- Используется для статистики и оптимизации

### Заголовок файла

```
Header (16 байт)
├── Size (int64, 8 байт)              // Количество элементов в массиве
├── Type (byte, 1 байт)               // Тип массива ('I', 'C', 'V')
└── StringLength (int32, 4 байт)      // Длина строки (для char/varchar)
```

---

## Обработка ошибок

### Система кодов ошибок

| Код | Константа | Описание | Причины |
|-----|-----------|---------|---------|
| -1 | ErrCodeFileNotFound | Файл не найден | Файл не существует, неверный путь |
| -2 | ErrCodeOutOfMemory | Нехватка памяти | Не достаточно оперативной памяти для кэша |
| -3 | ErrCodeIndexOutOfRange | Индекс вне диапазона | index < 0 или index >= size |
| -4 | ErrCodeFileOperation | Ошибка файловой операции | Ошибка чтения/записи, файл уже открыт, файл существует при Create |
| -5 | ErrCodeInvalidType | Неверный тип данных | Неверный тип при Create, неверное значение при Write |
| -6 | ErrCodeInsufficientDisk | Недостаточно дискового пространства | Нет места на диске для создания/записи |
| -7 | ErrCodeInvalidHandle | Неверный handle | Handle не существует, был закрыт |
| -8 | ErrCodePageNotFound | Страница не найдена | Ошибка при чтении страницы |

### Тип VMMError

```go
type VMMError struct {
    Code    int      // Код ошибки
    Message string   // Сообщение об ошибке
    Err     error    // Исходная ошибка (если есть)
}
```

### Обработка ошибок в коде

```go
// Создание ошибки с кодом
NewError(code int, message string) *VMMError

// Создание ошибки с оборачиванием другой ошибки
NewErrorWithWrapped(code int, message string, err error) *VMMError

// Получение кода из ошибки
GetErrorCode(err error) int
```

### Результат API

```go
type Result struct {
    Success   int       // 1 = успех, 0 = ошибка
    Data      [256]byte // Данные результата (строка)
    ErrorCode int       // Код ошибки (если Success = 0)
}

// Получить строку из Data
func (r *Result) String() string

// Проверить успешность
func (r *Result) IsSuccess() bool

// Получить сообщение об ошибке
func (r *Result) GetErrorMessage() string
```

---

## Кэширование

### Стратегия LRU (Least Recently Used)

**Принцип:**
- Кэш содержит N страниц (по умолчанию 10)
- При добавлении новой страницы:
  - Если есть место: добавить в начало (как самую свежую)
  - Если нет места: удалить последнюю (самую старую)

**Структура LRU:**
```
Head ↔ Самая новая ↔ ... ↔ Старая ↔ Tail
      (Front)                    (Back)

Новые обращения → перемещают страницу в начало
Вытеснение → удаляется со конца
```

### Процесс кэширования

#### При чтении (VMRead):
1. Проверить: есть ли страница в кэше?
2. Если ДА:
   - Переместить в начало (moveToFront)
   - Вернуть данные
3. Если НЕТ:
   - Читать со диска (ReadPage)
   - Добавить в кэш (Put)
   - Если кэш переполнен: старая страница вытесняется
   - Если старая страница Dirty: сохранить на диск перед вытеснением
   - Вернуть данные из новой страницы

#### При записи (VMWrite):
1. Получить страницу (как при чтении)
2. Записать значение в страницу
3. Отметить страницу как Dirty
4. Переместить в начало кэша (она свежая)
5. При вытеснении: Dirty страница автоматически сохраняется

### Dirty флаг

**Назначение:**
- Отслеживать какие страницы были изменены
- Оптимизировать сохранение: не писать на диск неизменные страницы

**Жизненный цикл:**
1. При создании: Dirty = false
2. При записи: Dirty = true
3. При вытеснении: если Dirty → сохранить на диск
4. После сохранения: Dirty = false
5. При Close: FlushDirtyPages сохраняет все оставшиеся грязные страницы

### Метрики кэша

**Отслеживаемые метрики:**
- Размер кэша (текущее количество страниц)
- Попадания в кэш (cache hits)
- Промахи в кэш (cache misses)
- Hit rate = hits / (hits + misses) * 100%

**Использование в GetStats:**
```
Cache Size: 5/10
Cache Hits: 45
Cache Misses: 15
Hit Rate: 75.00%
```

---

## Работа с файлами

### Структура файла виртуальной памяти (.vm)

```
Блок 0: Заголовок и сигнатура
├── Сигнатура (2 байта): "VM"
├── Size (8 байт): int64
├── Type (1 байт): 'I', 'C' или 'V'
└── StringLength (4 байт): int32

Блок 1+: Страницы данных
├── Page 0 (512 байт)
├── Page 1 (512 байт)
├── ...
└── Page N (512 байт)
```

### Структура файла строк (.vm.varchar)

Используется только для varchar массивов:
- Хранит строки переменной длины
- Индексируется по записям
- Не имеет фиксированной структуры

### Операции с файлами

#### Create
```
1. Открыть файл с флагами: CREATE | WRITE | TRUNCATE
2. Записать сигнатуру "VM"
3. Записать заголовок (Size, Type, StringLength)
4. Создать пустые страницы (нули) для каждого Page Count
5. Закрыть файл
```

#### Open
```
1. Открыть файл с флагами: READ | WRITE
2. Прочитать сигнатуру (проверка формата)
3. Прочитать заголовок (информация о массиве)
4. Загрузить метаинформацию
5. Файл остается открытым для дальнейших операций
```

#### ReadPage
```
1. Вычислить смещение в файле:
   offset = HeaderSize + (pageNum * PageSize)
2. Перейти на это смещение (Seek)
3. Прочитать PageSize байт
4. Разбить на Bitmap и Data
5. Вернуть структуру Page
```

#### WritePage
```
1. Вычислить смещение в файле (как в ReadPage)
2. Перейти на это смещение
3. Записать Bitmap
4. Записать Data
5. Синхронизировать с диском (Sync)
```

#### Close
```
1. Закрыть основной файл
2. Если есть varchar файл: закрыть и его
3. Очистить все структуры в памяти
4. Удалить из таблицы handles
```

### Обработка ошибок файловых операций

```go
// Ошибка создания
if file exists: return ErrCodeFileOperation

// Ошибка открытия
if file not exists: return ErrCodeFileNotFound
if cannot open: return ErrCodeFileOperation

// Ошибка чтения
if read fails: return ErrCodeFileOperation
if short read: return ErrCodeFileOperation

// Ошибка записи
if write fails: return ErrCodeFileOperation
if no disk space: return ErrCodeInsufficientDisk
```

---

## Конфигурация

### Файл config/config.go

```go
const (
    BitsPerPage      = 128          // Максимум элементов на странице
    BytesPerBitmap   = 16           // Размер bitmap (128 / 8)
    PhysicalPageSize = 512          // Размер физической страницы на диске

    MinCacheSize     = 3            // Минимум страниц в кэше
    MaxCacheSize     = 100          // Максимум страниц в кэше
    DefaultCacheSize = 10           // По умолчанию
)
```

### Константы

| Константа | Значение | Назначение |
|-----------|----------|-----------|
| BitsPerPage | 128 | Максимальное количество элементов в одной странице |
| BytesPerBitmap | 16 | Размер битовой карты (128 бит = 16 байт) |
| PhysicalPageSize | 512 | Размер страницы на диске |
| MinCacheSize | 3 | Минимальное количество страниц в кэше |
| MaxCacheSize | 100 | Максимальное количество страниц в кэше |
| DefaultCacheSize | 10 | Размер кэша по умолчанию |

### Вычисляемые значения

```go
// Размер данных в странице (без bitmap)
PageDataSize(elemSize) = BitsPerPage * elemSize

// Общий размер страницы (bitmap + данные, выравнено)
TotalPageSize(elemSize) = ceil((BytesPerBitmap + PageDataSize) / PhysicalPageSize) * PhysicalPageSize

// Количество страниц для массива
PageCount = ceil(ArraySize / BitsPerPage)
```

### Примеры расчета

```
Для int массива (elemSize = 4):
  PageDataSize = 128 * 4 = 512 байт
  TotalPageSize = (16 + 512 + 512 - 1) / 512 * 512 = 1024 байта

Для char(50) (elemSize = 50):
  PageDataSize = 128 * 50 = 6400 байт
  TotalPageSize = (16 + 6400 + 511) / 512 * 512 = 6912 байт

Для varchar (elemSize = 4):
  PageDataSize = 128 * 4 = 512 байт
  TotalPageSize = 1024 байт
```

---

## Примеры использования

### Пример 1: Работа с массивом целых чисел

```cpp
// C# код
using VirtualMemoryManagement;

// Создать массив из 10000 целых чисел
Result createResult = VMCreate("numbers.vm", 10000, "int", 0);
if (createResult.Success == 1) {
    Console.WriteLine("Created successfully");
}

// Открыть файл
Result openResult = VMOpen("numbers.vm");
int handle = int.Parse(openResult.Data.ToString());

// Записать значения
for (int i = 0; i < 100; i++) {
    VMWrite(handle, i, i.ToString());
}

// Прочитать значения
for (int i = 0; i < 100; i++) {
    Result readResult = VMRead(handle, i);
    if (readResult.Success == 1) {
        int value = int.Parse(readResult.Data.ToString());
        Console.WriteLine($"Element {i}: {value}");
    }
}

// Получить статистику
Result statsResult = VMStats(handle);
Console.WriteLine(statsResult.Data.ToString());

// Закрыть файл
VMClose(handle);
```

### Пример 2: Работа со строками переменной длины

```cpp
// Создать массив строк (макс 100 символов)
VMCreate("strings.vm", 1000, "varchar", 100);

// Открыть
Result openResult = VMOpen("strings.vm");
int handle = int.Parse(openResult.Data.ToString());

// Записать строки
VMWrite(handle, 0, "Hello");
VMWrite(handle, 1, "World");
VMWrite(handle, 2, "Virtual Memory");

// Прочитать
Result read0 = VMRead(handle, 0);  // "Hello"
Result read1 = VMRead(handle, 1);  // "World"

// Закрыть
VMClose(handle);
```

### Пример 3: Управление кэшем

```cpp
// Установить больший кэш для интенсивной работы
SetCacheSize(50);  // 50 страниц

VMCreate("large.vm", 100000, "int", 0);
int handle = VMOpen("large.vm").Handle;

// Кэш будет более эффективным
for (int i = 0; i < 100000; i += 128) {  // 128 элементов в странице
    VMWrite(handle, i, "42");
}

// Проверить статистику кэша
Result stats = VMStats(handle);
Console.WriteLine(stats.Data.ToString());

VMClose(handle);
```

### Пример 4: Обработка ошибок

```cpp
// Попытка открыть несуществующий файл
Result openResult = VMOpen("nonexistent.vm");
if (openResult.Success == 0) {
    Console.WriteLine($"Error: {openResult.ErrorCode}");
    Console.WriteLine(openResult.Data.ToString());  // "File not found"
}

// Попытка создать файл который уже существует
VMCreate("test.vm", 100, "int", 0);
Result createResult = VMCreate("test.vm", 100, "int", 0);
// Error 49 (ErrCodeFileOperation): "file already exists"

// Попытка открыть уже открытый файл
int handle1 = VMOpen("test.vm").Handle;
Result result = VMOpen("test.vm");
// Error 49: "file already opened"

VMClose(handle1);
```

### Пример 5: Одновременная работа с несколькими файлами

```cpp
// Создать и открыть несколько файлов
VMCreate("file1.vm", 1000, "int", 0);
VMCreate("file2.vm", 1000, "varchar", 50);
VMCreate("file3.vm", 500, "char", 20);

int h1 = VMOpen("file1.vm").Handle;  // h1 = 1
int h2 = VMOpen("file2.vm").Handle;  // h2 = 2
int h3 = VMOpen("file3.vm").Handle;  // h3 = 3

// Работать с каждым файлом независимо
VMWrite(h1, 0, "100");      // int в file1
VMWrite(h2, 0, "Hello");    // varchar в file2
VMWrite(h3, 0, "Name");     // char в file3

// Каждый handle имеет свой кэш и контекст
VMRead(h1, 0);  // "100"
VMRead(h2, 0);  // "Hello"
VMRead(h3, 0);  // "Name"

// Закрыть все
VMClose(h1);
VMClose(h2);
VMClose(h3);
```

---

## Тестирование

### Структура тестов

Тесты расположены в:
- `api/vmm_api_test.go` - Тесты API
- `storage/pagefile_test.go` - Тесты работы со страницами
- `storage/varcharfile_test.go` - Тесты работы со строками
- `cache/lru_test.go` - Тесты кэша
- `types/*/` - Тесты типов данных

### Запуск тестов

```bash
# Все тесты
go test ./...

# Тесты конкретного пакета
go test ./api
go test ./storage
go test ./cache

# С подробным выводом
go test -v ./...

# С расчетом покрытия
go test -cover ./...

# С профилированием
go test -cpuprofile=cpu.prof ./...
```

### Примеры тестов

#### Базовый тест создания

```go
func TestVMCreate(t *testing.T) {
    defer os.Remove("test_create.vm")
    
    result := VMCreate("test_create.vm", 100, "int", 0)
    if result.Success != 1 {
        t.Fatal("Create failed")
    }
}
```

#### Тест чтения-записи

```go
func TestVMReadWrite(t *testing.T) {
    defer os.Remove("test.vm")
    
    VMCreate("test.vm", 100, "int", 0)
    handle := VMOpen("test.vm").Handle
    
    // Записать
    VMWrite(handle, 0, "42")
    
    // Прочитать
    result := VMRead(handle, 0)
    if result.Data.ToString() != "42" {
        t.Fatal("Write/Read failed")
    }
    
    VMClose(handle)
}
```

#### Тест предотвращения двойного открытия

```go
func TestDoubleOpen(t *testing.T) {
    defer os.Remove("test.vm")
    
    VMCreate("test.vm", 100, "int", 0)
    h1 := VMOpen("test.vm")
    h2 := VMOpen("test.vm")
    
    if h2.Success == 1 {
        t.Fatal("Second open should fail")
    }
    if h2.ErrorCode != errors.ErrCodeFileOperation {
        t.Fatal("Wrong error code")
    }
    
    VMClose(int(h1.Data[0]))
}
```

### Тестовые утилиты

```go
// helpers.go - вспомогательные функции

func CreateTempFile(t *testing.T, name string) string {
    // Создать временный файл
}

func CleanupFile(filename string) {
    // Удалить файл
}

func AssertSuccess(t *testing.T, result Result) {
    if result.Success != 1 {
        t.Fatalf("Expected success, got error: %s", result.GetErrorMessage())
    }
}

func AssertError(t *testing.T, result Result, expectedCode int) {
    if result.Success == 1 {
        t.Fatal("Expected error")
    }
    if result.ErrorCode != expectedCode {
        t.Fatalf("Expected code %d, got %d", expectedCode, result.ErrorCode)
    }
}
```

---

## Сборка проекта

### Linux сборка (build.sh)

```bash
#!/bin/bash
# Скрипт для сборки DLL на Linux

GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  go build -o VirtualMemoryManagement.dll \
  -ldflags="-s -w" \
  .
```

### Windows PowerShell сборка (build.ps1)

```powershell
# Windows сборка с использованием PowerShell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "1"

go build -o VirtualMemoryManagement.dll `
  -ldflags="-s -w" `
  .
```

### Универсальная сборка (build-universal.sh)

```bash
#!/bin/bash
# Сборка для разных платформ

# Windows 64-bit
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o bin/vmm_windows_amd64.dll .

# Windows 32-bit
GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o bin/vmm_windows_386.dll .

# Linux 64-bit
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/vmm_linux_amd64.so .
```

### Требования для сборки

1. **Go 1.26+**
   ```bash
   go version
   ```

2. **CGO включен** (для компиляции DLL)
   ```bash
   # Linux: нужны cross-compile tools
   # Windows: встроено
   ```

3. **MinGW (для Linux → Windows)**
   ```bash
   # Ubuntu/Debian
   sudo apt-get install mingw-w64
   ```

### Проблемы сборки

#### Ошибка: "main is undeclared"
- Решение: Убедиться что в main.go есть функция main() пусть пустая

#### Ошибка CGO на Linux
- Решение: Установить mingw-w64 или отключить CGO

#### Размер DLL слишком большой
- Решение: Используйте флаги `-ldflags="-s -w"` для удаления debug информации

---

## Коды ошибок (Справка)

### Полный список кодов ошибок

```
Код    | Константа                  | Значение по-русски
-------|----------------------------|---------------------------
-1     | ErrCodeFileNotFound        | Файл не найден
-2     | ErrCodeOutOfMemory         | Нехватка памяти
-3     | ErrCodeIndexOutOfRange     | Индекс вне диапазона
-4     | ErrCodeFileOperation       | Ошибка файловой операции
-5     | ErrCodeInvalidType         | Неверный тип данных
-6     | ErrCodeInsufficientDisk    | Недостаточно дискового пространства
-7     | ErrCodeInvalidHandle       | Неверный handle
-8     | ErrCodePageNotFound        | Страница не найдена
-999   | (неизвестная ошибка)       | Неклассифицированная ошибка
```

### Когда возвращается ошибка 49 (ErrCodeFileOperation)

```
Операция        | Условие
----------------|------------------------------------------
VMCreate        | Файл уже существует
                | Ошибка при создании файла
                | Ошибка при записи
                |
VMOpen          | Файл не существует
                | Файл уже открыт (одновременно)
                | Ошибка при открытии
                |
VMClose         | Ошибка при сохранении
                | Ошибка при закрытии файла
                |
VMRead          | Ошибка при чтении со диска
                | Поврежденный файл
                |
VMWrite         | Недостаточно места (очень редко)
                | Ошибка при записи
                |
ReadPage/       | Ошибка I/O операции
WritePage       | Синхронизация с диском
```

### Диагностика ошибок

```cpp
// Проверить код ошибки
Result result = VMOpen("file.vm");

if (!result.IsSuccess()) {
    int errorCode = result.ErrorCode;
    string errorMsg = result.GetErrorMessage();
    
    switch (errorCode) {
        case -1:
            Console.WriteLine("File not found: " + errorMsg);
            break;
        case -3:
            Console.WriteLine("Index out of range: " + errorMsg);
            break;
        case -4:
            Console.WriteLine("File operation error: " + errorMsg);
            // Возможные причины:
            // - Файл уже открыт
            // - Нет доступа к файлу
            // - Файл поврежден
            break;
        case -7:
            Console.WriteLine("Invalid handle: " + errorMsg);
            break;
        default:
            Console.WriteLine("Unknown error: " + errorMsg);
            break;
    }
}
```

---

## Архитектурные решения

### 1. Выбор LRU для кэша

**Почему LRU (Least Recently Used)?**
- Простая в реализации
- Хорошая производительность в большинстве случаев
- Страницы которые часто используются остаются в памяти

**Альтернативы:**
- FIFO: проще но хуже производительность
- Clock: похожа на LRU но дешевле по CPU
- LFU: отслеживает частоту, но сложнее

### 2. Размер страницы 512 байт

**Почему?**
- Стандартный размер блока на диске
- Минимизирует фрагментацию
- Баланс между памятью и производительностью

### 3. Bitmap для отслеживания заполненности

**Почему битовая карта?**
- Минимальный расход памяти (1 бит на элемент)
- Быстрые операции AND/OR/XOR
- Легко найти первый свободный элемент

### 4. Отдельное хранилище для varchar

**Почему отдельный файл?**
- Строки переменной длины не имеют фиксированного размера
- Облегчает управление памятью
- Основной файл остается структурированным

### 5. Мьютекс на уровне API

**Почему?**
- Потокобезопасность для многопоточных приложений C#
- Предотвращает race conditions при открытии файлов
- Гарантирует целостность данных

---

## Производительность

### Оптимизация доступа

```
Сценарий                    | Оптимизация
----------------------------|------------------------------------------
Последовательный доступ     | LRU кэш - все страницы остаются
                            | Hit rate близка к 100%
                            |
Случайный доступ            | Зависит от размера кэша
                            | Для 1000 элементов, кэш 10 стр: 
                            | примерно 1-10% hit rate
                            |
Повторяющийся доступ        | Идеален для LRU
                            | Pages становятся "горячими"
                            | Hit rate может быть 90%+
```

### Примеры производительности

```
Операция                    | Время (на примерном ПК)
----------------------------|------------------------------------------
VMCreate (1000 элементов)   | ~10 мс
VMOpen                      | ~5 мс
VMRead (hit cache)          | ~0.1 мс
VMRead (miss cache)         | ~5 мс (reading from disk)
VMWrite (hit cache)         | ~0.2 мс
VMClose (no dirty pages)    | ~1 мс
VMClose (flush dirty pages) | ~10-50 мс (зависит от количества)
```

---

## Развертывание

### Подготовка DLL для использования в C#

1. **Скомпилировать**
   ```bash
   ./build.sh
   ```

2. **Скопировать DLL**
   ```bash
   cp VirtualMemoryManagement.dll /path/to/cs/project/
   ```

3. **Использовать в C#**
   ```csharp
   [DllImport("VirtualMemoryManagement.dll")]
   public static extern Result VMCreate(string filename, int size, string type, int stringLength);
   ```

### Развертывание на целевой машине

1. Скопировать:
   - `VirtualMemoryManagement.dll`
   - C# приложение
   - Все зависимости .NET

2. На Windows 11:
   - DLL должна быть в PATH или в папке приложения
   - Требуется .NET Runtime

---

## Развитие проекта

### Планируемые улучшения

1. **Оптимизация памяти**
   - Использование mmap для больших файлов
   - Pooling страниц

2. **Оптимизация производительности**
   - Асинхронный I/O
   - Параллельный кэш

3. **Новые типы данных**
   - Double/Float
   - Bool
   - Custom structures

4. **Расширенное кэширование**
   - Prefetching
   - Адаптивный размер кэша

5. **Мониторинг**
   - Логирование
   - Профилирование
   - Метрики производительности

---

## Лицензия и авторские права

Проект разработан в рамках лабораторной работы.

---

## Контактная информация

Для вопросов и предложений: обратитесь к разработчику проекта.

---

**Версия документации:** 1.0
**Дата обновления:** March 18, 2026
**Совместимость:** Go 1.26+, Windows/Linux

