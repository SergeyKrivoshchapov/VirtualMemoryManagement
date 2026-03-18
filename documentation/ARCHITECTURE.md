# Архитектурная документация - Virtual Memory Management System

## Обзор системы

Virtual Memory Management System (VMMS) - это комплексная система для управления виртуальной памятью, позволяющая работать с большими массивами данных через механизм пагинации и кэширования.

---

## 1. Слой API (api/vmm_api.go)

### Назначение
Экспортирует функции для использования из C# приложения через DLL.

### Основные компоненты

#### Глобальное состояние (Thread-safe)
```go
var (
    mu          sync.Mutex                    // Мьютекс для синхронизации
    handles     = make(map[int]*VirtualArray) // Таблица открытых файлов
    nextID      = 1                           // Счетчик для ID handles
    cacheSize   = config.DefaultCacheSize     // Размер кэша
    fileHandles = make(map[string]int)        // Отображение filename → handle
)
```

**fileHandles** - критическое дополнение!
- Предотвращает открытие одного файла дважды
- При попытке открыть файл, проверяется наличие в fileHandles
- При открытии: добавляется запись filename → handle
- При закрытии: удаляется из fileHandles

#### Функции управления кэшем
```go
func SetCacheSize(size int)  // Установить размер кэша (3-100)
func GetCacheSize() int       // Получить текущий размер
```

#### Функции работы с файлами
```go
func VMCreate(filename, size, typ, stringLength)  // Создать
func VMOpen(filename)                              // Открыть
func VMClose(handle)                               // Закрыть
```

#### Функции работы с данными
```go
func VMRead(handle, index)          // Чтение
func VMWrite(handle, index, value)  // Запись
```

#### Вспомогательные функции
```go
func VMStats(handle)      // Статистика
func GetHandle()           // Первый handle
func GetAllHandles()       // Все handles
func VMHelp(filename)      // Справка
```

### Инвариант: Одновременно открытый файл

```
Сценарий 1: ЗАПРЕЩЕНО
  h1 = VMOpen("file.vm")     // Успех, h1 = 1
  h2 = VMOpen("file.vm")     // Ошибка -4, "file already opened"

Сценарий 2: РАЗРЕШЕНО
  h1 = VMOpen("file1.vm")    // Успех, h1 = 1
  h2 = VMOpen("file2.vm")    // Успех, h2 = 2

Сценарий 3: РАЗРЕШЕНО (после закрытия)
  h1 = VMOpen("file.vm")     // Успех
  VMClose(h1)
  h1 = VMOpen("file.vm")     // Успех (файл больше не в fileHandles)
```

### Поток выполнения VMCreate

```
VMCreate("file.vm", 1000, "int", 0)
  ↓
1. Мьютекс LOCK
2. Проверить: файл не существует?
3. Преобразовать тип: "int" → array.TypeInt
4. Валидация размера: size > 0?
5. Разблокировка мьютекса
6. CreateWithCacheSize(filename, size, type, 0, cacheSize)
   ├─ Создать PageFile
   ├─ Инициализировать структуру виртуального массива
   └─ Загрузить начальные страницы
7. Закрыть файл (на этапе создания)
8. Вернуть Success = 1
```

### Поток выполнения VMOpen

```
VMOpen("file.vm")
  ↓
1. Мьютекс LOCK
2. Проверить: файл уже открыт в fileHandles?
   ├─ ДА → Вернуть Error -4 "file already opened"
   └─ НЕТ → Продолжить
3. Разблокировка мьютекса (ДО дорогой операции!)
4. OpenWithCacheSize(filename, cacheSize)
   ├─ Открыть файл для чтения/записи
   ├─ Прочитать заголовок
   ├─ Инициализировать кэш
   └─ Загрузить varchar индекс (если нужно)
5. Мьютекс LOCK (опять заблокировать)
6. Присвоить ID: id = nextID++
7. Добавить в таблицы:
   ├─ handles[id] = virtualArray
   └─ fileHandles[filename] = id
8. Разблокировка мьютекса
9. Вернуть Success = 1, Data = string(id)
```

**Важно:** Мьютекс разблокируется ПЕРЕД дорогой операцией открытия файла, чтобы не блокировать другие потоки.

### Поток выполнения VMClose

