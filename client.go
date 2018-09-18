package waffleio

type Client struct {
	token string
}

func New(token string) (*Client, error) {
	return &Client{
		token: token,
	}, nil
}
