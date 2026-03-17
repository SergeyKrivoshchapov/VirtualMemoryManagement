using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Runtime.InteropServices;
using System.Text;

namespace TimpLaba2_VirtualMemory.Models
{
    public class VirtualMemoryMock : IVirtualMemmoryFileWorker
    {
        private const string DllName = "vmm.dll";

        [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        public static extern Result VMCreate(string filename, int size, string typ, int? stringLength);

        [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        public static extern Result VMOpen(string filename);

        [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        public static extern Result VMClose(int handle);



        private static IVirtualMemmoryValueWorker? _openFile = null;

        private static int? _openFileHandle = null;



        private VirtualMemoryMock() { }

        private static readonly Lazy<VirtualMemoryMock> _instance =
            new Lazy<VirtualMemoryMock>(() => new VirtualMemoryMock());

        public static VirtualMemoryMock Instance => _instance.Value;

        public void CreateFile(string fileName, VMFileType valueType)
        {
            Result result = VMCreate(fileName, 10001, valueType.StringFileType, valueType.TypeLength);

            if (!result.IsSuccess())
            {
                throw new Exception($"Error {result.error_code}: {result.GetDataAsString()}");
            }
        }

        public IVirtualMemmoryValueWorker OpenFile(string fileName)
        {
            if (_openFile != null)
            {
                _openFile.Dispose();
                CloseFile();
            }

            Result result = VMOpen(fileName);
            
            if (!result.IsSuccess())
            {
                throw new Exception($"Error {result.error_code}: {result.GetDataAsString()}");
            }

            _openFileHandle = int.Parse(result.GetDataAsString());
            _openFile = new VMValueMock((int)_openFileHandle);

            return _openFile;
        }

        public void CloseFile()
        {
            if (_openFileHandle != null)
            {
                VMClose((int)_openFileHandle);
                _openFile = null;
                _openFileHandle = null;
            }
        }
    }

    public class VMValueMock : IVirtualMemmoryValueWorker
    {
        private const string DllName = "vmm.dll";

        [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        public static extern Result VMRead(int handle, int index);

        [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        public static extern Result VMWrite(int handle, int index, string value);

        private bool IsThreadClosed = false;

        private int _fileHandle;

        public VMValueMock(int fileHandle)
        {
            _fileHandle = fileHandle;
        }

        public void WriteValue(int index, string value)
        {
            if (IsThreadClosed)
            {
                throw new ObjectDisposedException(typeof(VMValueMock).FullName, 
                    "Cannot write to a closed file.");
            }

            Result result = VMWrite(_fileHandle, index, value);

            if (!result.IsSuccess())
            {
                throw new Exception($"Error {result.error_code}: {result.GetDataAsString()}");
            }
        }

        public string ReadValue(int index)
        {
            if (IsThreadClosed)
            {
                throw new ObjectDisposedException(typeof(VMValueMock).FullName, 
                    "Cannot read to a closed file.");
            }

            Result result = VMRead(_fileHandle, index);

            if (!result.IsSuccess())
            {
                throw new Exception($"Error {result.error_code}: {result.GetDataAsString()}");
            }

            return result.GetDataAsString();
        }

        public void Dispose()
        {
            IsThreadClosed = true;
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
