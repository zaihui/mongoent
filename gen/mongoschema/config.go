package mongoschema

import "go.mongodb.org/mongo-driver/mongo"

type config struct {
	mongo.Client
}
type Option func(*config)
func (c *config) options(opts ...Option) {	for _, opt := range opts {
		opt(c)
	}
}
func Driver(driver mongo.Client) Option {
	return func(c *config) {
		c.Client = driver
	}
}
