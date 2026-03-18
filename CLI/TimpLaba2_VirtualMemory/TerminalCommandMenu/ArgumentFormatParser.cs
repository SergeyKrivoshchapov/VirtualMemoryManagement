using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class ArgumentFormatParser : IArgumentParser
    {
        private readonly string _argumentFormat;

        public ArgumentFormatParser(string format)
        {
            _argumentFormat = format;
        }

        public string[] Parse(string arguments)
        {
            return FormatParser.Parse(_argumentFormat, arguments);
        }

        public bool TryParse(string arguments, out string[] result)
        {
            try
            {
                result = FormatParser.Parse(_argumentFormat, arguments);
                return true;

            }
            catch (Exception e)
            {
                result = [];
                return false;
            }
        }
    }
}
