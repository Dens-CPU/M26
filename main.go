package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
					log.Printf("Элемент %v был добавлен в буфер\n", element)
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
				log.Printf("Буффер пуст:")
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
	fmt.Printf("Enter values separated by spaces:")
	reader := bufio.NewReader(os.Stdin)
	values, _ := reader.ReadString('\n')
	n := 0
	for i := 0; i < len(values); i++ {
		if values[i] >= 48 && values[i] <= 57 {
			if n == 0 {
				n = int(rune(values[i]) - 48)
			} else {
				n = (n * 10) + int(rune(values[i]-48))
			}
		} else {
			if n != 0 {
				start = append(start, n)
				n = 0
			}
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
	log.Printf("Исходный канал сформирован\n")
	input := init(start)
	pipline := Filtr(NegativeNumber(input, done, &buf))
	for v := range pipline {
		log.Printf("Отфильтровонное число:%d\n", v)
	}
	log.Printf("Конец работы программы")
}
