package cloudconfig

import "github.com/pkg/errors"

type ClientOption func(c *CloudConfigClient) error

func WithBranch(branch string) ClientOption {
	return func(c *CloudConfigClient) error {
		if branch == "" {
			return errors.New("branch must not be empty")
		}

		c.branch = branch
		return nil
	}
}

func WithFormat(format Format) ClientOption {
	return func(c *CloudConfigClient) error {
		if !format.Valid() {
			return errors.Errorf("[%s] is not an acceptable format", format)
		}

		c.format = format
		return nil
	}
}

func WithBasicAuth(username, password string) ClientOption {
	return func(c *CloudConfigClient) error {
		if username == "" {
			return errors.New("username must not be empty")
		}

		c.basicAuth = &basicAuthInfo{
			username: username,
			password: password,
		}
		return nil
	}
}
