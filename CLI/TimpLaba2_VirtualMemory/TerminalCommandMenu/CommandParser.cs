using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class CommandParser : ICommandParser
    {
        public (string Command, string arguments) Parse(string input)
        {
            ArgumentNullException.ThrowIfNullOrWhiteSpace("Input string is null or white space", input);

            string[] parseString = input.Split(' ', StringSplitOptions.RemoveEmptyEntries);

            if (parseString.Length == 0) throw new ArgumentException("Input string is incorrect", input);

            return (parseString[0],  parseString.Length > 1 
                ? string.Join(' ', parseString[1..]) : "");
        }
    }
}
