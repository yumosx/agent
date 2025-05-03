package agent

import "github.com/yumosx/agent/internal/domain"

type Bash struct {
}

const bashDescription = `Execute a bash command in the terminal.
* Long running commands: For commands that may run indefinitely, it should be run in the background and the output should be redirected to a file, e.g. command = "python3 app.py > server.log 2>&1 &".
* Interactive: If a bash command returns exit code -1, this means the process is not yet finished. The assistant must then send a second call to terminal with an empty "command" (which will retrieve any additional logs), or it can send additional text (set "command" to the text) to STDIN of the running process, or it can send command="ctrl+c" to interrupt the process.
* Timeout: If a command execution result says "Command timed out. Sending SIGINT to the process", the assistant should retry running the command in the background.`

func (b *Bash) NewTool() domain.Tool {
	var t domain.Tool
	t.Type = "function"
	t.Function = domain.Function{
		Name:        "bash",
		Description: bashDescription,
		Parameters: &domain.FunctionParameters{
			Required: []string{"command"},
		},
	}

	return t
}
