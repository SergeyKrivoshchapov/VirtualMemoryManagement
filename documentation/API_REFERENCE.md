# API Quick Reference Card

## Функции управления файлами

### VMCreate
```
Создает новый виртуальный массив на диске

Result VMCreate(string filename, int size, string type, int stringLength)

Параметры:
  filename (string)    - Путь и имя файла (например "data.vm")
  size (int)          - Количество элементов в массиве (1-2147483647)
  type (string)       - Тип данных: "int"/"I", "char"/"C", "varchar"/"V"
  stringLength (int)  - Длина строки
                        - Для int: 0 (не используется)
                        - Для char: фиксированная длина (1-256)
                        - Для varchar: максимальная длина (1-256)

Возвращает:
  Success = 1, Data = "Created"     - Успешно создан
  Success = 0, ErrorCode = -4       - Файл уже существует
  Success = 0, ErrorCode = -5       - Неверный тип

Примеры:
  VMCreate("ints.vm", 10000, "int", 0)           // Int array
  VMCreate("names.vm", 100, "char", 20)          // Fixed strings
  VMCreate("description.vm", 1000, "varchar", 100)  // Variable strings
```

### VMOpen
```
Открывает существующий виртуальный массив

Result VMOpen(string filename)

Параметры:
  filename (string) - Путь и имя файла

Возвращает:
  Success = 1, Data = "handle_id"  - Успешно открыт (Data строка с числом)
  Success = 0, ErrorCode = -1      - Файл не найден
  Success = 0, ErrorCode = -4      - Файл уже открыт в этой сессии

Примечания:
  - Один файл может быть открыт только один раз одновременно
  - Handle используется для всех последующих операций
  - Вернуть handle: int handle = int.Parse(new string(result.Data).TrimEnd('\0'))

Примеры:
  Result res = VMOpen("data.vm");
  int handle = int.Parse(new string(res.Data).TrimEnd('\0'));
```

### VMClose
```
Закрывает виртуальный массив и сохраняет все изменения

Result VMClose(int handle)

Параметры:
  handle (int) - Идентификатор открытого файла (от VMOpen)

Возвращает:
  Success = 1, Data = "Closed"     - Успешно закрыт
  Success = 0, ErrorCode = -7      - Неверный handle

Действия:
  1. Проверка валидности handle
  2. Сохранение всех грязных страниц на диск
  3. Синхронизация файловой системы
  4. Закрытие файла
  5. Освобождение памяти

Примеры:
  VMClose(1);
  
Важно:
  - Всегда вызывайте VMClose после работы
  - Все несохраненные данные будут потеряны если не вызвать
```

---

## Функции работы с данными

### VMRead
```
Читает значение элемента по индексу

Result VMRead(int handle, int index)

Параметры:
  handle (int) - Идентификатор открытого файла
  index (int)  - Индекс элемента (0 до size-1)

Возвращает:
  Success = 1, Data = "value"      - Успешно прочитано
                                     Data содержит строковое представление
  Success = 0, ErrorCode = -7      - Неверный handle
  Success = 0, ErrorCode = -3      - Индекс вне диапазона
  Success = 0, ErrorCode = -4      - Ошибка при чтении с диска

Процесс:
  1. Определить номер страницы в памяти
  2. Проверить LRU кэш
  3. Если не в кэше - прочитать со диска
  4. Получить значение из страницы
  5. Преобразовать в строку

Примеры:
  // Прочитать int
  Result res = VMRead(1, 10);
  if (res.Success == 1) {
      int value = int.Parse(new string(res.Data).TrimEnd('\0'));
  }
  
  // Прочитать строку
  Result res = VMRead(2, 5);
  if (res.Success == 1) {
      string value = new string(res.Data).TrimEnd('\0');
  }

Примечания:
  - Первое чтение с индекса загружает страницу в кэш
  - Последующие чтения с той же страницы быстрые
  - Индекс 0-based (первый элемент это индекс 0)
```

### VMWrite
```
Записывает значение элемента по индексу

Result VMWrite(int handle, int index, string value)

Параметры:
  handle (int)   - Идентификатор открытого файла
  index (int)    - Индекс элемента (0 до size-1)
  value (string) - Значение для записи в строковом формате

Возвращает:
  Success = 1, Data = "Written"    - Успешно записано
  Success = 0, ErrorCode = -7      - Неверный handle
  Success = 0, ErrorCode = -3      - Индекс вне диапазона
  Success = 0, ErrorCode = -5      - Неверный формат для типа

Процесс:
  1. Преобразовать значение согласно типу массива:
     - int: ParseInt(value) → int32
     - char: обрезать/дополнить до нужной длины
     - varchar: сохранить в varchar хранилище
  2. Определить номер страницы
  3. Получить/загрузить страницу из кэша
  4. Записать значение в страницу
  5. Отметить страницу как грязную (Dirty)
  6. При необходимости вытеснить из кэша

Примеры:
  // Записать int
  VMWrite(1, 10, "42");
  VMWrite(1, 11, "100");
  
  // Записать строку
  VMWrite(2, 5, "Hello");
  VMWrite(2, 6, "World");
  
  // Записать пустое значение
  VMWrite(1, 12, "0");

Примечания:
  - Запись всегда в памяти (кэше)
  - На диск сохраняется при вытеснении из кэша
  - При Close вызывается FlushDirtyPages (сохранение всех)
  - Индекс 0-based
  - Длины строк должны соответствовать типу (для char)
```

