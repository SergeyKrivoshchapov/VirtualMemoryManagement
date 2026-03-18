using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class TerminalCommand : ITerminalCommand
    {
        public string Title { get; }

        private readonly IArgumentParser? _argumentParser;
        private readonly ICommand<string[]> _command;

        public event EventHandler CanExecuteChanged { add { } remove { } }

        public TerminalCommand(string title, 
            IArgumentParser? argumentParser, ICommand<string[]> command)
        {
            Title = title;

            _argumentParser = argumentParser;
            _command = command;
        }

        public bool CanExecute() => true;

        public bool TryParseArguments(string arguments, out string[] parseArguments)
        {
            if (_argumentParser == null)
            {
                parseArguments = [];

                if (arguments == null || arguments == string.Empty)
                {
                    return true;
                }
                else
                {
                    return false;
                }
            }

            string[] result;

            if (_argumentParser.TryParse(arguments, out result))
            {
                parseArguments = result;

                return true;
            }

            parseArguments = [];

            return false;
        }

        public void Execute(string[] arguments)
        {
            _command.Execute(arguments);
        }
    }
}
