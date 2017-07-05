/*
 * This file is part of arduino-cli.
 *
 * arduino-cli is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 BCMI LABS SA (http://www.arduino.cc/)
 */

package formatter

import "fmt"
import "errors"

// Formatter interface represents a generic formatter. It allows to print and format Messages.
type Formatter interface {
	Format(interface{}) (string, error) // Format formats a parameter if possible, otherwise it returns an error.
	Print(interface{}) error            // Print just prints specified parameter, returns error if it is not parsable.
}

var formatters map[string]Formatter
var defaultFormatter Formatter

func init() {
	formatters = make(map[string]Formatter, 2)
	AddCustomFormatter("text", TextPrinter(0))
	AddCustomFormatter("json", JSONPrinter(1))
	defaultFormatter = formatters["text"]
}

// SetFormatter sets the defaults format to the one specified, if valid. Otherwise it returns an error.
func SetFormatter(formatName string) error {
	_, formatterExists := formatters[formatName]
	if !formatterExists {
		return fmt.Errorf("Formatter for %s format not implemented", formatName)
	}
	defaultFormatter = formatters[formatName]
	return nil
}

// IsSupported returns whether the format specified is supported or not by the current set of formatters.
func IsSupported(formatName string) bool {
	_, supported := formatters[formatName]
	return supported
}

// AddCustomFormatter adds a custom formatter to the list of available formatters of this package.
//
// If a key is already present, it is replaced and old Value is returned.
//
// If format was not already added as supported, the custom formatter is
// simply added, and oldValue returns nil.
func AddCustomFormatter(formatName string, form Formatter) Formatter {
	oldValue := formatters[formatName]
	formatters[formatName] = form
	return oldValue
}

// Format formats a message formatted using a Formatter specified by SetFormatter(...) function.
func Format(msg interface{}) (string, error) {
	if defaultFormatter == nil {
		return "", errors.New("No formatter set")
	}
	return defaultFormatter.Format(msg)
}

// Print prints a message formatted using a Formatter specified by SetFormatter(...) function.
func Print(msg interface{}) error {
	if defaultFormatter == nil {
		return errors.New("No formatter set")
	}
	return defaultFormatter.Print(msg)
}

// printFunc is the base function of all Print methods of Formatters.
//
// It can be used for an unified implementation.
func printFunc(f Formatter, msg interface{}) error {
	val, err := f.Format(msg)
	if err != nil {
		fmt.Println(val)
	}
	return err
}
