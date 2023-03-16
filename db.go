package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Object interface {
	// GetKind returns the type of the object.
	GetKind() string
	// GetID returns a unique UUID for the object.
	GetID() string
	// GetName returns the name of the object. Names are not unique.
	GetName() string
	// SetID sets the ID of the object.
	SetID(string)
	// SetName sets the name of the object.
	SetName(string)
}

// Using the data store of your choice you need to provide an implementation of the following interface:
type ObjectDB interface {
	// Store will store the object in the data store. The object will have a
	// name and kind, and the Store method should create a unique ID.
	Store(ctx context.Context, object Object) error
	// GetObjectByID will retrieve the object with the provided ID.
	GetObjectByID(ctx context.Context, id string) (Object, error)
	// GetObjectByName will retrieve the object with the given name.
	GetObjectByName(ctx context.Context, name string) (Object, error)
	// ListObjects will return a list of all objects of the given kind.
	ListObjects(ctx context.Context, kind string) ([]Object, error)
	// DeleteObject will delete the object.
	DeleteObject(ctx context.Context, id string) error
}

// The data store should be able to store multiple objects of different kinds.
// For example, it should be able to store both Person and Animal objects.
func StoreObject(ctx context.Context, db ObjectDB, object Object) error {
	return db.Store(ctx, object)
}

// The data store should be able to retrieve objects by ID.
func GetObjectByID(ctx context.Context, db ObjectDB, id string) (Object, error) {
	return db.GetObjectByID(ctx, id)
}

func GetObjectByName(ctx context.Context, db ObjectDB, name string) (Object, error) {
	return db.GetObjectByName(ctx, name)
}

func ListObjects(ctx context.Context, db ObjectDB, kind string) ([]Object, error) {
	return db.ListObjects(ctx, kind)
}

func DeleteObject(ctx context.Context, db ObjectDB, id string) error {
	return db.DeleteObject(ctx, id)
}

type RedisObjectDB struct {
	client *redis.Client
	ObjectDB
}

func NewRedisObjectDB(client *redis.Client) *RedisObjectDB {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisObjectDB{client: client}
}

func (db *RedisObjectDB) Store(ctx context.Context, object Object) error {
	id, err := db.client.Incr(ctx, "object_counter").Result()
	if err != nil {
		return err
	}
	object.SetID(fmt.Sprintf("%d", id))

	key := fmt.Sprintf("%s:%s", object.GetKind(), object.GetID())
	value, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return db.client.Set(ctx, key, value, 0).Err()
}

func (db *RedisObjectDB) GetObjectByID(ctx context.Context, id string) (Object, error) {
	var obj Object

	keys, err := db.client.Keys(ctx, "*:"+id).Result()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("object not found")
	}

	data, err := db.client.Get(ctx, keys[0]).Bytes()
	if err != nil {
		return nil, err
	}

	switch keys[0][:len(keys[0])-len(id)-1] {
	case "person":
		obj = &Person{}
	case "animal":
		obj = &Animal{}
	default:
		return nil, fmt.Errorf("unknown object kind")
	}

	if err := json.Unmarshal(data, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (db *RedisObjectDB) GetObjectByName(ctx context.Context, name string) (Object, error) {
	var obj Object

	keys, err := db.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		data, err := db.client.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}

		switch strings.Split(key, ":")[0] {
		case "person":
			obj = &Person{}
		case "animal":
			obj = &Animal{}
		default:
			continue
		}

		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}

		if obj.GetName() == name {
			return obj, nil
		}
	}

	return nil, fmt.Errorf("object not found")
}

func (db *RedisObjectDB) ListObjects(ctx context.Context, kind string) ([]Object, error) {
	var objects []Object
	keys := db.client.Keys(ctx, kind+":*")
	if keys.Err() != nil {
		return nil, keys.Err()
	}
	for _, key := range keys.Val() {
		value, err := db.client.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}
		var obj Object
		switch kind {
		case "person":
			obj = &Person{}
		case "animal":
			obj = &Animal{}
		default:
			continue
		}

		err = json.Unmarshal(value, &obj)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (db *RedisObjectDB) DeleteObject(ctx context.Context, id string) error {
	keys, err := db.client.Keys(ctx, "*:"+id).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("object not found, id: %s", id)
	}

	if err := db.client.Del(ctx, keys[0]).Err(); err != nil {
		return err
	}
	return nil
}
