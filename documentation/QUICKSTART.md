# Quick Start Guide - Virtual Memory Management System

## Быстрый старт

### 1. Сборка DLL

#### Linux/macOS
```bash
cd /home/serg/GolandProjects/VirtualMemoryManagement
./build.sh
```

#### Windows PowerShell
```powershell
cd VirtualMemoryManagement
.\build.ps1
```

**Результат:** `VirtualMemoryManagement.dll` (или `.so` для Linux)

---

## API Reference - Краткая справка

### Основные функции

| Функция | Параметры | Возвращает | Описание |
|---------|-----------|-----------|---------|
| `VMCreate` | `filename, size, type, stringLength` | `Result` | Создать массив |
| `VMOpen` | `filename` | `Result(handle)` | Открыть массив |
| `VMClose` | `handle` | `Result` | Закрыть массив |
| `VMRead` | `handle, index` | `Result(value)` | Прочитать элемент |
| `VMWrite` | `handle, index, value` | `Result` | Записать элемент |
| `VMStats` | `handle` | `Result(stats)` | Получить статистику |

---

## Примеры кода

### Пример 1: Простое создание и использование

**C# код:**
```csharp
using System;
using System.Runtime.InteropServices;

class Program {
    [DllImport("VirtualMemoryManagement.dll")]
    public static extern Result VMCreate(string filename, int size, string type, int stringLength);
    
    [DllImport("VirtualMemoryManagement.dll")]
    public static extern Result VMOpen(string filename);
    
    [DllImport("VirtualMemoryManagement.dll")]
    public static extern Result VMClose(int handle);
    
    [DllImport("VirtualMemoryManagement.dll")]
    public static extern Result VMRead(int handle, int index);
    
    [DllImport("VirtualMemoryManagement.dll")]
    public static extern Result VMWrite(int handle, int index, string value);
    
    static void Main() {
        // Создать массив целых чисел (1000 элементов)
        Result createRes = VMCreate("numbers.vm", 1000, "int", 0);
        Console.WriteLine(createRes.Success == 1 ? "Created" : "Error");
        
        // Открыть
        Result openRes = VMOpen("numbers.vm");
        int handle = int.Parse(new string(openRes.Data).TrimEnd('\0'));
        
        // Записать значения
        for (int i = 0; i < 100; i++) {
            VMWrite(handle, i, i.ToString());
        }
        
        // Прочитать
        for (int i = 0; i < 100; i++) {
            Result readRes = VMRead(handle, i);
            Console.WriteLine($"Element {i}: {new string(readRes.Data).TrimEnd('\0')}");
        }
        
        // Закрыть
        VMClose(handle);
    }
}
```

---

### Пример 2: Работа со строками

**C# код:**
```csharp
// Создать массив строк переменной длины (макс 100 символов)
VMCreate("strings.vm", 500, "varchar", 100);
int handle = int.Parse(new string(VMOpen("strings.vm").Data).TrimEnd('\0'));

// Записать строки
VMWrite(handle, 0, "Hello");
VMWrite(handle, 1, "World");
VMWrite(handle, 2, "Virtual Memory");

// Прочитать
for (int i = 0; i < 3; i++) {
    string value = new string(VMRead(handle, i).Data).TrimEnd('\0');
    Console.WriteLine($"String {i}: {value}");
}

VMClose(handle);
```

---

### Пример 3: Обработка ошибок

**C# код:**
```csharp
// Попытка открыть несуществующий файл
Result res = VMOpen("nonexistent.vm");

if (res.Success == 0) {
    int errorCode = res.ErrorCode;
    string errorMsg = new string(res.Data).TrimEnd('\0');
    
    Console.WriteLine($"Error {errorCode}: {errorMsg}");
    
    // Обработка конкретных ошибок
    switch (errorCode) {
        case -1:  // ErrCodeFileNotFound
            Console.WriteLine("File doesn't exist");
            break;
        case -4:  // ErrCodeFileOperation
            Console.WriteLine("File operation failed");
            break;
        case -7:  // ErrCodeInvalidHandle
            Console.WriteLine("Invalid handle");
            break;
    }
}
```

---

### Пример 4: Работа с несколькими файлами

**C# код:**
```csharp
// Создать несколько массивов
VMCreate("file1.vm", 1000, "int", 0);
VMCreate("file2.vm", 1000, "varchar", 50);

// Открыть оба
int h1 = int.Parse(new string(VMOpen("file1.vm").Data).TrimEnd('\0'));
int h2 = int.Parse(new string(VMOpen("file2.vm").Data).TrimEnd('\0'));

// Работать с каждым
VMWrite(h1, 0, "42");
VMWrite(h2, 0, "Hello");

// Закрыть оба
VMClose(h1);
VMClose(h2);

// Переоткрыть и проверить
h1 = int.Parse(new string(VMOpen("file1.vm").Data).TrimEnd('\0'));
Console.WriteLine(new string(VMRead(h1, 0).Data).TrimEnd('\0'));  // "42"

VMClose(h1);
```

