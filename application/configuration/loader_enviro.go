// Sshwifty - A Web SSH client
//
// Copyright (C) 2019-2020 Rui NI <nirui@gmx.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/niruix/sshwifty/application/log"
)

const (
	enviroTypeName = "Environment Variable"
)

func parseEviro(name string) string {
	v := os.Getenv(name)

	if !strings.HasPrefix(v, "SSHWIFTY_ENV_RENAMED:") {
		return v
	}

	return os.Getenv(v[21:])
}

// Enviro creates an environment variable based configuration loader
func Enviro() Loader {
	return func(log log.Logger) (string, Configuration, error) {
		log.Info("Loading configuration from environment variables ...")

		dialTimeout, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_DIALTIMEOUT"), 10, 32)

		cfg, dialer, cfgErr := fileCfgCommon{
			HostName:       parseEviro("SSHWIFTY_HOSTNAME"),
			SharedKey:      parseEviro("SSHWIFTY_SHAREDKEY"),
			DialTimeout:    int(dialTimeout),
			Socks5:         parseEviro("SSHWIFTY_SOCKS5"),
			Socks5User:     parseEviro("SSHWIFTY_SOCKS5_USER"),
			Socks5Password: parseEviro("SSHWIFTY_SOCKS5_PASSWORD"),
		}.build()

		if cfgErr != nil {
			return enviroTypeName, Configuration{}, fmt.Errorf(
				"Failed to build the configuration: %s", cfgErr)
		}

		listenPort, listenPortErr := strconv.ParseUint(
			parseEviro("SSHWIFTY_LISTENPORT"), 10, 16)

		if listenPortErr != nil {
			return enviroTypeName, Configuration{}, fmt.Errorf(
				"Invalid \"SSHWIFTY_LISTENPORT\": %s", listenPortErr)
		}

		initialTimeout, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_INITIALTIMEOUT"), 10, 32)

		readTimeout, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_READTIMEOUT"), 10, 32)

		writeTimeout, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_WRITETIMEOUT"), 10, 32)

		heartbeatTimeout, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_HEARTBEATTIMEOUT"), 10, 32)

		readDelay, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_READDELAY"), 10, 32)

		writeDelay, _ := strconv.ParseUint(
			parseEviro("SSHWIFTY_WRITEELAY"), 10, 32)

		cfgSer := fileCfgServer{
			ListenInterface:       parseEviro("SSHWIFTY_LISTENINTERFACE"),
			ListenPort:            uint16(listenPort),
			InitialTimeout:        int(initialTimeout),
			ReadTimeout:           int(readTimeout),
			WriteTimeout:          int(writeTimeout),
			HeartbeatTimeout:      int(heartbeatTimeout),
			ReadDelay:             int(readDelay),
			WriteDelay:            int(writeDelay),
			TLSCertificateFile:    parseEviro("SSHWIFTY_TLSCERTIFICATEFILE"),
			TLSCertificateKeyFile: parseEviro("SSHWIFTY_TLSCERTIFICATEKEYFILE"),
		}

		return enviroTypeName, Configuration{
			HostName:    cfg.HostName,
			SharedKey:   cfg.SharedKey,
			Dialer:      dialer,
			DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
			Servers:     []Server{cfgSer.build()},
		}, nil
	}
}
