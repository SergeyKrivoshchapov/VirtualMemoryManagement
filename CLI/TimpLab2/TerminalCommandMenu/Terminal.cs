using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class Terminal : ITerminal
    {
        public void Write(string text)
        {
            Console.Write(text);
        }

        public string? Read()
        {
            return Console.ReadLine();
        }
    }
}
