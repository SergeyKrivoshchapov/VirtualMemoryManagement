using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu;
using TerminalCommandMenu.Abstractions;
using TimpLaba2_VirtualMemory.Presenters;

namespace TimpLaba2_VirtualMemory.Views
{
    public class TerminalView : TerminalInputer, ITerminalView
    {
        protected IVMTerminalPresenter<string[]> _presenter;

        public TerminalView(IVMTerminalPresenter<string[]> presenter, 
            string title, ITerminal terminal, ICommandParser commandParser,  IErrorSender errorSender)
            : base(title, terminal, commandParser, errorSender)
        {
            _presenter = presenter;

            ICommand<string[]> createFileCommand 
                = new Command((string[] args) => _presenter.CreateFile(args));
            ICommand<string[]> openFileCommand
                = new Command((string[] args) => _presenter.OpenFile(args));
            ICommand<string[]> inputValueCommand
                = new Command((string[] args) => _presenter.InputValue(args));
            ICommand<string[]> printValueCommand
                = new Command((string[] args) => _presenter.PrintValue(args));
            ICommand<string[]> helpCommand
                = new Command((string[] args) => PrintHelp());
            ICommand<string[]> exitTerminalCommand
                = new Command((string[] args) => Close());

            ITerminalCommand createInt = new TerminalCommand("Create",
                new ArgumentFormatParser("%a(%w)"), createFileCommand);
            ITerminalCommand createChar = new TerminalCommand("Create",
                new ArgumentFormatParser("%a(%w(%w))"), createFileCommand);
            ITerminalCommand open = new TerminalCommand("Open",
                new ArgumentFormatParser("%a"), openFileCommand);
            ITerminalCommand input = new TerminalCommand("Input",
                new ArgumentFormatParser("(%a,%s%a)"), inputValueCommand);
            ITerminalCommand print = new TerminalCommand("Input",
                new ArgumentFormatParser("(%a)"), printValueCommand);
            ITerminalCommand help = new TerminalCommand("Help",
                null, helpCommand);
            ITerminalCommand exit = new TerminalCommand("Exit",
                null, exitTerminalCommand);

            RegisterCommand(createInt);
            RegisterCommand(createChar);
            RegisterCommand(open);
            RegisterCommand(input);
            RegisterCommand(print);
            RegisterCommand(help);
            RegisterCommand(exit);
        }

        public void DisplayMessage(string message)
        {
            _terminal.Write(message);
        }

        public void DisplayError(string message)
        {
            _errorSender.NotifyOnError(message);
        }

        private void PrintHelp()
        {
             _terminal.Write(
                "Available commands:\n" +
                "  Create <fileName>(int | char(<len>) | varchar(<maxLen>))\n" +
                "      — создает все необходимые структуры в памяти и файлы на диске.\n" +
                "\n" +
                "  Open <fileName>\n" +
                "      — открывает указанный файл и связанные с ним файлы в режиме rw,\n" +
                "        создает структуры в памяти и считывает заданное количество\n" +
                "        страниц (>= 3), устанавливая абсолютный номер, статус и время записи.\n" +
                "\n" +
                "  Input(<index>, <value>)\n" +
                "      — записывает значение в элемент массива с номером <index>.\n" +
                "        Строковое значение нужно заключать в кавычки.\n" +
                "\n" +
                "  Print(<index>)\n" +
                "      — выводит на экран значение элемента массива с номером <index>.\n" +
                "\n" +
                "  Help [fileName]\n" +
                "      — выводит список команд на экран или в указанный файл.\n" +
                "\n" +
                "  Exit\n" +
                "      — закрывает все файлы и завершает программу.\n" +
                "        Файлы виртуального массива не уничтожаются автоматически.\n"
            );
        }
    }
}
