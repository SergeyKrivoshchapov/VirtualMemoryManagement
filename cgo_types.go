// +build dll

package main

/*
#include <stdint.h>

typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;
*/
import "C"


