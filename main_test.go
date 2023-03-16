package main

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestRedisObjectDB(t *testing.T) {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Create RedisObjectDB instance
	db := &RedisObjectDB{client: rdb}

	// Create and store objects
	person1 := &Person{Name: "Alice", Age: 25}
	if err := db.Store(context.Background(), person1); err != nil {
		t.Errorf("failed to store person: %v", err)
	}

	person2 := &Person{Name: "Bob", Age: 30}
	if err := db.Store(context.Background(), person2); err != nil {
		t.Errorf("failed to store person: %v", err)
	}

	animal1 := &Animal{Name: "Rover", Type: "Golden Retriever"}
	if err := db.Store(context.Background(), animal1); err != nil {
		t.Errorf("failed to store animal: %v", err)
	}

	// Retrieve objects by ID
	id := person1.GetID()
	obj, err := db.GetObjectByID(context.Background(), id)
	if err != nil {
		t.Errorf("failed to get person by ID: %v", err)
	}
	if obj.GetName() != "Alice" {
		t.Errorf("unexpected person name: %s", obj.GetName())
	}

	id = animal1.GetID()
	obj, err = db.GetObjectByID(context.Background(), id)
	if err != nil {
		t.Errorf("failed to get animal by ID: %v", err)
	}
	if obj.GetName() != "Rover" {
		t.Errorf("unexpected animal name: %s", obj.GetName())
	}

	// Retrieve objects by name
	name := person2.GetName()
	obj, err = db.GetObjectByName(context.Background(), name)
	if err != nil {
		t.Errorf("failed to get person by name: %v", err)
	}
	if obj.GetID() == "" {
		t.Errorf("person ID is empty")
	}

	name = "Fluffy"
	_, err = db.GetObjectByName(context.Background(), name)
	if err == nil {
		t.Errorf("expected error for missing object")
	}

	// List objects
	objects, err := db.ListObjects(context.Background(), "person")
	if err != nil {
		t.Errorf("failed to list persons: %v", err)
	}
	if len(objects) != 2 {
		t.Errorf("unexpected number of persons: %d", len(objects))
	}

	objects, err = db.ListObjects(context.Background(), "animal")
	if err != nil {
		t.Errorf("failed to list animals: %v", err)
	}
	if len(objects) != 1 {
		t.Errorf("unexpected number of animals: %d", len(objects))
	}

	// Delete object
	id = person2.GetID()
	if err := db.DeleteObject(context.Background(), id); err != nil {
		t.Errorf("failed to delete person: %v", err)
	}
	obj, err = db.GetObjectByID(context.Background(), id)
	if err == nil {
		t.Errorf("expected error for missing object")
	}

	db.client.FlushAllAsync(context.Background())
}
