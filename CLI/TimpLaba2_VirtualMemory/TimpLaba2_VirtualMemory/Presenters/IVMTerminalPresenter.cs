using System;
using System.Collections.Generic;
using System.Text;
using TimpLaba2_VirtualMemory.Views;

namespace TimpLaba2_VirtualMemory.Presenters
{
    public interface IVMTerminalPresenter<T>
    {
        ITerminalView View { set; }

        void CreateFile(T arguments);

        void OpenFile(T arguments);

        void InputValue(T arguments);

        void PrintValue(T arguments);
    }
}
