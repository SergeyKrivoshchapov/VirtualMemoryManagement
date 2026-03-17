namespace TerminalCommandMenu.Abstractions
{
    public interface IErrorSender
    {
        void NotifyOnError(string message);
    }
}