```
VMClose(handle)
  ↓
1. Мьютекс LOCK
2. Проверить: handle существует в handles?
   ├─ НЕТ → Вернуть Error -7 "Invalid handle"
   └─ ДА → Продолжить
3. Найти и удалить из fileHandles (по handle)
4. Удалить из handles[handle]
5. Разблокировка мьютекса
6. FlushDirtyPages() - сохранить все грязные страницы
7. Close() - закрыть файл
8. Вернуть Success = 1
```

### Поток выполнения VMRead

```
VMRead(handle, 10)
  ↓
1. Мьютекс LOCK (быстрая проверка)
2. Получить VirtualArray из handles[handle]
3. Разблокировка мьютекса
4. va.Read(10)
   ├─ Вычислить pageNum = 10 / elementsPerPage
   ├─ Получить/загрузить страницу (из кэша или диска)
   ├─ Получить элемент из страницы по смещению
   └─ Вернуть value
5. Преобразовать значение в строку (int32 → "42", string → "string")
6. Вернуть Success = 1, Data = "значение"
```

### Поток выполнения VMWrite

```
VMWrite(handle, 10, "42")
  ↓
1. Мьютекс LOCK (быстрая проверка)
2. Получить VirtualArray из handles[handle]
3. Получить arrayInfo
4. Разблокировка мьютекса
5. Преобразовать "42" согласно типу:
   ├─ int: ParseInt("42") → int32(42)
   ├─ char: "42" → string
   └─ varchar: "42" → string
6. va.Write(10, value)
   ├─ Вычислить pageNum = 10 / elementsPerPage
   ├─ Получить/загрузить страницу
   ├─ Записать значение в страницу
   ├─ Отметить Dirty = true
   └─ Переместить страницу в начало кэша
7. Вернуть Success = 1
```

---

## 2. Слой Virtual Array (virtualmemory/virtual_array.go)

### Назначение
Реализует логику виртуальной памяти - управление страницами, кэшем, чтением/записью.

### Структура VirtualArray

```go
type VirtualArray struct {
    pageStorage  storage.PageStorage   // Хранилище страниц (на диск)
    varcharStore storage.VarcharStorage // Хранилище строк (для varchar)
    arrayInfo    *array.Info            // Информация о массиве
    pageCache    cache.Cache            // LRU кэш страниц
    varcharIndex map[int]int64          // Индекс для varchar строк
    cacheSize    int                    // Размер кэша
}
```

### Ключевые методы

#### Create/Open

```go
func CreateWithCacheSize(filename, size, typ, stringLength, cacheSize)
  ↓
1. Валидация размера (size > 0)
2. Ограничение cacheSize (3-100)
3. Создать PageFile и инициализировать
4. Создать arrayInfo (из заголовка)
5. Создать LRU кэш с нужным размером
6. Если varchar: создать VarcharFile
7. Загрузить начальные страницы (loadInitialPages)
8. Вернуть VirtualArray
```

```go
func OpenWithCacheSize(filename, cacheSize)
  ↓
1. Ограничить cacheSize
2. Открыть PageFile (читать заголовок)
3. Получить arrayInfo из файла
4. Создать LRU кэш
5. Если varchar: открыть VarcharFile
6. Загрузить начальные страницы
7. Загрузить индекс varchar (если нужно)
8. Вернуть VirtualArray
```

#### Read

```go
func (va *VirtualArray) Read(index int) (interface{}, error)
  ↓
1. Валидация: 0 <= index < va.arrayInfo.Size
2. Вычислить:
   ├─ pageNum = index / BitsPerPage
   ├─ elemOffset = index % BitsPerPage
   └─ byteOffset = elemOffset * ElementSize
3. Получить страницу: getPage(pageNum)
   ├─ Проверить кэш: pageCache.Get(pageNum)
   ├─ Если НЕ найдена:
   │  ├─ Прочитать со диска: pageStorage.ReadPage(pageNum)
   │  ├─ Добавить в кэш: pageCache.Put(page)
   │  │   └─ Если вытеснена старая страница (Dirty):
   │  │       └─ WritePage(oldPage) - сохранить на диск
   │  └─ Обновить accessCount++
   └─ Вернуть страницу
4. Проверить bitmap: заполнен ли элемент?
   ├─ НЕТ → вернуть zero value
   └─ ДА → продолжить
5. Получить значение из Data:
   ├─ int: binary.Read → int32
   ├─ char: string из Data
   └─ varchar: получить из VarcharFile по индексу
6. Вернуть значение
```

