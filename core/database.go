package core

// Get retrives a string key from the database.
func (c *Client) Get(bucket, key string) string {
	val := c.GetBytes(bucket, key)
	return BytesToString(val)
}

// GetBytes retrives a []byte key from the database.
func (c *Client) GetBytes(bucket, key string) []byte {
	tx := c.DB.NewTransaction(false)
	defer tx.Discard()

	item, err := tx.Get(StringToBytes(bucket + key))
	if err != nil {
		panic(err)
	}

	val, err := item.Value()
	if err != nil {
		panic(err)
	}

	return val
}

// Set sets a string key-value pair in the database.
func (c *Client) Set(bucket, key, value string) {
	c.SetBytes(bucket, key, StringToBytes(value))
}

// SetBytes sets a []byte key in the database.
func (c *Client) SetBytes(bucket, key string, value []byte) {
	tx := c.DB.NewTransaction(true)
	defer tx.Discard()

	err := tx.Set(StringToBytes(bucket+key), value)
	if err != nil {
		panic(err)
	}

	err = tx.Commit(nil)
	if err != nil {
		panic(err)
	}
}
