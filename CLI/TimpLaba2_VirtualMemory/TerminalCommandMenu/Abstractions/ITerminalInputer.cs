using System;
using System.Collections.Generic;
using System.Text;

namespace TerminalCommandMenu.Abstractions
{
    public interface ITerminalInputer
    {
        void Show();

        void Close();
    }
}
