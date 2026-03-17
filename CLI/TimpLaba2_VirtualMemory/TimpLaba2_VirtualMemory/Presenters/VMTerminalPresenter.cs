using System;
using System.Collections.Generic;
using System.Security.AccessControl;
using System.Text;
using TimpLaba2_VirtualMemory.Models;
using TimpLaba2_VirtualMemory.Views;

namespace TimpLaba2_VirtualMemory.Presenters
{
    public class VMTerminalPresenter : IVMTerminalPresenter<string[]>
    {
        protected ITerminalView? _view;
        protected IVirtualMemmoryFileWorker _fileWorker;
        protected IVirtualMemmoryValueWorker? _valueWorker = null;

        public ITerminalView View { set { _view = value; } }

        public VMTerminalPresenter(ITerminalView? view, IVirtualMemmoryFileWorker fileWorker)
        {
            _view = view;
            _fileWorker = fileWorker;
        }

        public void CreateFile(string[] arguments)
        {
            if (arguments.Length < 2 || arguments.Length > 3)
            {
                _view?.DisplayError("Not enough arguments for Create command.");
                return;
            }

            VMFileType.FileType? fileType;

            if (VMFileType.TryConvertStringToFileType(arguments[1], out fileType))
            {
                if (arguments.Length == 2)
                {
                    try
                    {
                        _fileWorker.CreateFile(arguments[0], 
                            new VMFileType((VMFileType.FileType)fileType, null));
                    }
                    catch (Exception ex)
                    {
                        _view?.DisplayError(ex.Message);
                    }
                }
                else
                {
                    int typeLength;

                    if (!int.TryParse(arguments[2], out typeLength))
                    {
                        _view?.DisplayError("Invalid type length argument for Create command.");
                        return;
                    }

                    try
                    {
                        _fileWorker.CreateFile(arguments[0], 
                            new VMFileType((VMFileType.FileType)fileType, typeLength));
                    }
                    catch (Exception ex)
                    {
                        _view?.DisplayError(ex.Message);
                    }
                }
            }            
        }

        public void OpenFile(string[] arguments)
        {
            if (arguments.Length != 1)
            {
                _view?.DisplayError("Not enough arguments for Open command.");
                return;
            }

            try
            {
                _valueWorker = _fileWorker.OpenFile(arguments[0]);
            }
            catch (Exception ex)
            {
                _view?.DisplayError(ex.Message);
            }
        }

        public void InputValue(string[] arguments)
        {
            if (_valueWorker == null)
            {
                _view?.DisplayError("File doesn't open.");
                return;
            }

            if (arguments.Length != 2)
            {
                _view?.DisplayError("Not enough arguments for Input command.");
                return;
            }

            int index;
            string value;

            if (!int.TryParse(arguments[0], out index))
            {
                _view?.DisplayError("Index should be integer");
                return;
            }

            value = arguments[1];
            
            try
            {
                _valueWorker.WriteValue(index, value);
            }
            catch (Exception ex)
            {
                _view?.DisplayError(ex.Message);
            }
        }

        public void PrintValue(string[] arguments)
        {
            if (_valueWorker == null)
            {
                _view?.DisplayError("File doesn't open.");
                return;
            }

            if (arguments.Length != 1)
            {
                _view?.DisplayError("Not enough arguments for Print command.");
                return;
            }

            int index;

            if (!int.TryParse(arguments[0], out index))
            {
                _view?.DisplayError("Index should be integer");
                return;
            }

            try
            {
                _valueWorker.ReadValue(index);
            }
            catch (Exception ex)
            {
                _view?.DisplayError(ex.Message);
            }
        }
    }
}
