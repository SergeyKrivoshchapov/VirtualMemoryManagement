using System;
using System.Collections.Generic;
using System.Text;

namespace TerminalCommandMenu
{
    public static class FormatParser
    {
        private static readonly HashSet<char> Separators = new HashSet<char>
        {
            '(', ')', '[', ']', '{', '<', '>',
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
                    else if (spec == 'w')
                    {
                        fi++;

                        char? nextDelimiter = GetNextDelimiter(format, fi);

                        string arg = ReadWordArgument(
                            input,
                            ref ii,
                            nextDelimiter
                        );

                        result.Add(arg);
                    }
                    else if (spec == 'd')
                    {
                        fi++;

                        char? nextDelimiter = GetNextDelimiter(format, fi);

                        string arg = ReadNumberArgument(
                            input,
                            ref ii,
                            nextDelimiter
                        );

                        result.Add(arg);
                    }
                    else if (spec == 'q')
                    {
                        fi++;

                        string arg = ReadQuotedArgument(
                            input,
                            ref ii
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

                    if (s != 'a' && s != 's' && s != 'w' && s != 'd' && s != 'q')
                        throw new Exception("Invalid specifier");

                    bool isArg =
                        s == 'a' ||
                        s == 'w' ||
                        s == 'd' ||
                        s == 'q';

                    if (isArg && lastWasArg)
                        throw new Exception("Cannot use two arguments in a row");

                    lastWasArg = isArg;

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

        private static string ReadWordArgument(
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

            while (index < input.Length)
            {
                char c = input[index];

                if (c == delimiter)
                    break;

                if (IsOpen(c) || IsClose(c))
                    throw new Exception("Nested brackets not allowed in %w");

                index++;
            }

            return input.Substring(start, index - start);
        }

        private static string ReadNumberArgument(
            string input,
            ref int index,
            char? delimiter
        )
        {
            int start = index;

            if (index < input.Length && input[index] == '-')
                index++;

            bool hasDigit = false;

            while (index < input.Length)
            {
                char c = input[index];

                if (delimiter != null && c == delimiter)
                    break;

                if (!char.IsDigit(c))
                    throw new Exception("Invalid number");

                hasDigit = true;
                index++;
            }

            if (!hasDigit)
                throw new Exception("Invalid number");

            return input.Substring(start, index - start);
        }

        private static string ReadQuotedArgument(
            string input,
            ref int index
        )
        {
            if (index >= input.Length || input[index] != '"')
                throw new Exception("Expected quote");

            index++; // skip "

            int start = index;

            while (index < input.Length)
            {
                if (input[index] == '"')
                {
                    string result =
                        input.Substring(start, index - start);

                    index++; // skip closing "
                    return result;
                }

                index++;
            }

            throw new Exception("Unclosed quote");
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