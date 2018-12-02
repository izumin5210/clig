package foobar

import (
	"github.com/izumin5210/clig/pkg/cli"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

type Ctx struct {
	WorkingDir cli.Path
	IO         cli.IO
	FS         afero.Fs
	Viper      *afero.Viper
	Exec       exec.Interface

	Build  cli.Build
	Config *Config
}

func (c *Ctx) Init() error {
	c.Viper.SetFs(c.FS)

	var err error

	err = c.loadConfig()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Ctx) loadConfig() error {
	c.Viper.SetConfigName(c.Build.AppName)

	err := c.Viper.ReadInConfig()
	if err != nil {
		zap.L().Info("failed to find a config file", zap.Error(err))
		return nil
	}

	err = c.Viper.Unmarshal(c.Config)
	if err != nil {
		zap.L().Warn("failed to parse the config file", zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}

