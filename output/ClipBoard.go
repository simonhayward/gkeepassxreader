package output

import (
	"io"
	"os/exec"
	"runtime"
)

type execCommand struct{}

func (ec *execCommand) Process(cmds []string, text string) error {
	var args []string
	c := cmds[0]

	if len(cmds) > 1 {
		args = cmds[1:]
	}

	cmd := exec.Command(c, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, text)
	}()

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// ClipBoard interface
type ClipBoard interface {
	CopyProcess(string) error
}

type oSXClipBoard struct {
	execCommand
}
type linuxClipBoard struct {
	execCommand
}
type windowsClipBoard struct {
	execCommand
}

func (osxC *oSXClipBoard) CopyProcess(text string) error {
	return osxC.Process([]string{"pbcopy"}, text)
}

func (osxC *linuxClipBoard) CopyProcess(text string) error {
	return osxC.Process([]string{"xclip", "-selection", "clipboard"}, text)

}

func (osxC *windowsClipBoard) CopyProcess(text string) error {
	return osxC.Process([]string{"clip"}, text)
}

//GetClipboard for OS
func GetClipboard() ClipBoard {
	goos := runtime.GOOS
	switch {
	case goos == "windows":
		return &windowsClipBoard{}
	case goos == "linux":
		return &linuxClipBoard{}
	case goos == "darwin":
		return &oSXClipBoard{}
	}
	return nil
}
