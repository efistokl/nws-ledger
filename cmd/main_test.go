package main_test

import (
	"io"
	"os"
	"os/exec"
	"testing"

	main "github.com/efistokl/nws-ledger/cmd"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	binName := "./main"
	assert.NoError(t, run("go", "build", "-o", binName))

	assert.NoError(t, os.Rename(binName, os.TempDir()+binName))
	assert.NoError(t, os.Chdir(os.TempDir()))
	defer func() {
		assert.NoError(t, os.Remove(binName))
		assert.NoError(t, os.Remove(main.DefaultStoreFile))
	}()

	t.Run("Wrong args", func(t *testing.T) {
		assert.Error(t, run(binName))
		assert.Error(t, run(binName, "nonexistentcommand"))
		assert.Error(t, run(binName, "add"))
		assert.Error(t, run(binName, "add", "--nws", "test"))
		assert.Error(t, run(binName, "add", "--amount", "-10"))
		assert.Error(t, run(binName, "list", "args"))
	})

	t.Run("Add + List + Summary", func(t *testing.T) {
		assert.NoError(t, run(binName, "add", "--amount", "250", "--name", "Groceries - supermarket", "--nws", "needs", "--domain", "groceries"))

		output := runAndGetOutput(t, binName, "list")

		assert.Equal(t, `name,amount,nws,domain
Groceries - supermarket,250,needs,groceries
`, output)

		assert.NoError(t, run(binName, "add", "--amount", "700", "--name", "new iPhone", "--nws", "wants", "--domain", "shopping"))

		output = runAndGetOutput(t, binName, "list")

		assert.Equal(t, `name,amount,nws,domain
Groceries - supermarket,250,needs,groceries
new iPhone,700,wants,shopping
`, output)

		assert.NoError(t, run(binName, "add", "--amount", "700", "--name", "another iPhone", "--nws", "wants", "--domain", "shopping"))

		output = runAndGetOutput(t, binName, "summary")

		assert.Equal(t, `nws,amount
needs,250
wants,1400
savings,0
total,1650
`, output)
	})
}

func run(binName string, args ...string) error {
	cmd := exec.Command(binName, args...)
	return cmd.Run()
}

func runAndGetOutput(t *testing.T, binName string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binName, args...)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	assert.NoError(t, err)

	assert.NoError(t, cmd.Start())

	output, err := io.ReadAll(stdout)
	assert.NoError(t, err)

	assert.NoError(t, cmd.Wait())

	return string(output)
}
