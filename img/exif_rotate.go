package img

import (
	"os/exec"

	"github.com/pkg/errors"
)

func ExifRotate(path string) error {
	bin, err := exec.LookPath("exiftran")
	if err != nil {
		return errors.Wrap(err, "no exiftran for auto rotate")
	}

	cmdLine := []string{"-a", "-i", path}
	cmd := exec.Command(bin, cmdLine...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to auto rotate image: %s", output)
	}
	return nil
}
