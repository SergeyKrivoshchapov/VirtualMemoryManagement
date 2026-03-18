using System;
using System.Collections.Generic;
using System.Text;

namespace TerminalCommandMenu.Abstractions
{
    public interface IArgumentParser
    {
        string[] Parse(string arguments);

        bool TryParse(string arguments, out string[] result);
    }
}
