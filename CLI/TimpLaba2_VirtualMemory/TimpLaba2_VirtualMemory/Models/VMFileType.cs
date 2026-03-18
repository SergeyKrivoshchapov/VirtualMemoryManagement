using System;
using System.Collections.Generic;
using System.Text;

namespace TimpLaba2_VirtualMemory.Models
{
    public class VMFileType
    {
        private readonly FileType _fileType;

        public string StringFileType { get
            {
                switch (_fileType)
                {
                    case FileType.integer:
                        return "int";
                    case FileType.chareble:
                        return "char";
                    case FileType.varchar:
                        return "varchar";
                    default:
                        throw new InvalidOperationException("Invalid file type.");
                }
            } }

        public int? TypeLength { get; }

        public enum FileType
        {
            integer, chareble, varchar
        }

        static public bool TryConvertStringToFileType(string fileType, out FileType? outFileType)
        {
            switch (fileType)
            {
                case "int":
                    outFileType = FileType.integer;
                    return true;
                case "char":
                    outFileType = FileType.chareble;
                    return true;
                case "varchar":
                    outFileType = FileType.varchar;
                    return true;
                default:
                    outFileType = null;
                    return false;
            }
        }

        public VMFileType(FileType filetype, int? typeLength)
        {
            _fileType = filetype;
            TypeLength = typeLength;
        }

    }
}
