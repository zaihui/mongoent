package go_mongo

type Client struct {
	config
	User *UserClient
	UserInfo *UserInfoClient
}

type UserClient struct {
	config
}
func NewUser(c config) *UserClient {
	return &UserClient{ config: c }
}
func(c *UserClient) Query() *UserQuery {
	return &UserQuery{ config: c.config }
}
type UserInfoClient struct {
	config
}
func NewUserInfo(c config) *UserInfoClient {
	return &UserInfoClient{ config: c }
}
func(c *UserInfoClient) Query() *UserInfoQuery {
	return &UserInfoQuery{ config: c.config }
}
