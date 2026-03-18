# Быстрая инструкция - тест VMWrite CHAR

## TL;DR

```bash
# 1. Запустить тест с сохранением файла
KEEP_TEST_FILES=1 go test -v -run TestVMWriteCharSimple ./api

# 2. Скопировать путь из вывода, например:
# File location: /tmp/vmm_test_3605693800/test_char_simple

# 3. Посмотреть файл
hexdump -C /tmp/vmm_test_3605693800/test_char_simple | head -30

# 4. Или извлечь текст
strings /tmp/vmm_test_3605693800/test_char_simple

# 5. Очистить
rm -rf /tmp/vmm_test_*
```

## Что записывает тест

Массив из 20 элементов CHAR с длиной строки 30 байт:
```
[0] = "Hello"
[1] = "World"
[2] = "Test String"
[3] = "1234567890"
[4] = "Short"
[5] = "This is a longer test string"
[6-19] = пусто
```

## Структура файла

```
0-1:   Сигнатура "VM"
2-9:   Size = 20 (int64)
10:    Type = 'C' (CHAR)
11-14: StringLength = 30 (int32)
15+:   Данные (1 страница = 4096 байт)
```

## Проверить в файле

```bash
# Сигнатура
head -c 2 test_char_simple | od -t x1
# Выход: 56 4d (VM в hex)

# Все значения
strings test_char_simple
# Выход содержит все записанные строки

# Размер
ls -lh test_char_simple
# ~4.1K (размер 1 страницы с выравниванием)
```

## Для Windows PowerShell

```powershell
$env:KEEP_TEST_FILES="1"
go test -v -run TestVMWriteCharSimple ./api

# После теста в выводе будет путь к файлу
# Можешь открыть его в любом hex-редакторе
```

---

**Готово!** Файл остается после теста, ты можешь посмотреть его содержимое вручную.

