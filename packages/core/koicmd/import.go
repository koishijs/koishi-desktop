package koicmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/koierr"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/compress"
)

func koiCmdImport(i *do.Injector) *proto.Response {
	var err error

	l := do.MustInvoke[*logger.Logger](i)
	command := do.MustInvoke[*proto.CommandRequest](i)
	config := do.MustInvoke[*koiconfig.Config](i)

	l.Debug("Trigger KoiCmd import")

	// Parse command
	path, ok := command.Flags["path"].(string)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}
	name, ok := command.Flags["name"].(string)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}
	force, ok := command.Flags["force"].(bool)
	if !ok {
		return proto.NewErrorResult(koierr.ErrBadRequest)
	}
	if path == "" {
		return proto.NewErrorResult(koierr.NewErrBadRequest(errors.New("must provide path")))
	}

	// Auto generate name
	if name == "" {
		name, err = generateInstanceName(i)
		if err != nil {
			return proto.NewErrorResult(koierr.NewErrInternalError(err))
		}
	}

	targetPath := filepath.Join(config.Computed.DirInstance, name)
	exists, err := isInstanceExists(i, name)
	if err != nil {
		return proto.NewErrorResult(koierr.NewErrInternalError(err))
	}
	if exists {
		if force {
			err = os.RemoveAll(targetPath)
			if err != nil {
				return proto.NewErrorResult(koierr.NewErrInternalError(err))
			}
		} else {
			return proto.NewErrorResult(koierr.NewErrInstanceExists(name))
		}
	}

	l.Infof("Importing instance %s\nUsing Koishi bundle: %s", name, path)
	err = os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return proto.NewErrorResult(koierr.NewErrInternalError(fmt.Errorf("cannot create target path %s: %w", targetPath, err)))
	}

	err = compress.ExtractZipFile(path, targetPath)
	if err != nil {
		return proto.NewErrorResult(koierr.NewErrInternalError(fmt.Errorf("failed to unzip %s: %w", path, err)))
	}

	return proto.NewErrorResult(koierr.ErrSuccess)
}
