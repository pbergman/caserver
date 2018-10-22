package config

import (
	"crypto/sha256"
	"crypto/x509/pkix"
	"errors"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type AppConfig struct {
	Path        string `default:"/var/lib/caserver"`
	Address     string `default:":8080"`
	Key         [32]byte
	CaNotAfter  [3]int `default:"10"`
	PemNotAfter [3]int `default:"10"`
}

type Config struct {
	AppConfig `ini:"app"`
	CaSubject *pkix.Name `ini:"ca"`
}

func (c *Config) parseIntArray(value string, dst *[3]int) {
	parts := strings.Split(value, ",")
	switch len(parts) {
	case 1:
		y, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		*dst = [3]int{y, 0, 0}
	case 2:
		y, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		m, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		*dst = [3]int{y, m, 0}
	case 3:
		y, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		m, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		d, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		*dst = [3]int{y, m, d}
	}
}

func (c *Config) Read(source string) error {

	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, source)

	if err != nil {
		return err
	}

	if section, err := cfg.GetSection("app"); err == nil {
		if err := c.readAppSection(section); err != nil {
			return err
		}
	}

	if section, err := cfg.GetSection("ca"); err == nil {
		if c.CaSubject == nil {
			c.CaSubject = new(pkix.Name)
		}
		if err := c.readCaSection(section, c.CaSubject); err != nil {
			return err
		}
	} else {
		return errors.New("missing required `ca` section in config")
	}

	return nil
}

func (c *Config) readCaSection(conf *ini.Section, ca *pkix.Name) error {
	if conf.HasKey("country") {
		ca.Country = []string{conf.Key("country").String()}
	}
	if conf.HasKey("organization") {
		ca.Organization = []string{conf.Key("organization").String()}
	}
	if conf.HasKey("organizational_unit") {
		ca.OrganizationalUnit = []string{conf.Key("organizational_unit").String()}
	}
	if conf.HasKey("locality") {
		ca.Locality = []string{conf.Key("locality").String()}
	}
	if conf.HasKey("province") {
		ca.Province = []string{conf.Key("province").String()}
	}
	if conf.HasKey("street_address") {
		ca.StreetAddress = []string{conf.Key("street_address").String()}
	}
	if conf.HasKey("postal_code") {
		ca.PostalCode = []string{conf.Key("postal_code").String()}
	}
	if conf.HasKey("serial_number") {
		ca.SerialNumber = conf.Key("serial_number").String()
	}

	if conf.HasKey("common_name") {

		ca.CommonName = conf.Key("common_name").String()
	}

	if ca.CommonName == "" {
		return errors.New("missing required `ca.commen_name` field in config section")
	} else {
		return nil
	}
}

func (c *Config) readAppSection(conf *ini.Section) error {
	if conf.HasKey("path") {
		c.Path = conf.Key("path").String()
	}
	if conf.HasKey("address") {
		c.Address = conf.Key("address").String()
	}
	if conf.HasKey("key") {
		c.Key = sha256.Sum256([]byte(conf.Key("key").String()))
	}
	if conf.HasKey("ca_not_after") {
		c.parseIntArray(conf.Key("ca_not_after").String(), &c.CaNotAfter)
	}
	if conf.HasKey("pem_not_after") {
		c.parseIntArray(conf.Key("pem_not_after").String(), &c.PemNotAfter)
	}
	return nil
}
