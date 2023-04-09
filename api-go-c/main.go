package main

/*
#include <stdio.h>

int sum(int a, int b) // 덧셈 함수 작성
{
	return a + b;
}

void hello() // Hello, world! 출력 함수 작성
{
	printf("Hello, world!\n");
}
*/
import "C"

import "fmt"

func main() {
	var a, b int = 1, 2
	r := C.sum(C.int(a), C.int(b)) // C 언어 함수 sum 호출
	fmt.Println(r)                 // 3

	C.hello() // Hello, world!
}
