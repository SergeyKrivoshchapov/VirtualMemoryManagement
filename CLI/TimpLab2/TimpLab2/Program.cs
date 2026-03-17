using System;
using TerminalCommandMenu;
using TerminalCommandMenu.Abstractions;

namespace Prog
{
    class Program
    {
        static void Main()
        {
            ITerminal terminal = new Terminal();
            ICommandParser commandParse = new CommandParser();
            IArgumentParser argumentParser = new ArgumentSeparatorParser(" ");
            IErrorSender errorSender = new ErrorSender(terminal);

            TerminalInputer terminalInputer = new("$root", terminal, commandParse, errorSender);

            ICommand<string[]> CreateCommand = new Command((string[] x) => 
            {
                VmNative.VMCreate(x[0], Convert.ToInt32(x[1]), x[2], Convert.ToInt32(x[3])); 
            });
            ITerminalCommand createComm = new TerminalCommand("Create", argumentParser, CreateCommand);
            terminalInputer.RegisterCommand(createComm);



            ICommand<string[]> OpenCommand = new Command((string[] x) =>
            {
                VmNative.VMOpen(x[0]);
            });
            ITerminalCommand openComm = new TerminalCommand("Open", argumentParser, OpenCommand);
            terminalInputer.RegisterCommand(openComm);



            ICommand<string[]> InputCommand = new Command((string[] x) =>
            {
                VmNative.VMWrite(Convert.ToInt32(x[0]), Convert.ToInt32(x[1]), x[2]);
            });
            ITerminalCommand inputComm = new TerminalCommand("Input", argumentParser, InputCommand);
            terminalInputer.RegisterCommand(inputComm);



            ICommand<string[]> PrintCommand = new Command((string[] x) =>
            {
                VmNative.VMRead(Convert.ToInt32(x[0]), Convert.ToInt32(x[1]));
            });
            ITerminalCommand printComm = new TerminalCommand("Print", argumentParser, PrintCommand);
            terminalInputer.RegisterCommand(printComm);



            ICommand<string[]> HelpCommand = new Command((string[] x) =>
            {
                VmNative.VMHelp(null);
            });
            ITerminalCommand helpComm = new TerminalCommand("Help", argumentParser, HelpCommand);
            terminalInputer.RegisterCommand(helpComm);






            ICommand<string[]> exitComm = new Command((string[] x) => { terminalInputer.Close(); });
            ITerminalCommand exitCommand = new TerminalCommand("Exit", argumentParser, exitComm);
            terminalInputer.RegisterCommand(exitCommand);

            terminalInputer.Show();
        }
    }
}