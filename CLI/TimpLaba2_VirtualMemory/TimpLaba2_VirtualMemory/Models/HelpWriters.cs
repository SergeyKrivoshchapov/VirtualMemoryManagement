using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using System.Text;

namespace TimpLaba2_VirtualMemory.Models
{
    public interface IHelpWriter<T>
    {
        void WriteHelp(T args);
    }

    public class HelpWriter : IHelpWriter<string[]>
    {
        private const string DllName = "vmm.dll";

        [DllImport(DllName, CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        public static extern Result VMHelp(string filename, string helpText);

        public void WriteHelp(string[] args)
        {
            if (args.Length != 2)
            {
                throw new Exception("Incorrect args count");
            }

            Result result = VMHelp(args[0], args[1]);

            if (!result.IsSuccess())
            {
                throw new Exception($"Error {result.error_code}: {result.data}");
            }
        }
    }
}