#### Write

```go
func (va *VirtualArray) Write(index int, value interface{}) error
  ↓
1. Валидация индекса
2. Преобразовать value в нужный тип
3. Вычислить pageNum и смещения (как в Read)
4. Получить страницу (из кэша или диска)
5. Обновить bitmap: отметить элемент как заполненный
6. Записать значение в Data:
   ├─ int: binary.Write → bytes
   ├─ char: copy bytes в Data
   └─ varchar: write в VarcharFile, сохранить индекс
7. Отметить страницу: Dirty = true
8. Переместить в начало кэша (moveToFront)
9. Вернуть nil (без ошибок)
```

#### Flush и Close

```go
func (va *VirtualArray) FlushDirtyPages() error
  ↓
1. Получить все страницы из кэша
2. Для каждой страницы:
   ├─ Если Dirty = true:
   │  ├─ WritePage(page) - сохранить на диск
   │  ├─ Dirty = false
   │  └─ Обновить accessCount
   └─ Продолжить
3. Если VarcharFile: сохранить varchar индекс
4. Вернуть nil
```

```go
func (va *VirtualArray) Close() error
  ↓
1. Вызвать FlushDirtyPages() - сохранить все грязные
2. Закрыть PageFile
3. Если VarcharFile: закрыть его
4. Очистить кэш
5. Вернуть nil
```

---

## 3. Слой Cache (cache/lru.go)

### Архитектура LRU

```
Структура двусвязного списка:

  Head ↔ Node1 ↔ Node2 ↔ Node3 ↔ Tail
 (dummy) (самая   ...   (самая  (dummy)
         новая)         старая)

pageMap = {pageNum → Node}

Операция Get(pageNum):
  ├─ Найти Node в pageMap
  ├─ Переместить в начало (после Head)
  └─ Вернуть page

Операция Put(page):
  ├─ Если существует:
  │  ├─ Обновить page
  │  └─ Переместить в начало
  └─ Иначе:
     ├─ Если кэш полон:
     │  ├─ Найти последний Node (перед Tail)
     │  ├─ Удалить из списка
     │  ├─ Удалить из pageMap
     │  └─ Вернуть evicted page
     └─ Создать новый Node
        ├─ Добавить в pageMap
        ├─ Добавить в начало списка
        └─ Вернуть nil (nothing evicted)
```

### Методы

```go
func (lru *LRUCache) Get(pageNumber int) *page.Page
  ├─ O(1) - поиск в map
  ├─ Переместить в начало - O(1)
  └─ Вернуть page или nil

func (lru *LRUCache) Put(p *page.Page) *page.Page
  ├─ O(1) - все операции в map и list
  └─ Вернуть evicted page или nil

func (lru *LRUCache) Contains(pageNumber int) bool
  └─ O(1) - check map

func (lru *LRUCache) All() []*page.Page
  └─ O(n) - собрать все страницы
```

### Использование в VirtualArray

```
Read/Write:
  1. Попробовать Get из кэша
  2. Если miss:
     ├─ ReadPage со диска
     └─ Put в кэш
        ├─ Если вытеснена (Dirty):
        │  └─ WritePage на диск
        └─ Вернуть новую страницу
  3. Работать с страницей

Close:
  1. Получить All() все страницы
  2. Для каждой (Dirty):
     └─ WritePage на диск
  3. Очистить кэш
```

---

## 4. Слой Storage (storage/)

### PageFile

#### Структура файла
```
Offset  Size   Content
──────────────────────────────────
0       2      Signature "VM"
2       8      Size (int64) - кол-во элементов
10      1      Type (byte) - 'I', 'C', 'V'
11      4      StringLength (int32)
15      ?      Page 0
? +512  512    Page 1
? +512  512    Page N
...
```

#### Header.WriteTo
```go
func (h *Header) WriteTo(w io.Writer) error
  ├─ Write Size (8 bytes, little-endian)
  ├─ Write Type (1 byte)
  └─ Write StringLength (4 bytes)
```

