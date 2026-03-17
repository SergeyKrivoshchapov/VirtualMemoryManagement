using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using System.Text;

namespace TimpLaba2_VirtualMemory.Models
{
    public class VirtualMemoryMock : IVirtualMemmoryFileWorker
    {
        public void CreateFile(string fileName, VMFileType valueType)
        {
            Console.WriteLine("CreateFile");
        }

        public void CloseFile()
        {
            Console.WriteLine("CloseFile");
        }

        public IVirtualMemmoryValueWorker OpenFile(string fileName)
        {
            Console.WriteLine("OpenFile");

            return new VMValueMock();
        }
    }

    public class VMValueMock : IVirtualMemmoryValueWorker
    {
        public void Dispose()
        {

        }

        public string ReadValue(int index)
        {
            Console.WriteLine("ReadValue");

            return "";
        }

        public void WriteValue(int index, string value)
        {
            Console.WriteLine("WriteValue");
        }
    }

    [StructLayout(LayoutKind.Sequential)]
    public struct Result
    {
        public int success;

        public int error_code;

        [MarshalAs(UnmanagedType.ByValArray, SizeConst = 256)]
        public byte[] data;
    }

    public static class ResultExtensions
    {
        public static bool IsSuccess(this Result result) => result.success != 0;

        public static string GetDataAsString(this Result result)
        {
            if (result.data == null) return string.Empty;
            int length = Array.IndexOf(result.data, (byte)0);
            if (length < 0) length = result.data.Length;
            return System.Text.Encoding.UTF8.GetString(result.data, 0, length);
        }
    }
}
