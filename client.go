package go_mongo

import (
	"cc/go-mongo/user"
	"cc/go-mongo/userinfo"
)

type Client struct {
	config
	User *UserClient
	UserInfo *UserInfoClient
}

func (c *Client) init(){
	c.User = NewUserClient(c.config)
	c.UserInfo = NewUserInfoClient(c.config)
}
func NewClient(opts ...Option) *Client {
	cfg := config{}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}
type UserClient struct {
	config
	dbName string
}
func (c *UserClient)SetDBName(dbName string)*UserClient{
	c.dbName=dbName
	return c
}
func NewUserClient(c config) *UserClient {
	return &UserClient{ config: c }
}
func(c *UserClient) Query() *UserQuery {
	return &UserQuery{ 
		config: c.config,
		Predicates: []user.UserPredicate{},
		dbName: c.dbName,
	}
}
type UserInfoClient struct {
	config
	dbName string
}
func (c *UserInfoClient)SetDBName(dbName string)*UserInfoClient{
	c.dbName=dbName
	return c
}
func NewUserInfoClient(c config) *UserInfoClient {
	return &UserInfoClient{ config: c }
}
func(c *UserInfoClient) Query() *UserInfoQuery {
	return &UserInfoQuery{ 
		config: c.config,
		Predicates: []userinfo.UserInfoPredicate{},
		dbName: c.dbName,
	}
}
