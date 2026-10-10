package command

type CommandRegistry interface {
	Resolve(code CmdCode) (Command, bool)
}

type commandRegistry struct {
	commands map[CmdCode]Command
}

func NewCommandRegistry() CommandRegistry {
	registry := commandRegistry{}
	registry.commands = make(map[CmdCode]Command)
	registry.commands[CWD] = cwdCommand{}
	registry.commands[EPRT] = eprtCommand{}
	registry.commands[PASS] = passCommand{}
	registry.commands[PORT] = portCommand{}
	registry.commands[QUIT] = quitCommand{}
	registry.commands[STOR] = storCommand{}
	registry.commands[SYST] = systCommand{}
	registry.commands[TYPE] = typeCommand{}
	registry.commands[USER] = userCommand{}
	return registry
}

func (r commandRegistry) Resolve(code CmdCode) (Command, bool) {
	cmd, ok := r.commands[code]
	return cmd, ok
}
