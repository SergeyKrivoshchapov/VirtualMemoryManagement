// +build !dll

package main

import (
	"VirtualMemoryManagement/api"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var activeHandle int = -1

	fmt.Print("VM> ")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Print("VM> ")
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "create":
			handleCreate(parts)
			activeHandle = api.GetHandle()

		case "open":
			if len(parts) < 2 {
				fmt.Println("Usage: Open <filename>")
			} else {
				result := api.VMOpen(parts[1])
				if result.IsSuccess() {
					activeHandle = api.GetHandle()
					fmt.Println("Opened:", result.String())
				} else {
					fmt.Println("Error:", result.String())
				}
			}

		case "close":
			if len(parts) < 2 {
				if activeHandle == -1 {
					fmt.Println("No active handle")
				} else {
					result := api.VMClose(activeHandle)
					if result.IsSuccess() {
						fmt.Println("Closed")
						activeHandle = -1
					} else {
						fmt.Println("Error:", result.String())
					}
				}
			} else {
				handle, _ := strconv.Atoi(parts[1])
				result := api.VMClose(handle)
				if result.IsSuccess() {
					fmt.Println("Closed")
					if activeHandle == handle {
						activeHandle = -1
					}
				} else {
					fmt.Println("Error:", result.String())
				}
			}

		case "read", "print":
			if len(parts) < 2 {
				fmt.Println("Usage: Read <index>")
			} else {
				if activeHandle == -1 {
					fmt.Println("No active handle")
				} else {
					index, err := strconv.Atoi(parts[1])
					if err != nil {
						fmt.Println("Invalid index:", parts[1])
					} else {
						result := api.VMRead(activeHandle, index)
						if result.IsSuccess() {
							fmt.Println("Value:", result.String())
						} else {
							fmt.Println("Error:", result.String())
						}
					}
				}
			}

		case "write", "input":
			if activeHandle == -1 {
				fmt.Println("No active handle")
			} else if len(parts) < 3 {
				fmt.Println("Usage: Write <index> <value>")
			} else {
				index, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("Invalid index:", parts[1])
				} else {
					value := extractStringValue(line, parts[0])
					result := api.VMWrite(activeHandle, index, value)
					if result.IsSuccess() {
						fmt.Println("Written")
					} else {
						fmt.Println("Error:", result.String())
					}
				}
			}

		case "stats":
			if activeHandle == -1 {
				fmt.Println("No active handle")
			} else {
				result := api.VMStats(activeHandle)
				if result.IsSuccess() {
					fmt.Println(result.String())
				} else {
					fmt.Println("Error:", result.String())
				}
			}

		case "help":
			var filename string
			if len(parts) > 1 {
				filename = parts[1]
			}
			result := api.VMHelp(filename)
			fmt.Println(result.String())

		case "exit", "quit":
			handles := api.GetAllHandles()
			for _, h := range handles {
				api.VMClose(h)
			}
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command:", cmd)
		}

		fmt.Print("VM> ")
	}
}

func handleCreate(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Usage: Create <filename> <type> [<length>]")
		fmt.Println("Types: int, char(length), varchar(maxlength)")
		return
	}

	filename := parts[1]
	typeStr := parts[2]

	var arrayType string
	var stringLength int

	if strings.Contains(typeStr, "int") {
		arrayType = "int"
	} else if strings.HasPrefix(typeStr, "char") {
		arrayType = "char"
		re := regexp.MustCompile(`char\((\d+)\)`)
		matches := re.FindStringSubmatch(typeStr)
		if len(matches) > 1 {
			stringLength, _ = strconv.Atoi(matches[1])
		} else if len(parts) > 3 {
			stringLength, _ = strconv.Atoi(parts[3])
		}
	} else if strings.HasPrefix(typeStr, "varchar") {
		arrayType = "varchar"
		re := regexp.MustCompile(`varchar\((\d+)\)`)
		matches := re.FindStringSubmatch(typeStr)
		if len(matches) > 1 {
			stringLength, _ = strconv.Atoi(matches[1])
		} else if len(parts) > 3 {
			stringLength, _ = strconv.Atoi(parts[3])
		}
	}

	size := 10000
	if len(parts) > 4 {
		size, _ = strconv.Atoi(parts[4])
	}

	result := api.VMCreate(filename, size, arrayType, stringLength)
	if result.IsSuccess() {
		fmt.Println("Created:", result.String())
	} else {
		fmt.Println("Error:", result.String())
	}
}

func extractStringValue(line string, cmdName string) string {
	index := strings.Index(line, cmdName)
	if index == -1 {
		return ""
	}

	remaining := strings.TrimLeft(line[index+len(cmdName):], " \t")

	parts := strings.SplitN(remaining, " ", 2)
	if len(parts) < 2 {
		return ""
	}

	valueStr := strings.TrimSpace(parts[1])

	if strings.HasPrefix(valueStr, "\"") && strings.HasSuffix(valueStr, "\"") {
		return valueStr[1 : len(valueStr)-1]
	}

	return valueStr
}
