package account

type NotificationMessage int

const (
	UserLoggedIn NotificationMessage = iota + 1
	UserProfileUpdated
)
