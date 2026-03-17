using System;
using System.Collections.Generic;
using System.Text;
using TerminalCommandMenu.Abstractions;

namespace TerminalCommandMenu
{
    public class ErrorSender : IErrorSender
    {
        private IWriteble _writer;

        public ErrorSender(IWriteble writer)
        {
            _writer = writer;
        }

        public void NotifyOnError(string message)
        {
            _writer.Write(message + "\n");
        }
    }
}
