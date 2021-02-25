package hello

import "fmt"

func Bottle99() {
	for bottle := 99; bottle >= 0; bottle-- {
		switch {
		case bottle >= 2:
			fmt.Printf("%d bottles of beer on the wall, %d bottles of beer.\n", bottle, bottle)
			s := "bottles"
			if bottle == 2 {
				s = "bottle"
			}
			fmt.Printf("Take one down, pass it around, %d %s of beer on the wall.\n", bottle-1, s)
		case bottle == 1:
			fmt.Printf("%d bottle of beer on the wall, %d bottles of beer.\n", bottle, bottle)
			fmt.Printf("Take one down, pass it around, No more bottles of beer on the wall.\n")
		default:
			fmt.Printf("No more bottles of beer on the wall, no more bottles of beer.\n")
			fmt.Printf("Go to the store and buy some more, 99 bottles of beer on the wall.\n")
		}

	}
}

func Fizzbuzz() {
	for i := 1; i <= 100; i++ {
		switch {
		case i%3 == 0 && i%5 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}
