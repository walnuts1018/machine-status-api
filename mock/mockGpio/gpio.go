package mockGpio

import "github.com/walnuts1018/machine-status-api/domain"

type mockClient struct {
	isAlicePwOn bool
}

func NewClient() domain.GPIOClient {
	return &mockClient{}
}

func (c *mockClient) StartAlice() error {
	c.isAlicePwOn = true
	return nil
}

func (c *mockClient) StopAlice() error {
	c.isAlicePwOn = false
	return nil
}

func (c *mockClient) IsPwOn() (bool, error) {
	return c.isAlicePwOn, nil
}
