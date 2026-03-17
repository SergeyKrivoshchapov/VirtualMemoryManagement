using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class Command : ICommand<string[]>
    {
        Action<string[]> _action;

        public Command(Action<string[]> action)
        {
            _action = action;
        }

        public void Execute(string[] arguments)
        {
            _action(arguments);
        }
    }
}
