/* -*- indent-tabs-mode:nil; coding: utf-8 -*-
 * Copyleft (C) 2015
 * "Tag bao" known as "wackonline" <bb.qnyd@gmail.com>
 * Goway is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License and GNU
 * Lesser General Public License published by the Free Software
 * Foundation, either version 3 of the License, or (at your option)
 * any later version.
 * Goway is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License and GNU Lesser General Public License
 * for more details.
 * You should have received a copy of the GNU General Public License
 * and GNU Lesser General Public License along with this program.
 * If not, see <http://www.gnu.org/licenses/>.
 */
package goway

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	E_ALL     = 0
	E_ERROR   = 1
	E_WARNING = 2
	E_STRICT  = 3
	E_NOTICE  = 4
)

type Loggers interface {
	SetLogLevel(string)
}
type logger struct {
	lvs     int
	All     bool
	Error   bool
	Warning bool
	Strict  bool
	Notice  bool
}

// Determine whether excluded set log level
// Example:
//   all := A|B|C|D|E
//   all & C not eq 0 (C in ALL)
//   all1 := A|B|C
//   all1 & D eq 0 (D not in ALL)
func (lg logger) isLv(v int) bool {
	perm := lg.lvs & v
	if perm == 0 {
		return true
	}
	return false
}

//Excluded error message is set to false
func (lg *logger) SetLogLevel(lv string) {
	str := strings.Split(lv, "|")
	if len(str) >= 1 {
		for _, v := range str {
			if v == "E_NOTICE" {
				lg.Notice = lg.isLv(E_NOTICE)
			} else if v == "E_ERROR" {
				lg.Error = lg.isLv(E_ERROR)
			} else if v == "E_WARNING" {
				lg.Warning = lg.isLv(E_WARNING)
			} else if v == "E_STRICT" {
				lg.Strict = lg.isLv(E_STRICT)
			} else {
				lg.All = lg.isLv(E_ALL)
			}
		}
	}
	lg.All = false
}

func InitLogger() Loggers {
	logs := &logger{}
	// An operation to get all the mistakes
	logs.lvs = E_ALL | E_ERROR | E_WARNING | E_STRICT | E_NOTICE
	return logs
}

func Logger() Handler {
	return func(res http.ResponseWriter, req *http.Request, c Context, log *log.Logger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		log.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

		rw := res.(ResponseWriter)
		c.Next()

		log.Printf("Completed %v %s, Content-Length: %v bytes in %v\n", rw.Status(), http.StatusText(rw.Status()), rw.Size(), time.Since(start))
	}
}