package account

type AccountCreatedEventType struct{}

func (AccountCreatedEventType) Name() string {
	return "AccountCreated"
}

type AccountDeactivatedEventType struct{}

func (AccountDeactivatedEventType) Name() string {
	return "AccountDeactivated"
}

var (
	AccountCreated     = AccountCreatedEventType{}
	AccountDeactivated = AccountDeactivatedEventType{}
)