---

### Пример 5: Изменение размера кэша

**C# код:**
```csharp
[DllImport("VirtualMemoryManagement.dll")]
public static extern void SetCacheSize(int size);

[DllImport("VirtualMemoryManagement.dll")]
public static extern int GetCacheSize();

// Установить большой кэш для интенсивной работы
SetCacheSize(50);

VMCreate("large.vm", 100000, "int", 0);
int handle = int.Parse(new string(VMOpen("large.vm").Data).TrimEnd('\0'));

// Интенсивная работа с большим кэшем
for (int i = 0; i < 10000; i++) {
    VMWrite(handle, i, (i * 2).ToString());
}

VMClose(handle);

// Вернуть к норме
SetCacheSize(10);
```

---

### Пример 6: Получение статистики

**C# код:**
```csharp
[DllImport("VirtualMemoryManagement.dll")]
public static extern Result VMStats(int handle);

VMCreate("test.vm", 1000, "int", 0);
int handle = int.Parse(new string(VMOpen("test.vm").Data).TrimEnd('\0'));

// Выполнить операции
for (int i = 0; i < 100; i++) {
    VMWrite(handle, i, i.ToString());
}

// Получить статистику
Result statsRes = VMStats(handle);
string stats = new string(statsRes.Data).TrimEnd('\0');

Console.WriteLine("Statistics:");
Console.WriteLine(stats);

VMClose(handle);
```

---

## Типы данных

### Array Types

```
Тип         | Код   | Применение
------------|-------|---------------------------
"int" / "I" | I     | Целые числа (int32)
"char" / "C"| C     | Строки фиксированной длины
"varchar"/"V"| V    | Строки переменной длины
```

### Примеры создания

```csharp
// Целые числа
VMCreate("ints.vm", 10000, "int", 0);

// Строки фиксированной длины (20 символов)
VMCreate("names.vm", 100, "char", 20);

// Строки переменной длины (макс 100 символов)
VMCreate("description.vm", 1000, "varchar", 100);
```

---

## Коды ошибок

| Код | Ошибка | Причина |
|-----|--------|---------|
| -1 | FileNotFound | Файл не существует |
| -2 | OutOfMemory | Нехватка памяти |
| -3 | IndexOutOfRange | Индекс вне диапазона |
| **-4** | **FileOperation** | **Файл уже открыт, не существует, ошибка I/O** |
| -5 | InvalidType | Неверный тип данных |
| -6 | InsufficientDisk | Нет места на диске |
| -7 | InvalidHandle | Неверный handle |
| -8 | PageNotFound | Ошибка при чтении страницы |

### Ошибка -4 (FileOperation) - наиболее частая

```
Сценарий                           | Решение
-----------------------------------|----------------------------
VMCreate на уже существующий файл | Удалить файл перед созданием
VMOpen на файл который уже открыт | Закрыть первый handle
VMOpen на несуществующий файл      | Сначала создать (VMCreate)
Ошибка чтения/записи на диск       | Проверить доступ к диску
```

---

## Конфигурация

### Параметры кэша

```csharp
SetCacheSize(size);  // 3 ≤ size ≤ 100

// Рекомендации:
// 3-5:    Минимальное использование памяти
// 10:     Default, хороший баланс
// 20-50:  Интенсивная работа
// 100:    Максимальная производительность
```

### Размер массива

```csharp
// Рекомендации по размеру:
// 100-1000:      Малые массивы (в памяти)
// 1000-100000:   Средние массивы (смешанно)
// 100000+:       Большие массивы (на диске)

// Формула для page count:
// pages = (size + 127) / 128
```

---

## Оптимизация

### Советы производительности

1. **Последовательный доступ**
   ```csharp
   // Хорошо
   for (int i = 0; i < 1000; i++) {
       VMRead(handle, i);  // Последовательно
   }
   ```

2. **Увеличьте кэш для интенсивной работы**
   ```csharp
   SetCacheSize(50);
   // ... интенсивная работа ...
   SetCacheSize(10);  // Вернуть к норме
   ```

3. **Группировка операций**
   ```csharp
   // Хорошо
   for (int i = 0; i < 100; i++) {
       VMWrite(handle, i, value);  // Все в одном проходе
   }
   
   // Плохо
   VMWrite(handle, 0, value);
   // ... другие операции ...
   VMWrite(handle, 1, value);
   ```

4. **Использование правильного типа**
   ```csharp
   // Хорошо для строк переменной длины
   VMCreate("data.vm", 1000, "varchar", 100);
   
   // Не эффективно для строк переменной длины
   VMCreate("data.vm", 1000, "char", 100);  // Тратит память
   ```

---

## Обработка исключений

### Правильный обработчик ошибок

