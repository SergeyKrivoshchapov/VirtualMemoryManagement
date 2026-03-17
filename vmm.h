#ifndef VMM_H
#define VMM_H

#include <stdint.h>

/* Result structure for API responses */
typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;

/* Virtual Memory Manager Functions */

/* Create a new virtual array
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMCreate(const char* filename, int size, const char* typ, int stringLength);

/* Open an existing virtual array file
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMOpen(const char* filename);

/* Close a virtual array
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMClose(int handle);

/* Read a value from the array
   Returns: Result struct with data and error code */
extern Result __cdecl VMRead(int handle, int index);

/* Write a value to the array
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMWrite(int handle, int index, const char* value);

/* Get help information
   Returns: Result struct with help text */
extern Result __cdecl VMHelp(const char* filename);

#endif /* VMM_H */
