using System;
using System.Collections.Generic;
using System.Text;

namespace TerminalCommandMenu
{
    using System;
    using System.Collections.Generic;

    public static class FormatParser
    {
        private static readonly HashSet<char> Separators = new HashSet<char>
            {
                '(', ')', '[', ']', '{', '}', '<', '>',
                '/', '\\', '|', ',', '.', ':', ';',
                '-', '+', '*', '=', '_'
            };

        public static string[] Parse(string format, string input)
        {
            ValidateFormat(format);

            List<string> result = new List<string>();

            int fi = 0;
            int ii = 0;

            while (fi < format.Length)
            {
                char f = format[fi];

                if (f == '%')
                {
                    fi++;

                    if (fi >= format.Length)
                        throw new Exception("Invalid format");

                    char spec = format[fi];

                    if (spec == 'a')
                    {
                        fi++;

                        char? nextDelimiter = GetNextDelimiter(format, fi);

                        string arg = ReadArgument(
                            input,
                            ref ii,
                            nextDelimiter
                        );

                        result.Add(arg);
                    }
                    else if (spec == 's')
                    {
                        if (ii >= input.Length || input[ii] != ' ')
                            throw new Exception("Expected space");

                        ii++;
                        fi++;
                    }
                    else
                    {
                        throw new Exception("Unknown specifier");
                    }
                }
                else
                {
                    if (ii >= input.Length || input[ii] != f)
                        throw new Exception("Delimiter mismatch");

                    ii++;
                    fi++;
                }
            }

            if (ii != input.Length)
                throw new Exception("Extra characters");

            return result.ToArray();
        }

        private static void ValidateFormat(string format)
        {
            bool lastWasArg = false;

            for (int i = 0; i < format.Length; i++)
            {
                char c = format[i];

                if (c == '%')
                {
                    if (i + 1 >= format.Length)
                        throw new Exception("Invalid %");

                    char s = format[i + 1];

                    if (s != 'a' && s != 's')
                        throw new Exception("Invalid specifier");

                    if (s == 'a' && lastWasArg)
                        throw new Exception("Cannot use %a%a");

                    lastWasArg = s == 'a';

                    i++;
                }
                else
                {
                    if (!Separators.Contains(c))
                        throw new Exception("Invalid character in format");

                    lastWasArg = false;
                }
            }
        }

        private static char? GetNextDelimiter(string format, int pos)
        {
            for (int i = pos; i < format.Length; i++)
            {
                if (format[i] == '%')
                    i++;
                else
                    return format[i];
            }

            return null;
        }

        private static string ReadArgument(
            string input,
            ref int index,
            char? delimiter
        )
        {
            int start = index;

            if (delimiter == null)
            {
                index = input.Length;
                return input.Substring(start);
            }

            int depth = 0;

            while (index < input.Length)
            {
                char c = input[index];

                if (c == delimiter && depth == 0)
                    break;

                if (IsOpen(c))
                {
                    depth++;
                }
                else if (IsClose(c))
                {
                    depth--;
                    if (depth < 0)
                        throw new Exception("Unbalanced brackets");
                }

                index++;
            }

            return input.Substring(start, index - start);
        }

        private static bool IsOpen(char c)
        {
            return c == '(' || c == '[' || c == '{' || c == '<';
        }

        private static bool IsClose(char c)
        {
            return c == ')' || c == ']' || c == '}' || c == '>';
        }
    }
}
