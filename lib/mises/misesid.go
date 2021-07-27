package mises

type User struct {
	ID string
}

type Client interface {
	Auth(misesid, code string) error
}

type ClientImpl struct {
}

// TODO mises auth
func (c *ClientImpl) Auth(misesid, code string) error {
	return nil
}

func New() Client {
	return &ClientImpl{}
}
