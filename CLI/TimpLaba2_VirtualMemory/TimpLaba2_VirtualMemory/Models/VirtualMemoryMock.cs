using System;
using System.Collections.Generic;
using System.Text;

namespace TimpLaba2_VirtualMemory.Models
{
    public class VirtualMemoryMock : IVirtualMemmoryFileWorker
    {
        public void CloseFile()
        {
            Console.WriteLine("CloseFile");
        }

        public void CreateFile(string fileName, VMFileType valueType)
        {
            Console.WriteLine("CreateFile");
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
}
