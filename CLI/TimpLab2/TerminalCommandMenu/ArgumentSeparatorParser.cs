using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class ArgumentSeparatorParser : IArgumentParser
    {
        private string _separator;

        public ArgumentSeparatorParser(string separator)
        {
            _separator = separator;
        }

        public string[] Parse(string arguments)
        {
            return arguments.Split(_separator);
        }

        public bool TryParse(string arguments, out string[] result)
        {
            try
            {
                result = arguments.Split(_separator);

                return true;
            }
            catch (Exception e)
            {
                result = [];
            }

            return false;
        }
    }
}
