# Тестирование VMWrite для CHAR типа

## Как запустить тест с сохранением файла

### Linux/macOS

```bash
cd /home/serg/GolandProjects/VirtualMemoryManagement

# Запустить тест с сохранением файлов
KEEP_TEST_FILES=1 go test -v -run TestVMWriteCharSimple ./api
```

### Windows PowerShell

```powershell
cd VirtualMemoryManagement
$env:KEEP_TEST_FILES="1"; go test -v -run TestVMWriteCharSimple ./api
```

## Где находится файл после теста

В выводе теста будет строка:
```
File location: /tmp/vmm_test_3605693800/test_char_simple
Keeping test files at: /tmp/vmm_test_3605693800
```

Файл находится там! Можно открыть и посмотреть содержимое.

## Структура .vm файла для CHAR

```
Offset  Размер  Содержимое
──────────────────────────────────────
0       2       Сигнатура: "VM" (0x56 0x4D)
2       8       Size = 20 (0x14 в little-endian)
10      1       Type = 'C' (0x43 для CHAR)
11      4       StringLength = 30 (0x1E в little-endian)
15      ~4096   Данные элементов (страницы)
```

## Просмотр содержимого файла

### Просмотр в hex формате

```bash
# Первые 256 байт
hexdump -C /tmp/vmm_test_3605693800/test_char_simple | head -20

# Все содержимое
hexdump -C /tmp/vmm_test_3605693800/test_char_simple
```

### Строки в файле

```bash
# Извлечь все строки
strings /tmp/vmm_test_3605693800/test_char_simple

# Вывод будет содержать:
# Hello
# World
# Test String
# 1234567890
# Short
# This is a longer test string
```

### Размер файла

```bash
ls -lh /tmp/vmm_test_3605693800/test_char_simple

# Должно быть:
# -rw-r--r-- 1 serg serg 4.1K /tmp/vmm_test_3605693800/test_char_simple
```

## Понимание формата

### Header (первые 15 байт)

```
Байты   Значение
──────────────────
0-1     "VM" (сигнатура)
2-9     Size = 20 элементов (int64)
10      Type = 'C' (CHAR)
11-14   StringLength = 30 (int32)
```

### Data (начиная с байта 15)

Для CHAR массива с 20 элементами и длиной строки 30:
- Каждый элемент занимает 30 байт
- Всего элементов на странице = 128
- Страниц нужно = ceil(20/128) = 1 страница
- Размер страницы = 16 (bitmap) + 128*30 (data) = 3856 байт
- Выравнивается до 4096

## Тестовые данные в файле

После запуска теста записанные данные:

```
Index 0: "Hello"              (заполнено, остальное нули)
Index 1: "World"              (заполнено, остальное нули)
Index 2: "Test String"        (заполнено, остальное нули)
Index 3: "1234567890"         (заполнено, остальное нули)
Index 4: "Short"              (заполнено, остальное нули)
Index 5: "This is a longer test string" (заполнено)
Index 6-19: пусто (нули)
```

## Проверка целостности

Когда откроешь файл вручную, можешь:

1. **Проверить сигнатуру:**
   ```bash
   head -c 2 /tmp/vmm_test_3605693800/test_char_simple | od -A x -t x1
   # Должно быть: 56 4d (VM в hex)
   ```

2. **Проверить Size:**
   ```bash
   od -A x -N 10 -t x1 /tmp/vmm_test_3605693800/test_char_simple
   # После 56 4d должно быть: 14 00 00 00 00 00 00 00 (20 в little-endian)
   ```

3. **Проверить Type:**
   ```bash
   od -A x -N 11 -t c /tmp/vmm_test_3605693800/test_char_simple | tail -1
   # Должно быть: C
   ```

4. **Извлечь первое значение "Hello":**
   ```bash
   dd if=/tmp/vmm_test_3605693800/test_char_simple skip=15 bs=1 count=5 2>/dev/null
   # Вывод: Hello
   ```

## Очистка после тестирования

```bash
# Удалить все временные тестовые файлы
rm -rf /tmp/vmm_test_*
```

---

**Теперь ты можешь:**
1. Запустить тест
2. Открыть `.vm` файл вручную
3. Проверить структуру и содержимое
4. Убедиться что CHAR данные сохранились корректно

