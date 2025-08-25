package main_test

import (
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddListCmd(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	binName := "./main"
	run(t, "go", "build", "-o", binName)

	run(t, binName, "add", "--amount", "250", "--name", "Groceries - supermarket", "--nws", "needs")

	output := runAndGetOutput(t, binName, "list")

	assert.Equal(t, `name,amount,nws
"Groceries - supermarket",250,needs`, output)

	run(t, binName, "add", "--amount", "700", "--name", "new iPhone", "--nws", "wants")

	output = runAndGetOutput(t, binName, "list")

	assert.Equal(t, `name,amount,nws
"Groceries - supermarket",250,needs
"new iPhone",700,wants`, output)
}

func run(t *testing.T, binName string, args ...string) {
	cmd := exec.Command(binName, args...)
	err := cmd.Run()
	assert.NoError(t, err)
}

func runAndGetOutput(t *testing.T, binName string, args ...string) string {
	cmd := exec.Command(binName, args...)

	stdout, err := cmd.StdoutPipe()
	assert.NoError(t, err)

	assert.NoError(t, cmd.Start())

	output, err := io.ReadAll(stdout)
	assert.NoError(t, err)

	assert.NoError(t, cmd.Wait())

	return string(output)
}
