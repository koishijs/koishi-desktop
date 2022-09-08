//nolint:wrapcheck
package instance

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
)

func Instances(i *do.Injector) ([]string, error) {
	var err error

	config := do.MustInvoke[*koiconfig.Config](i)

	entries, err := os.ReadDir(config.Computed.DirInstance)
	if err != nil {
		return nil, err
	}

	instances := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		instance := entry.Name()
		exists, err := IsInstanceExists(i, instance)
		if err != nil {
			return nil, err
		}
		if exists {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func IsInstanceExists(i *do.Injector, name string) (bool, error) {
	var err error

	config := do.MustInvoke[*koiconfig.Config](i)
	instanceDir := filepath.Join(config.Computed.DirInstance, name)

	_, err = os.Stat(instanceDir)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check instance %s: %w", name, err)
	}

	for _, f := range []string{
		// Deprecated koishi.config.yml are not supported
		filepath.Join(instanceDir, "koishi.yml"),
		filepath.Join(instanceDir, "package.json"),
	} {
		_, err = os.Stat(f)
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}

	return true, nil
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
