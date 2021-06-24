package script

// Test runs goconvey
func Test() []string {
	commands := make([]string, 3)
	commands = append(commands, "go get github.com/smartystreets/goconvey")
	commands = append(commands, "goconvey")
	return commands
}
