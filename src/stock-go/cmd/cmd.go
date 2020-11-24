package cmd

type CommandInterface interface {
	Run(param map[string]string) error
}

var (
	CommandMap map[string]CommandInterface
)

func setCommand(cmd string, data CommandInterface) {
	if CommandMap == nil {
		CommandMap = map[string]CommandInterface{}
	}
	CommandMap[cmd] = data
}
