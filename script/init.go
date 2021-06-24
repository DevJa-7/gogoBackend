package script

// InitApp runs initialize commands by scrips
func InitApp() []string {
	commands := make([]string, 8)

	commands = append(commands, "go get")
	commands = append(commands, "brew install glide")
	commands = append(commands, "glide install")
	commands = append(commands, "glide up")
	return commands
}
