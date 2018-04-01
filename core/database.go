package core

import (
	"encoding/binary"
)

// Get retrives a string value from the database with a string key.
func (c *Client) Get(bucket, key string) (string, error) {
	val, err := c.GetBytes(bucket, key)
	if err != nil {
		return "", err
	}

	return BytesToString(val), nil
}

// GetByID retrives a string value from the database with a uint64 key.
func (c *Client) GetByID(bucket string, key uint64) (string, error) {
	byteKey := make([]byte, len(bucket)+LenUint64)
	copy(byteKey, bucket)
	binary.LittleEndian.PutUint64(byteKey[len(bucket):], key)

	val, err := c.rawGet(byteKey)
	if err != nil {
		return "", err
	}

	return BytesToString(val), nil
}

// GetBytes retrieves a []byte value from the database with a string key.
func (c *Client) GetBytes(bucket, key string) ([]byte, error) {
	val, err := c.rawGet(StringToBytes(bucket + key))
	return val, err
}

func (c *Client) rawGet(key []byte) ([]byte, error) {
	tx := c.DB.NewTransaction(false)
	defer tx.Discard()

	item, err := tx.Get(key)
	if err != nil {
		return nil, err
	}

	val, err := item.Value()
	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set sets a string value with a string key in the database.
func (c *Client) Set(bucket, key, value string) error {
	return c.SetBytes(bucket, key, StringToBytes(value))
}

// SetByID sets a string value with a uint64 key in the database.
func (c *Client) SetByID(bucket string, key uint64, value string) error {
	byteKey := make([]byte, len(bucket)+LenUint64)
	copy(byteKey, bucket)
	binary.LittleEndian.PutUint64(byteKey[len(bucket):], key)

	return c.rawSet(byteKey, StringToBytes(value))
}

// SetBytes sets a []byte value with a string key in the database.
func (c *Client) SetBytes(bucket, key string, value []byte) error {
	return c.rawSet(StringToBytes(bucket+key), value)
}

func (c *Client) rawSet(key, value []byte) error {
	tx := c.DB.NewTransaction(true)
	defer tx.Discard()

	err := tx.Set(key, value)
	if err != nil {
		return err
	}

	err = tx.Commit(nil)
	if err != nil {
		return err
	}

	return nil
}
