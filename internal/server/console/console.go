package console

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// Console represents the server console and can be used to interact with the minecraft server.
type Console struct {
	stdin  *bufio.Writer
	stderr *bufio.Reader
	stdout *bufio.Reader

	subscriber map[string]func(line string)
}

// NewConsole creates a new console
func NewConsole(stdin io.Writer, stderr io.Reader, stdout io.Reader) *Console {
	return &Console{
		stdin:  bufio.NewWriter(stdin),
		stderr: bufio.NewReader(stderr),
		stdout: bufio.NewReader(stdout),

		subscriber: make(map[string]func(line string)),
	}
}

// Start will start the console, it will follow the logs
func (c *Console) Start() {
	for {
		if l, err := c.readLine(); err == nil {
			go func() {
				for _, s := range c.subscriber {
					s(l)
				}
			}()
			fmt.Print(l)
		}
	}
}

// Subscribe subscribes to the log stream
func (c *Console) Subscribe(name string, handler func(line string)) {
	c.subscriber[name] = handler
}

// Unsubscribe unsubscribes from the log stream
func (c *Console) Unsubscribe(name string) {
	delete(c.subscriber, name)
}

// SendCommand sends a command to the minecraft server
func (c *Console) SendCommand(cmd string) error {
	logrus.WithField("cmd", cmd).Debug("sending command")
	if _, err := c.stdin.WriteString(fmt.Sprintf("%s\r\n", cmd)); err != nil {
		logrus.WithError(err).Error("sending command failed")
		return err
	}
	logrus.Trace("send command")
	return c.stdin.Flush()
}

func (c *Console) readLine() (string, error) {
	if c.stdout == nil {
		return "", io.EOF
	}
	l, err := c.stdout.ReadString('\n')
	if err == io.EOF {
		return "", io.EOF
	}
	return l, nil
}