#### PageFile.ReadPage

```
ReadPage(pageNum):
  ├─ Вычислить offset:
  │  └─ offset = HeaderSize + (pageNum * PageSize)
  ├─ Seek(offset)
  ├─ Read PageSize bytes
  ├─ Разбить:
  │  ├─ Первые 16 байт → Bitmap
  │  └─ Остальное → Data
  └─ Вернуть Page{bitmap, data, pageNum, false}

Размер:
  HeaderSize = 2 + 8 + 1 + 4 = 15 bytes
  PageSize = 512 bytes
```

#### PageFile.WritePage

```
WritePage(page):
  ├─ Вычислить offset
  ├─ Seek(offset)
  ├─ Write Bitmap (16 bytes)
  ├─ Write Data (496 bytes)
  └─ Sync() - синхронизировать с диском
```

### VarcharFile

#### Назначение
Хранилище для строк переменной длины в varchar массивах.

#### Структура

```
Файл: filename.vm.varchar

Содержимое:
  [String 0]
  [String 1]
  [String 2]
  ...
  [String N]

varcharIndex = {
  0: 0,           // String 0 начинается с offset 0
  1: 15,          // String 1 начинается с offset 15
  2: 30,          // String 2 начинается с offset 30
  ...
}
```

#### Методы

```go
func (vf *VarcharFile) WriteString(value string) int64
  ├─ Получить текущее смещение
  ├─ Написать length (4 байта, int32)
  ├─ Написать value (len bytes)
  ├─ Синхронизировать
  └─ Вернуть offset начала

func (vf *VarcharFile) ReadString(offset int64) (string, error)
  ├─ Seek(offset)
  ├─ Прочитать length (int32)
  ├─ Прочитать bytes (length)
  └─ Вернуть string
```

### BinaryIO

Вспомогательные функции для бинарных операций.

```go
func ReadInt32(data []byte, offset int) int32
  └─ binary.LittleEndian.Uint32 + cast

func WriteInt32(data []byte, offset int, value int32) error
  └─ binary.LittleEndian.PutUint32

func ReadString(data []byte, offset, length int) string
  └─ string(data[offset:offset+length])

func WriteString(data []byte, offset int, value string, length int) error
  └─ copy(data[offset:], value)
```

---

## 5. Типы данных (types/)

### Array Type

```go
type Type byte

const (
    TypeInt     Type = 'I'  // Целые числа (int32)
    TypeChar    Type = 'C'  // Строки фиксированной длины
    TypeVarchar Type = 'V'  // Строки переменной длины
)

type Info struct {
    Size         int  // Кол-во элементов
    Type         Type // Тип
    StringLength int  // Длина строки
    ElementSize  int  // Размер одного элемента в байтах
    PageCount    int  // Кол-во страниц
}
```

**ElementSize вычисляется:**
```
TypeInt:     4 (sizeof int32)
TypeChar:    stringLength (фиксированно)
TypeVarchar: 4 (sizeof указателя на строку)
```

**PageCount вычисляется:**
```
PageCount = ceil(Size / BitsPerPage)
         = (Size + BitsPerPage - 1) / BitsPerPage
         = (Size + 127) / 128
```

### Page

```go
type Page struct {
    Bitmap         [16]byte         // Bitmap (128 бит)
    Data           []byte           // Данные элементов
    AbsoluteNumber int              // Номер страницы
    Dirty          bool             // Был ли изменен
    AccessCount    int              // Кол-во обращений
}
```

**Bitmap:**
- 128 бит = 16 байт
- Бит i = 1 если элемент i заполнен
- Бит i = 0 если элемент i пуст (zero value)

**Data:**
- Размер = BitsPerPage * ElementSize = 128 * ElementSize
- Элемент i находится по смещению i * ElementSize

### Bitmap

```go
type Bitmap struct {
    bits [16]byte  // 128 бит
}

func (b *Bitmap) Set(index int)
  └─ bits[index / 8] |= (1 << (index % 8))

func (b *Bitmap) Get(index int) bool
  └─ (bits[index / 8] & (1 << (index % 8))) != 0

func (b *Bitmap) Clear(index int)
  └─ bits[index / 8] &= ^(1 << (index % 8))
```

### Result

