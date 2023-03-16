package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create a new Redis client.
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create a new RedisObjectDB instance.
	db := &RedisObjectDB{
		client: client,
	}

	// Create a new context for the Redis operations.
	ctx := context.Background()

	// Create some test objects.
	person1 := &Person{Name: "Alice", Age: 30, ID: "1"}
	person2 := &Person{Name: "Bob", Age: 25, ID: "2"}
	animal1 := &Animal{Name: "Fluffy", Type: "cat", ID: "3"}
	animal2 := &Animal{Name: "Rover", Type: "dog", ID: "4"}

	// Store the objects in the database.
	if err := db.Store(ctx, person1); err != nil {
		log.Fatalf("failed to store person1: %v", err)
	}
	if err := db.Store(ctx, person2); err != nil {
		log.Fatalf("failed to store person2: %v", err)
	}
	if err := db.Store(ctx, animal1); err != nil {
		log.Fatalf("failed to store animal1: %v", err)
	}
	if err := db.Store(ctx, animal2); err != nil {
		log.Fatalf("failed to store animal2: %v", err)
	}

	// Get an object by ID.
	obj, err := db.GetObjectByID(ctx, person1.GetID())
	if err != nil {
		log.Fatalf("failed to get object by ID:`` %s, %v", person1.GetID(), err)
	}
	person := obj.(*Person)
	fmt.Printf("got person by ID: %q: \n%+v\n", person1.GetID(), person)

	// Get an object by name.
	obj, err = db.GetObjectByName(ctx, animal1.GetName())
	if err != nil {
		log.Fatalf("failed to get object by name: %v", err)
	}
	animal := obj.(*Animal)
	fmt.Printf("got animal by name: %q:  \n%+v\n", animal1.GetName(), animal)

	// List objects of a certain kind.
	objects, err := db.ListObjects(ctx, "person")
	if err != nil {
		log.Fatalf("failed to list objects: %v", err)
	}
	fmt.Println("list of person type of objects")
	for _, obj := range objects {
		fmt.Printf("%+v\n", obj.(*Person))
	}

	// Delete an object.
	if err := db.DeleteObject(ctx, animal2.GetID()); err != nil {
		log.Fatalf("failed to delete object: %v", err)
	}
}
