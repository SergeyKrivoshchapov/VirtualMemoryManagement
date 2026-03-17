using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class TerminalInputer : ITerminalInputer
    {
        private bool _isTerminalActive = false;
        private Dictionary<string, List<ITerminalCommand>> _commands = [];
        private ITerminal _terminal;
        private ICommandParser _commandParser;
        private IErrorSender _errorSender;

        private string Title { get; set {
                ArgumentException.ThrowIfNullOrEmpty(value, nameof(Title));
                
                field = value;
            } 
        }

        public TerminalInputer(string title, 
            ITerminal terminal, ICommandParser commandParser, IErrorSender errorSender)
        {
            Title = title;
            _terminal = terminal;
            _commandParser = commandParser;
            _errorSender = errorSender;
        }

        public void RegisterCommand(ITerminalCommand command)
        {
            if (_commands.ContainsKey(command.Title))
            {
                // Need check to dublicate command

                _commands[command.Title].Add(command);
            }
            else
            {
                _commands[command.Title] = new List<ITerminalCommand> { command };
            }
        }

        public void Show()
        {
            _isTerminalActive = true;

            while (_isTerminalActive)
            {
                _terminal.Write($"{Title}>");
                string? input = _terminal.Read();

                if (input == null || input == string.Empty)
                {
                    _errorSender.NotifyOnError("Input cannot be empty.");
                    continue;
                }

                (string command, string arguments) = _commandParser.Parse(input);

                if (!_commands.ContainsKey(command))
                {
                    _errorSender.NotifyOnError("Command name is incorrect.");
                    continue;
                }

                int countOfReflectionTryParse = 0;
                ITerminalCommand comm = _commands[command][0];
                string[] parseCommands = [];

                foreach (var a in _commands[command])
                {
                    if (a.TryParseArguments(arguments, out parseCommands))
                    {
                        countOfReflectionTryParse++;
                        comm = a;
                    }

                    if (countOfReflectionTryParse > 1)
                    {
                        break;
                    }
                }

                if (countOfReflectionTryParse == 0)
                {
                    _errorSender.NotifyOnError("Arguments is incorrect.");
                    continue;
                }

                if (countOfReflectionTryParse > 1)
                {
                    _errorSender.NotifyOnError("Incorect command syntaxys. Some command can be call.");
                    continue;
                }

                comm.Execute(parseCommands);
            }
        }

        public void Close()
        {
            _isTerminalActive = false;
        }
    }
}
