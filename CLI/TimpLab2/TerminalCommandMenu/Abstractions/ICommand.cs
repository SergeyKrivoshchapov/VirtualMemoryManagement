namespace TerminalCommandMenu.Abstractions
{
    public interface ICommand<in T> 
    {
        void Execute(T arguments);
    }
}
