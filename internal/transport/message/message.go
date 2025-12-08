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

	WrongId          = "Provided wrong id, entry not found"
	WrongActionId    = "No action corresponds to given input, try again!"
	PasswordMismatch = "Entered password are not matching"
	UnhandledError   = "Encountered error: "
	NoEntriesFound   = "No entries found"
	WrongPassword    = "Wrong password"

	NewLine       = "\n"
	LongSeparator = "--------------------------------------------\n"
)
