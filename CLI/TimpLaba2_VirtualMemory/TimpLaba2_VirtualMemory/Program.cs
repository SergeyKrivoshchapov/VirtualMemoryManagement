using System;
using TerminalCommandMenu;
using TimpLaba2_VirtualMemory.Models;
using TimpLaba2_VirtualMemory.Presenters;
using TimpLaba2_VirtualMemory.Views;

namespace MainProgrma
{
    class Program
    {
        static void Main(string[] args)
        {

            IVMTerminalPresenter<string[]> presenter = new VMTerminalPresenter(null,
                new VirtualMemoryMock());
            var view = new TerminalView(presenter, "$root", new Terminal(),
                new CommandParser(), new ErrorSender(new Terminal()));

            presenter.View = view;

            view.Show();
        }
    }
}