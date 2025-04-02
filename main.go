package main

import (
	"fmt"
	"sync"
	"time"
)

// RING BUFFER
// ///////////////////////////////////////////////////////////////////////////
// Struct BufRing
type BufRing struct {
	size int
	data []*int
	head int
	tail int
}

// Make BufRing
func NewBufRing(size int) BufRing {
	return BufRing{size: size, data: make([]*int, size), head: 0, tail: 0}
}

var mutex = sync.RWMutex{}

// AddElement to BufRing
func (b *BufRing) Add(element int) {
	b.data[b.tail] = &element
	b.tail = (b.tail + 1) % b.size
	if b.tail == b.head {
		b.head = (b.head + 1) % b.size
	}
}

func NegativeNumber(input <-chan int, done chan int, b *BufRing) *BufRing {
	go func() {
		for {
			select {
			case <-done:
				return
			case element, ok := <-input:
				if !ok {
					return
				}
				if element > 0 {
					mutex.Lock()
					fmt.Printf("Элемент %v был добавлен в буфер\n", element)
					b.Add(element)
					mutex.Unlock()
					time.Sleep(time.Second)
				}
			}
		}
	}()
	return b
}

func Filtr(b *BufRing) <-chan int {
	output := make(chan int, len(b.data))
	go func() {
		defer close(output)
		for {
			time.Sleep(time.Second * time.Duration(timeInteraval))
			mutex.Lock()
			if b.data[b.head] == nil {
				fmt.Println("Буффер пуст")
				fmt.Println(b.data)
				return
			}
			element := int(*b.data[b.head])
			// fmt.Printf("Элемент равен %d\n", element)
			if element%3 == 0 && element != 0 {
				output <- element
				b.data[b.head] = nil
				b.head = (b.head + 1) % b.size
			} else {
				b.data[b.head] = nil
				b.head = (b.head + 1) % b.size
			}
			mutex.Unlock()
		}

	}()
	return output
}

const bufferSize int = 5
const timeInteraval int = 1

func main() {
	buf := NewBufRing(bufferSize)
	done := make(chan int)
	start := []int{}
	var value int
ENTER:
	for {
		fmt.Println("Enter the value")
		_, err := fmt.Scanln(&value)
		if err != nil {
			fmt.Println("Value is not correct")
			break
		}
		start = append(start, value)
		var action string
		fmt.Println("Coutinue enter: Press Enter/Stop enter: Enter Stop")
		fmt.Scanf("%s", &action)
		switch action {
		case "":
			continue
		case "Stop":
			break ENTER
		case "stop":
			break ENTER
		default:
			fmt.Println("Unknow command")
			break ENTER
		}
	}
	//Формирование исходного канала
	init := func(array []int) <-chan int {
		output := make(chan int, len(array))
		go func() {
			defer close(output)
			for _, element := range array {
				select {
				case <-done:
					return
				case output <- element:
				}
			}
		}()
		return output
	}

	input := init(start)
	pipline := Filtr(NegativeNumber(input, done, &buf))
	for v := range pipline {
		fmt.Println(v)
	}
}
