using System;
using System.Collections.Generic;
using System.Text;

namespace TerminalCommandMenu.Abstractions
{
    public interface ICommandParser
    {
        (string Command, string arguments) Parse(string input);
    }

    public interface ICommandsParser
    {
        List<(string Command, string arguments)> Parse(string input);
    }
}
