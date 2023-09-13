package domain

type GPIOClient interface {
	StartAlice() error
	StopAlice() error
	IsPwOn() (bool, error)
}
