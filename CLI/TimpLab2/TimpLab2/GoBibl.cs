using System;
using System.Runtime.InteropServices;

/// <summary>
/// Нативные функции из библиотеки VM (написана на Go).
/// </summary>
internal static class VmNative
{
    // Имя библиотеки без расширения — платформа сама добавит .dll, .so или .dylib.
    private const string DllName = "vmm";

    [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
    public static extern int VMCreate(string filename, int size, string typ, int stringLength);

    [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
    public static extern int VMOpen(string filename);

    [DllImport(DllName, CallingConvention = CallingConvention.Cdecl)]
    public static extern int VMClose(int handle);

    [DllImport(DllName, CallingConvention = CallingConvention.Cdecl)]
    public static extern Result VMRead(int handle, int index);

    [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
    public static extern int VMWrite(int handle, int index, string value);

    // Параметр filename может быть null — в Go это будет соответствовать нулевому указателю.
    [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
    public static extern Result VMHelp(string filename);
}

/// <summary>
/// Структура, возвращаемая функциями VMRead и VMHelp.
/// Соответствует C.Result из Go.
/// </summary>
[StructLayout(LayoutKind.Sequential)]
public struct Result
{
    /// <summary>Признак успеха (ненулевое значение — успех).</summary>
    public int success;

    /// <summary>Код ошибки, если success == 0.</summary>
    public int error_code;

    /// <summary>Буфер данных (до 256 байт). Для строк используется UTF-8, возможен завершающий ноль.</summary>
    [MarshalAs(UnmanagedType.ByValArray, SizeConst = 256)]
    public byte[] data;
}

/// <summary>
/// Методы расширения для удобной работы с Result.
/// </summary>
public static class ResultExtensions
{
    /// <summary>Проверяет, успешно ли завершилась операция.</summary>
    public static bool IsSuccess(this Result result) => result.success != 0;

    /// <summary>Извлекает строку из поля data (ожидается UTF-8, обрезается по первому нулю).</summary>
    public static string GetDataAsString(this Result result)
    {
        if (result.data == null) return string.Empty;
        int length = Array.IndexOf(result.data, (byte)0);
        if (length < 0) length = result.data.Length;
        return System.Text.Encoding.UTF8.GetString(result.data, 0, length);
    }
}