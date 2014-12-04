// Copyright 2014 Canonical Ltd.

package config_test

import (
	"io/ioutil"
	"path"
	"testing"

	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/CanonicalLtd/blues-identity/config"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type configSuite struct {
	jujutesting.IsolationSuite
}

var _ = gc.Suite(&configSuite{})

const testConfig = `
mongo-addr: localhost:23456
api-addr: 1.2.3.4:5678
foo: 1
bar: false
`

func (s *configSuite) readConfig(c *gc.C, content string) (*config.Config, error) {
	// Write the configuration content to file.
	path := path.Join(c.MkDir(), "config.yaml")
	err := ioutil.WriteFile(path, []byte(content), 0666)
	c.Assert(err, gc.IsNil)

	// Read the configuration.
	return config.Read(path)
}

func (s *configSuite) TestRead(c *gc.C) {
	conf, err := s.readConfig(c, testConfig)
	c.Assert(err, gc.IsNil)
	c.Assert(conf, jc.DeepEquals, &config.Config{
		MongoAddr: "localhost:23456",
		APIAddr:   "1.2.3.4:5678",
	})
}

func (s *configSuite) TestReadConfigError(c *gc.C) {
	cfg, err := config.Read(path.Join(c.MkDir(), "no-such-file.yaml"))
	c.Assert(err, gc.ErrorMatches, ".* no such file or directory")
	c.Assert(cfg, gc.IsNil)
}

func (s *configSuite) TestValidateConfigError(c *gc.C) {
	cfg, err := s.readConfig(c, "")
	c.Assert(err, gc.ErrorMatches, "missing fields mongo-addr, api-addr in config file")
	c.Assert(cfg, gc.IsNil)
}