---

## Функции информации и управления

### VMStats
```
Получает статистику использования виртуального массива

Result VMStats(int handle)

Параметры:
  handle (int) - Идентификатор открытого файла

Возвращает:
  Success = 1, Data = "statistics" - Статистика в текстовом формате
  Success = 0, ErrorCode = -7      - Неверный handle

Содержимое статистики:
  - Размер массива
  - Тип данных
  - Количество страниц
  - Размер кэша и использованные страницы
  - Статистика обращений к кэшу
  - Статистика попаданий/промахов

Примеры:
  Result statsRes = VMStats(1);
  if (statsRes.Success == 1) {
      string stats = new string(statsRes.Data).TrimEnd('\0');
      Console.WriteLine(stats);
  }

Примечания:
  - Возвращает различную информацию для отладки
  - Полезно для оптимизации кэша
  - Может содержать подробности о производительности
```

### SetCacheSize
```
Устанавливает размер LRU кэша

void SetCacheSize(int size)

Параметры:
  size (int) - Желаемый размер кэша (количество страниц)

Действия:
  - Если size < 3: устанавливается 3
  - Если size > 100: устанавливается 100
  - Иначе: используется указанный размер

Рекомендуемые значения:
  3-5:    Минимальное использование памяти
  10:     По умолчанию, хороший баланс
  20-50:  Для интенсивной работы
  100:    Максимум для лучшей производительности

Примеры:
  SetCacheSize(10);   // По умолчанию
  SetCacheSize(50);   // Увеличить для интенсивной работы
  SetCacheSize(3);    // Минимум

Примечание:
  - Влияет на все последующие операции Open
  - Для открытых файлов действует их собственный кэш размер
```

### GetCacheSize
```
Получает текущий размер кэша

int GetCacheSize()

Параметры:
  (нет)

Возвращает:
  int - Текущий размер кэша в страницах (3-100)

Примеры:
  int currentSize = GetCacheSize();
  Console.WriteLine($"Cache size: {currentSize}");
```

### VMHelp
```
Получает справку по системе

Result VMHelp(string filename)

Параметры:
  filename (string) - Может быть любое имя или пусто

Возвращает:
  Success = 1, Data = "help text" - Справка в текстовом формате

Содержимое:
  - Список доступных команд
  - Параметры каждой команды
  - Краткое описание

Примеры:
  Result helpRes = VMHelp("");
  Console.WriteLine(new string(helpRes.Data).TrimEnd('\0'));
```

---

## Таблица типов данных

### Поддерживаемые типы

```
┌─────────────┬──────┬──────────────────┬─────────────────────────┐
│ Тип         │ Код  │ Размер элемента  │ Использование           │
├─────────────┼──────┼──────────────────┼─────────────────────────┤
│ "int"/"I"   │ I    │ 4 байта          │ Целые числа int32       │
│ "char"/"C"  │ C    │ stringLength     │ Строки фиксированной    │
│             │      │ байт             │ длины                   │
│ "varchar"/" │ V    │ 4 байта (указ.)  │ Строки переменной       │
│  V"         │      │ + отдельное      │ длины                   │
│             │      │ хранилище        │                         │
└─────────────┴──────┴──────────────────┴─────────────────────────┘
```

### Создание по типам

```go
// Целые числа
VMCreate("ints.vm", 10000, "int", 0);

// Строки фиксированной длины (20 символов)
VMCreate("names.vm", 100, "char", 20);

// Строки переменной длины (макс 100 символов)
VMCreate("descriptions.vm", 1000, "varchar", 100);
```

---

## Коды ошибок

### Основные ошибки

```
┌──────┬──────────────────────────┬─────────────────────────────────┐
│ Код  │ Ошибка                   │ Что делать                      │
├──────┼──────────────────────────┼─────────────────────────────────┤
│ -1   │ FileNotFound             │ Создать файл: VMCreate          │
│ -3   │ IndexOutOfRange          │ Проверить индекс (0 < idx < sz) │
│ -4   │ FileOperation ⚠️ ЧАСТАЯ   │ Проверить:                      │
│      │                          │  - Файл существует?             │
│      │                          │  - Файл не открыт дважды?       │
│      │                          │  - Диск доступен?               │
│ -5   │ InvalidType              │ Проверить тип данных            │
│ -7   │ InvalidHandle            │ Проверить handle от VMOpen      │
└──────┴──────────────────────────┴─────────────────────────────────┘
```

