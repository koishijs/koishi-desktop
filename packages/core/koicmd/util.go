package koicmd

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

func isInstanceExists(i *do.Injector, name string) (bool, error) {
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

func generateInstanceName(i *do.Injector) (string, error) {
	var err error
	exists := true
	prefix := "default"
	name := prefix
	exists, err = isInstanceExists(i, name)
	if err != nil {
		return "", err
	}
	if !exists {
		return name, nil
	}
	for index := 1; index < math.MaxUint16; index++ {
		name = prefix + strconv.Itoa(index)
		exists, err = isInstanceExists(i, name)
		if err != nil {
			return "", err
		}
		if !exists {
			return name, nil
		}
	}

	return "", errors.New("max instance count exceed")
}
