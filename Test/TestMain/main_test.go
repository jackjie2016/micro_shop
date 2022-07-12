package main

import (
	"testing"
	"fmt"
)

func testPrint1(t *testing.T)  {

	res:=Print1to20()
	fmt.Println("hey")
	if res!=200{
		t.Errorf("wrong result of Print1to20")
	}
}

func testPrint2(t *testing.T)  {

	res:=Print1to20()
	fmt.Println("hey")
	if res!=200{
		t.Errorf("wrong result of Print1to20")
	}
}

func TestAll(t *testing.T)  {
	t.Run("testPrint1",testPrint1)
	t.Run("testPrint2",testPrint2)
}
func TestMain(m *testing.M) {
	fmt.Println("Test main begin")
	m.Run()
}

func aaa(n int) int {
	for n>0 {
		n--
	}
	return n
}
func  BenchmarkALl(b *testing.B)  {
	for n:=0;n<b.N;n++ {
		aaa(n)
	}
}