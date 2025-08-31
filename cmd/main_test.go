package main_test

import (
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	binName := "./main"
	assert.NoError(t, run("go", "build", "-o", binName))

	t.Run("Wrong args", func(t *testing.T) {
		assert.Error(t, run(binName))
		assert.Error(t, run(binName, "nonexistentcommand"))
		assert.Error(t, run(binName, "add"))
		assert.Error(t, run(binName, "add", "--nws", "test"))
		assert.Error(t, run(binName, "add", "--amount", "-10"))
		assert.Error(t, run(binName, "list", "args"))
	})

	t.Run("Add + List", func(t *testing.T) {
		assert.NoError(t, run(binName, "add", "--amount", "250", "--name", "Groceries - supermarket", "--nws", "needs"))

		output := runAndGetOutput(t, binName, "list")

		assert.Equal(t, `name,amount,nws
	"Groceries - supermarket",250,needs`, output)

		assert.NoError(t, run(binName, "add", "--amount", "700", "--name", "new iPhone", "--nws", "wants"))

		output = runAndGetOutput(t, binName, "list")

		assert.Equal(t, `name,amount,nws
	"Groceries - supermarket",250,needs
	"new iPhone",700,wants`, output)
	})
}

func run(binName string, args ...string) error {
	cmd := exec.Command(binName, args...)
	return cmd.Run()
}

func runAndGetOutput(t *testing.T, binName string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binName, args...)

	stdout, err := cmd.StdoutPipe()
	assert.NoError(t, err)

	assert.NoError(t, cmd.Start())

	output, err := io.ReadAll(stdout)
	assert.NoError(t, err)

	assert.NoError(t, cmd.Wait())

	return string(output)
}
