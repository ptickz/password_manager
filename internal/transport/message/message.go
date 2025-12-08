package message

const (
	FirstTimeEnter         = "Storage not found\nSeems its your first time using this app\nEnter your master password: "
	RepeatMasterPassword   = "Repeat master password: "
	AuthWithPassword       = "Enter master password: "
	AwaitInput             = "Await input: "
	TimeoutMessage         = "Timeout reached, quitting"
	RequestInputId         = "Enter entry id or name: "
	RequestServiceName     = "Enter service name: "
	RequestServiceLogin    = "Enter login: "
	RequestServicePassword = "Enter password: "
	CreationSuccess        = "Entry created successfully"
	PasswordMismatch       = "Entered password are not matching"

	WrongId        = "Provided wrong id, entry not found\n"
	WrongActionId  = "No action corresponds to given input, try again!\n"
	UnhandledError = "Encountered error: "

	NewLine       = "\n"
	LongSeparator = "--------------------------------------------\n"
)
