using System;
using System.Collections.Generic;
using System.Text;

namespace TimpLaba2_VirtualMemory.Views
{
    public interface ITerminalView
    {
        void DisplayMessage(string message);

        void DisplayError(string message);
    }
}
