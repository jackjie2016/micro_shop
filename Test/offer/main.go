package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var flag1 bool
var flag2 bool

func main() {
	var s string
	fmt.Scanln(&s)
	s1 := strings.Split(s, ",")
	if len(s1) > 14 {
		log.Fatal("数字不能超过14")
	}
	//使用map基数 判断同一个数字是否超过4此
	m := make(map[string]int)
	for _, k := range s1 {
		if _, ok := m[k]; ok {
			m[k] += 1
			if m[k] > 4 {
				log.Fatal("同一个数字不能超过4次")
			}
		} else {
			m[k] = 1
		}
	}
	flag1 = false
	//判断是否至少一个k的次数=2
	for _, mv := range m {
		if mv == 2 {
			flag1 = true
			//如果有一对的，那去吧切片中删除这个
			//for k, v := range s1 {
			//	if v == mk {
			//			//		s1 = append(s1[:k-1], s1[k+1:len(s1)]...)
			//			//	}
			//}
		}

	}

	if !flag1 {
		log.Fatal("NO")
	} else {

		if s1[0] == s1[1] {
			if s1[1] == s1[2] {
				s1[2] = "0"
			}
			s1[0] = "0"
			s1[1] = "0"
		} else {
			a, _ := strconv.Atoi(s1[0])
			b, _ := strconv.Atoi(s1[1])
			c, _ := strconv.Atoi(s1[2])
			if a+1 != b || b+1 != c {
				log.Fatal("NO")
			}
			s1[0] = "0"
			s1[1] = "0"
			s1[2] = "0"
		}

		//不等于0 说明对子出现

		//fmt.Println(s1)
		for j := 0; j < 3; j++ {
			app1(s1, 2+3*j)
		}

		if s1[12] != s1[13] {
			log.Fatal("NO")
		} else {
			s1[12] = "0"
			s1[13] = "0"
		}

		for _, k := range s1 {
			if k != "0" {
				log.Fatal("NO")
			}
		}
		log.Fatal("YES")
		//不等于0 说明对子出现
		//if s1[2] != "0" {
		//	if s1[2] == s1[3] {
		//		if s1[3] != s1[4] {
		//			log.Fatal("NO")
		//		}
		//	} else {
		//		a, _ := strconv.Atoi(s1[2])
		//		b, _ := strconv.Atoi(s1[3])
		//		c, _ := strconv.Atoi(s1[4])
		//		if a+1 != b || b+1 != c {
		//			log.Fatal("NO")
		//		}
		//		s1[2] = "0"
		//		s1[3] = "0"
		//		s1[4] = "0"
		//	}
		//} else {
		//	if s1[2] == s1[3] {
		//		if s1[3] == s1[4] {
		//			s1[4] = "0"
		//		}
		//		s1[2] = "0"
		//		s1[3] = "0"
		//	} else {
		//		a, _ := strconv.Atoi(s1[2])
		//		b, _ := strconv.Atoi(s1[3])
		//		c, _ := strconv.Atoi(s1[4])
		//		if a+1 != b || b+1 != c {
		//			log.Fatal("NO")
		//		}
		//		s1[2] = "0"
		//		s1[3] = "0"
		//		s1[4] = "0"
		//	}
		//}

	}
}
func app1(s1 []string, i int) {
	if s1[i] != "0" {
		if s1[i] == s1[i+1] {
			if s1[i+1] != s1[i+2] {
				log.Fatal("NO")
			}
		} else {
			a, _ := strconv.Atoi(s1[i])
			b, _ := strconv.Atoi(s1[i+1])
			c, _ := strconv.Atoi(s1[i+2])
			if a+1 != b || b+1 != c {
				log.Fatal("NO")
			}
			s1[i] = "0"
			s1[i+1] = "0"
			s1[i+2] = "0"
		}
	} else {
		i = i + 1
		if s1[i] == s1[i+1] {
			if s1[i+1] == s1[i+2] {
				s1[i+2] = "0"
			}
			s1[i] = "0"
			s1[i+1] = "0"
		} else {
			a, _ := strconv.Atoi(s1[i])
			b, _ := strconv.Atoi(s1[i+1])
			c, _ := strconv.Atoi(s1[i+2])
			if a+1 != b || b+1 != c {
				log.Fatal("NO")
			}
			s1[i] = "0"
			s1[i+1] = "0"
			s1[i+2] = "0"
		}
	}
}

//func app2(s []string) {
//
//}