```csharp
public class VMHelper {
    public static int OpenFile(string filename) {
        Result res = VMOpen(filename);
        
        if (res.Success != 1) {
            throw new Exception(
                $"Failed to open {filename}: {GetErrorMessage(res.ErrorCode)}"
            );
        }
        
        return int.Parse(new string(res.Data).TrimEnd('\0'));
    }
    
    public static string GetErrorMessage(int errorCode) {
        return errorCode switch {
            -1 => "File not found",
            -2 => "Out of memory",
            -3 => "Index out of range",
            -4 => "File operation failed",
            -5 => "Invalid type",
            -6 => "Insufficient disk space",
            -7 => "Invalid handle",
            -8 => "Page not found",
            _ => $"Unknown error {errorCode}"
        };
    }
}

// Использование
try {
    int handle = VMHelper.OpenFile("data.vm");
    // ... работа ...
    VMClose(handle);
} catch (Exception ex) {
    Console.WriteLine($"Error: {ex.Message}");
}
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

# С отчетом о покрытии
go test -cover ./...
```

### Примеры тестов на Go

```go
func TestBasicCreateOpen(t *testing.T) {
    defer os.Remove("test.vm")
    
    // Create
    res := api.VMCreate("test.vm", 100, "int", 0)
    if !res.IsSuccess() {
        t.Fatal("Create failed")
    }
    
    // Open
    res = api.VMOpen("test.vm")
    if !res.IsSuccess() {
        t.Fatal("Open failed")
    }
    
    handle, _ := strconv.Atoi(res.String())
    
    // Close
    res = api.VMClose(handle)
    if !res.IsSuccess() {
        t.Fatal("Close failed")
    }
}

func TestReadWrite(t *testing.T) {
    defer os.Remove("test.vm")
    
    api.VMCreate("test.vm", 100, "int", 0)
    res := api.VMOpen("test.vm")
    handle, _ := strconv.Atoi(res.String())
    
    // Write
    res = api.VMWrite(handle, 0, "42")
    if !res.IsSuccess() {
        t.Fatal("Write failed")
    }
    
    // Read
    res = api.VMRead(handle, 0)
    if !res.IsSuccess() {
        t.Fatal("Read failed")
    }
    
    if res.String() != "42" {
        t.Fatalf("Expected 42, got %s", res.String())
    }
    
    api.VMClose(handle)
}
```

---

## FAQ - Часто задаваемые вопросы

### Q: Какой максимальный размер массива?
**A:** Теоретически неограниченный, ограничен только размером диска. Практически - до 2^31 элементов.

### Q: Как изменить тип данных существующего массива?
**A:** Нельзя напрямую. Нужно создать новый массив и скопировать данные.

### Q: Что произойдет если открыть один файл дважды?
**A:** Вторая попытка вернет ошибку -4 (FileOperation).

### Q: Где хранятся данные?
**A:** В файле `filename.vm` (основной) и `filename.vm.varchar` (для varchar).

### Q: Как я могу улучшить производительность?
**A:** Увеличьте размер кэша (SetCacheSize), используйте последовательный доступ.

### Q: Потокобезопасна ли система?
**A:** Да, все операции защищены мьютексом на уровне API.

### Q: Что случится если я потеряю handle?
**A:** Файл останется открытым, но вы не сможете закрыть его. Вам нужно перезапустить приложение.

---

## Утилиты и инструменты

### Просмотр структуры файла

```bash
# Linux
hexdump -C file.vm | head -20

# Windows PowerShell
Get-Content file.vm -Encoding Byte | Format-Hex -Path file.vm
```

### Удаление файлов

```bash
# Удалить основной файл и varchar индекс
rm file.vm file.vm.varchar
```

### Проверка целостности файла

```go
// Проверить сигнатуру файла
func ValidateFile(filename string) bool {
    f, _ := os.Open(filename)
    defer f.Close()
    
    sig := make([]byte, 2)
    f.Read(sig)
    return string(sig) == "VM"
}
```

---

## Troubleshooting

### Problem: Error -4 при VMCreate

**Решение 1:** Файл уже существует
```csharp
System.IO.File.Delete("file.vm");
System.IO.File.Delete("file.vm.varchar");
VMCreate("file.vm", 1000, "int", 0);
```

**Решение 2:** Нет доступа к папке
```csharp
// Использовать абсолютный путь или папку с правами
VMCreate(@"C:\Users\YourUser\Documents\file.vm", 1000, "int", 0);
```

### Problem: Error -3 при VMRead

**Причина:** Индекс >= размер массива
```csharp
VMCreate("file.vm", 100, "int", 0);
int handle = int.Parse(new string(VMOpen("file.vm").Data).TrimEnd('\0'));

// Неправильно (индекс 100 не существует для size=100)
VMRead(handle, 100);  // Error -3

// Правильно
VMRead(handle, 99);  // OK, индекс 0-99
```

### Problem: Медленная производительность

**Решение:**
```csharp
// Увеличить кэш
SetCacheSize(50);

// Использовать последовательный доступ вместо случайного
for (int i = 0; i < 10000; i++) {
    VMRead(handle, i);  // Хорошо
}

// Избежать частых Open/Close
VMOpen once → do many operations → VMClose once
```

---

**Версия:** 1.0
**Последнее обновление:** March 18, 2026

