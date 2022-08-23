//nolint:wrapcheck
package instance

import (
	"errors"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
)

func IsInstanceExists(i *do.Injector, name string) (bool, error) {
	config := do.MustInvoke[*koiconfig.Config](i)
	_, err := os.Stat(filepath.Join(config.Computed.DirInstance, name))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func GenerateInstanceName(i *do.Injector) (string, error) {
	var err error

	var exists bool
	prefix := "default"
	name := prefix
	exists, err = IsInstanceExists(i, name)
	if err != nil {
		return "", err
	}
	if !exists {
		return name, nil
	}
	for index := 1; index < math.MaxUint16; index++ {
		name = prefix + strconv.Itoa(index)
		exists, err = IsInstanceExists(i, name)
		if err != nil {
			return "", err
		}
		if !exists {
			return name, nil
		}
	}

	return "", errors.New("max instance count exceed")
}
