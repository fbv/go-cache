package main

import (
	"log"
	"time"

	"github.com/fbv/go-cache"
)

type Person struct {
	Name string
}

func newPerson(name string) (*Person, error) {
	log.Println("cache miss:", name)
	return &Person{
		Name: name,
	}, nil
}

func tryStrategy(s cache.ExpirationStrategy[*Person], names []string) {
	people := cache.New[string, *Person](s, newPerson)
	for _, name := range names {
		p, err := people.Get(name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(p)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	names := []string{"James", "John", "Robert", "John", "Michael", "Michael", "Michael", "James"}

	log.Println("No expiration")
	tryStrategy(cache.NoExpiration[*Person](), names)

	log.Println("----------")
	log.Println("Expire after 1000ms")
	tryStrategy(cache.ExpireAfter[*Person](1*time.Second), names)

	log.Println("----------")
	log.Println("Expire after 1000ms not used")
	tryStrategy(cache.LastAccess[*Person](1*time.Second), names)
}
