namespace TerminalCommandMenu.Abstractions
{
    public interface ITerminalCommand
    {
        string Title { get; }

        event EventHandler CanExecuteChanged { add { } remove { } }

        bool TryParseArguments(string arguments, out string[] parseArguments);

        bool CanExecute();

        void Execute(string[] arguments);

    }
}