### Обработка в коде

```csharp
Result result = VMRead(handle, index);

if (result.Success == 1) {
    // Успех - обработать данные
    string value = new string(result.Data).TrimEnd('\0');
} else {
    // Ошибка - обработать в зависимости от кода
    string errorMsg = new string(result.Data).TrimEnd('\0');
    switch (result.ErrorCode) {
        case -1:
            Console.WriteLine("File not found");
            break;
        case -3:
            Console.WriteLine("Index out of range");
            break;
        case -4:
            Console.WriteLine("File operation error");
            break;
        case -7:
            Console.WriteLine("Invalid handle");
            break;
        default:
            Console.WriteLine($"Unknown error: {errorMsg}");
            break;
    }
}
```

---

## Типичный рабочий цикл

```csharp
// 1. Создание
Result createRes = VMCreate("data.vm", 1000, "int", 0);
if (createRes.Success != 1) {
    Console.WriteLine("Create failed");
    return;
}

// 2. Открытие
Result openRes = VMOpen("data.vm");
if (openRes.Success != 1) {
    Console.WriteLine("Open failed");
    return;
}
int handle = int.Parse(new string(openRes.Data).TrimEnd('\0'));

// 3. Работа с данными
for (int i = 0; i < 100; i++) {
    // Записать
    VMWrite(handle, i, i.ToString());
    
    // Прочитать
    Result readRes = VMRead(handle, i);
    if (readRes.Success == 1) {
        string value = new string(readRes.Data).TrimEnd('\0');
        Console.WriteLine($"{i}: {value}");
    }
}

// 4. Получение статистики
Result statsRes = VMStats(handle);
if (statsRes.Success == 1) {
    Console.WriteLine(new string(statsRes.Data).TrimEnd('\0'));
}

// 5. Закрытие (ОБЯЗАТЕЛЬНО!)
VMClose(handle);
```

---

## Распространенные ошибки

```
❌ ОШИБКА: Double Open
  VMOpen("file.vm")   // OK
  VMOpen("file.vm")   // ERROR -4 "file already opened"
  
✅ ПРАВИЛЬНО:
  h1 = VMOpen("file1.vm")
  h2 = VMOpen("file2.vm")

❌ ОШИБКА: Индекс вне диапазона
  VMCreate("file.vm", 100, "int", 0)  // size = 100
  VMRead(handle, 100)  // ERROR -3, валидные индексы 0-99
  
✅ ПРАВИЛЬНО:
  VMRead(handle, 99)   // OK

❌ ОШИБКА: Забыли закрыть
  h = VMOpen("file.vm")
  // ... работа ...
  // Забыли VMClose - данные не сохранены!
  
✅ ПРАВИЛЬНО:
  h = VMOpen("file.vm")
  // ... работа ...
  VMClose(h)  // Обязательно!

❌ ОШИБКА: Неправильное преобразование типа
  VMWrite(handle, 0, "not a number")  // Для int array
  
✅ ПРАВИЛЬНО:
  VMWrite(handle, 0, "42")  // Строка которая parses как число
```

---

## Оптимизационные советы

### Производительность

```
Быстро:
  1. Последовательные чтения (0, 1, 2, 3, ...)
  2. Большой кэш (50-100 страниц)
  3. Минимум Open/Close операций
  4. Группировка операций

Медленно:
  1. Случайные обращения
  2. Маленький кэш (3-5 страниц)
  3. Частые Open/Close
  4. Разреженный доступ
```

### Код оптимизации

```csharp
// Перед интенсивной работой
SetCacheSize(50);  // Больше кэша

// Открыть один раз
int handle = int.Parse(new string(VMOpen("data.vm").Data).TrimEnd('\0'));

// Выполнить много операций
for (int i = 0; i < 100000; i++) {
    VMWrite(handle, i, (i * 2).ToString());
}

// Закрыть один раз
VMClose(handle);

// Вернуть к норме
SetCacheSize(10);
```

---

## Шпаргалка для копирования

```csharp
// Полный шаблон использования
using System;
using System.Runtime.InteropServices;

class VMM {
    [DllImport("VirtualMemoryManagement.dll")]
    static extern void SetCacheSize(int size);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMCreate(string f, int s, string t, int sl);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMOpen(string f);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMWrite(int h, int i, string v);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMRead(int h, int i);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMStats(int h);
    
    [DllImport("VirtualMemoryManagement.dll")]
    static extern Result VMClose(int h);
    
    static void Main() {
        SetCacheSize(10);
        VMCreate("test.vm", 100, "int", 0);
        int h = int.Parse(new string(VMOpen("test.vm").Data).TrimEnd('\0'));
        VMWrite(h, 0, "42");
        Console.WriteLine(new string(VMRead(h, 0).Data).TrimEnd('\0'));
        VMClose(h);
    }
}
```

---

**Версия:** 1.0
**Дата:** March 18, 2026
**Совместимость:** Virtual Memory Management System

