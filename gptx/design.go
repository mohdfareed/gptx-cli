package main

type AppTheme struct {
	// colors
	Reset, Bold, Dim, Black, White          string
	Red, Green, Yellow, Blue, Magenta, Cyan string
	// symbols
	AppIcon, ChatIcon, HistoryIcon, SettingsIcon string
	ToolIcon, WebSearchIcon, ShellIcon           string
	FileIcon, FolderIcon, RemoteIcon             string
	// status
	On, Off, Success, Error       string
	Warning, Info, Unknown, Debug string
}

var Theme AppTheme = PlainTheme
var PlainTheme = AppTheme{}
var RichTheme = AppTheme{
	// MARK: Colors
	Reset:   "\033[0m",
	Bold:    "\033[1m",
	Dim:     "\033[2m",
	Black:   "\033[30m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	White:   "\033[37m",
	// MARK: Symbols
	AppIcon:       " ",
	ChatIcon:      "󰭹 ",
	HistoryIcon:   " ",
	SettingsIcon:  " ",
	ToolIcon:      " ",
	WebSearchIcon: " ",
	ShellIcon:     " ",
	FileIcon:      " ",
	FolderIcon:    " ",
	RemoteIcon:    " ",
}

func init() { // runtime init
	// MARK: Status
	RichTheme.On = RichTheme.Bold + " " + RichTheme.Reset
	RichTheme.Off = RichTheme.Dim + " " + RichTheme.Reset
	RichTheme.Success = RichTheme.Green + " " + RichTheme.Reset
	RichTheme.Error = RichTheme.Red + " " + RichTheme.Reset
	RichTheme.Warning = RichTheme.Yellow + " " + RichTheme.Reset
	RichTheme.Info = RichTheme.Blue + " " + RichTheme.Reset
	RichTheme.Unknown = RichTheme.White + RichTheme.Dim + " " + RichTheme.Reset
	RichTheme.Debug = RichTheme.Magenta + " " + RichTheme.Reset
}

// MARK: Legacy
// ============================================================================

// const (
// 	Reset   = "\033[0m"
// 	Bold    = "\033[1m"
// 	Dim     = "\033[2m"
// 	Black   = "\033[30m"
// 	Red     = "\033[31m"
// 	Green   = "\033[32m"
// 	Yellow  = "\033[33m"
// 	Blue    = "\033[34m"
// 	Magenta = "\033[35m"
// 	Cyan    = "\033[36m"
// 	White   = "\033[37m"
// )

// const (
// 	AppIcon       = " "
// 	ChatIcon      = "󰭹 "
// 	HistoryIcon   = " "
// 	SettingsIcon  = " "
// 	ToolIcon      = " "
// 	WebSearchIcon = " "
// 	ShellIcon     = " "
// 	FileIcon      = " "
// 	FolderIcon    = " "
// 	RemoteIcon    = " "

// 	On      = Bold + " " + Reset
// 	Off     = Dim + " " + Reset
// 	Success = Green + " " + Reset
// 	Error   = Red + " " + Reset
// 	Warning = Yellow + " " + Reset
// 	Info    = Blue + " " + Reset
// 	Unknown = White + Dim + " " + Reset
// 	Debug   = Magenta + " " + Reset
// )
