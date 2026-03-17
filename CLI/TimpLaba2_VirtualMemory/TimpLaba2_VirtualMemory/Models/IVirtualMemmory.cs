using System;
using System.Collections.Generic;
using System.Text;

namespace TimpLaba2_VirtualMemory.Models
{
    public interface IVirtualMemmoryFileWorker
    {
        void CreateFile(string fileName, VMFileType valueType);

        IVirtualMemmoryValueWorker OpenFile(string fileName);

        void CloseFile();
    }

    public interface IVirtualMemmoryValueWorker : IDisposable
    {
        void WriteValue(int index, string value);

        string ReadValue(int index);
    }
}
