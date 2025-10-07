package store


type clickHouse struct{}

func NewClickHouse(addr, db, user, pass string, codec map[string]string) (Store, error) {
	return &clickHouse{}, nil
}
func (c *clickHouse) Close() error { return nil }
