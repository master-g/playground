package landmine

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Entry function for package
func Entry() {
	err := mine3(true)
	if err != nil {
		fmt.Println(err)
	}
}

func print(pi *int) {
	fmt.Printf("defer print(&i) = %v\n", *pi)
}

// lopp variable is scoped outside loop
func mine1() {
	for i := 0; i < 3; i++ {
		defer fmt.Printf("defer fmt.Println(i) = %v\n", i)
		defer func() { fmt.Printf("defer func() { fmt.Println(i) }() = %v\n", i) }()
		defer func(i int) { fmt.Printf("defer func(i int) { fmt.Println(i) }(i) = %v\n", i) }(i)
		defer print(&i)
		go fmt.Printf("go fmt.Println(i) = %v\n", i)
		go func() { fmt.Printf("go func(){ fmt.Println(i) }() = %v\n", i) }()
	}
	time.Sleep(time.Second * 1)
	fmt.Println("--- sleep over ---")
}

type Cat interface {
	Meow()
}

type Tabby struct{}

func (*Tabby) Meow() { fmt.Println("meow") }

func GetACat() Cat {
	var myTabby *Tabby
	// Oops, something went wrong, no real value here
	return myTabby
}

var ErrRua = errors.New("rua")

func mine3(forReal bool) (err error) {
	if forReal {
		result, err := doit()
		if err != nil || result != "ok" {
			err = ErrRua
		}
	}
	return err
}

func doit() (string, error) {
	if rand.Intn(100) > 50 {
		return "ok", nil
	}
	return "rua", errors.New("shit happens")
}