```go
type Result struct {
    Success   int        // 1 = успех, 0 = ошибка
    Data      [256]byte  // Данные результата (строка)
    ErrorCode int        // Код ошибки
}

func Success(data string) Result
  └─ Success=1, copy Data

func Error(err error) Result
  ├─ Success=0
  ├─ ErrorCode=GetErrorCode(err)
  └─ Data=err.Error()

func ErrorWithCode(code int, message string) Result
  ├─ Success=0
  ├─ ErrorCode=code
  └─ Data=message
```

---

## 6. Конфигурация (config/)

```go
const (
    BitsPerPage      = 128  // Элементов на странице
    BytesPerBitmap   = 16   // Размер bitmap (128/8)
    PhysicalPageSize = 512  // Размер физической страницы

    MinCacheSize     = 3    // Минимум страниц в кэше
    MaxCacheSize     = 100  // Максимум страниц в кэше
    DefaultCacheSize = 10   // По умолчанию
)
```

---

## 7. Обработка ошибок (errors/)

```go
type VMMError struct {
    Code    int    // Код ошибки
    Message string // Сообщение
    Err     error  // Исходная ошибка
}

const (
    ErrCodeFileNotFound     = -1
    ErrCodeOutOfMemory      = -2
    ErrCodeIndexOutOfRange  = -3
    ErrCodeFileOperation    = -4
    ErrCodeInvalidType      = -5
    ErrCodeInsufficientDisk = -6
    ErrCodeInvalidHandle    = -7
    ErrCodePageNotFound     = -8
)
```

**Преобразование ошибок:**
```
os.ErrNotExist → ErrCodeFileNotFound (-1)
io.ErrUnexpectedEOF → ErrCodeFileOperation (-4)
VMMError → VMMError.Code
другое → -999
```

---

## 8. DLL Wrapper (dll_wrapper.go)

### Назначение
Экспортировать функции для использования из C# через CGO.

### Типы CGO

```go
// cgo_types.go

type CResult struct {
    Success   C.int
    Data      [256]C.char
    ErrorCode C.int
}
```

### Функции экспорта

```go
// Каждая публичная функция в api должна быть обернута:

//export VMCreate
func VMCreate(cFilename *C.char, cSize C.int, cType *C.char, cStringLength C.int) CResult {
    filename := C.GoString(cFilename)
    typ := C.GoString(cType)
    
    result := api.VMCreate(filename, int(cSize), typ, int(cStringLength))
    
    return CResult{
        Success:   C.int(result.Success),
        ErrorCode: C.int(result.ErrorCode),
        // копировать Data...
    }
}
```

---

## 9. Интеграция с C# (vmm.h)

```c
// Заголовок для импорта в C#

typedef struct {
    int Success;
    char Data[256];
    int ErrorCode;
} Result;

Result VMCreate(const char* filename, int size, const char* type, int stringLength);
Result VMOpen(const char* filename);
Result VMClose(int handle);
Result VMRead(int handle, int index);
Result VMWrite(int handle, int index, const char* value);
Result VMStats(int handle);
```

---

## 10. Поток данных

### Создание файла

```
C# Application
  ↓
VMCreate("file.vm", 100, "int", 0)
  ↓ (DLL wrapper)
api.VMCreate(...)
  ↓
virtualmemory.CreateWithCacheSize(...)
  ├─ pageFile.Create()
  │  ├─ os.Create("file.vm")
  │  ├─ write signature "VM"
  │  ├─ write header (size, type, stringLength)
  │  └─ write empty pages (zeros)
  ├─ if varchar: varcharFile.Create()
  ├─ loadInitialPages()
  │  ├─ ReadPage(0)
  │  ├─ ReadPage(1)
  │  └─ Put в кэш
  └─ Close (при создании)
  ↓
Result{Success: 1, Data: "Created"}
```

### Чтение элемента

```
C# Application
  ↓
VMRead(1, 10)  // handle=1, index=10
  ↓ (DLL wrapper)
api.VMRead(1, 10)
  ├─ mu.Lock()
  ├─ va = handles[1]
  ├─ mu.Unlock()
  ├─ va.Read(10)
  │  ├─ pageNum = 10 / 128 = 0
  │  ├─ elemOffset = 10 % 128 = 10
  │  ├─ pageCache.Get(0)
  │  │  ├─ найдена в кэше → return page
  │  │  └─ не найдена:
  │  │     ├─ pageStorage.ReadPage(0)
  │  │     │  ├─ seek(15 + 0*512) = 15
  │  │     │  ├─ read 512 bytes
  │  │     │  └─ parse bitmap + data
  │  │     ├─ pageCache.Put(page)
  │  │     │  ├─ если кэш полон:
  │  │     │  │  └─ evict old page (если Dirty → save)
  │  │     │  └─ add to pageMap
  │  │     └─ return page
  │  ├─ bitmap.Get(10) → check if filled
  │  ├─ get value from page.Data[10*4:(10+1)*4]
  │  ├─ parse int32
  │  └─ return value
  ├─ convert to string "value"
  └─ Result{Success: 1, Data: "value"}
```

### Запись элемента

```
C# Application
  ↓
VMWrite(1, 10, "42")
  ↓
api.VMWrite(1, 10, "42")
  ├─ Parse type, convert "42" → int32(42)
  ├─ va.Write(10, 42)
  │  ├─ pageNum = 0, elemOffset = 10
  │  ├─ pageCache.Get(0) [как выше]
  │  ├─ bitmap.Set(10) → отметить заполненным
  │  ├─ write int32(42) в page.Data[10*4:(10+1)*4]
  │  ├─ page.Dirty = true
  │  ├─ pageCache.moveToFront() → обновить позицию
  │  └─ return
  └─ Result{Success: 1, Data: "Written"}

После VMClose(1):
  ├─ va.FlushDirtyPages()
  │  └─ для каждой page (Dirty = true):
  │     ├─ pageStorage.WritePage(page)
  │     │  ├─ seek(15 + pageNum*512)
  │     │  ├─ write bitmap
  │     │  ├─ write data
  │     │  └─ sync
  │     └─ page.Dirty = false
  └─ закрыть файл
```

---

## 11. Инварианты

### 1. Один файл - одна ручка одновременно

```
invariant: ∀file ∈ openFiles: count(handles[h].filename == file) ≤ 1

Поддержка: fileHandles[filename] → handle
```

### 2. Целостность данных

```
invariant: ∀page: if page ∈ cache ∧ page.Dirty:
           ∀close: page будет сохранена на диск перед закрытием

Поддержка: FlushDirtyPages() при Close
```

### 3. Консистентность bitmap

```
invariant: ∀page, index: bitmap[index] = 1 ↔ data[index] != zero

Поддержка: bitmap.Set при Write, bitmap.Clear при Delete
```

### 4. Размер кэша в пределах

```
invariant: MinCacheSize ≤ cacheSize ≤ MaxCacheSize

Поддержка: SetCacheSize ограничивает значение
```

---

## 12. Производительность

### Сложность операций

```
Операция              | Сложность | Примечания
──────────────────────┼───────────┼──────────────────
VMCreate              | O(n)      | n = pageCount (write pages)
VMOpen                | O(1)      | Читает только заголовок
VMClose               | O(m)      | m = dirty pages
VMRead (hit cache)    | O(1)      | Из памяти
VMRead (miss cache)   | O(1)      | Disk I/O + O(1) cache operations
VMWrite (hit cache)   | O(1)      | Запись в памяти
VMWrite (miss cache)  | O(1)      | Disk I/O + O(1) cache operations
```

### Cache Hit Rate

```
Сценарий                | Hit Rate
──────────────────────┼──────────
Последовательный доступ | ~95%+
Часто повторяющиеся    | ~80-90%
Случайный доступ (N элементов, 
  K страниц в кэше)     | ~K/(N/128) * 100%
```

---

## 13. Расширения и улучшения

### Возможные оптимизации

1. **Prefetching** - загружать соседние страницы
2. **Adaptive caching** - динамический размер кэша
3. **Read-ahead** - прогнозирующее чтение
4. **Compression** - сжатие на диске
5. **Async I/O** - асинхронные операции чтения/записи

### Новые типы данных

- Double/Float
- Boolean
- Byte/Short
- Structures/Tuples

### Мониторинг

- Логирование операций
- Профилирование CPU/Memory
- Метрики производительности

---

**Версия:** 1.0
**Дата:** March 18, 2026

