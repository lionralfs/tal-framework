package tal

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/lionralfs/tal-page-strategies"
)

// DeviceConfig represents a device configuration (on disk as json)
type DeviceConfig struct {
	PageStrategy string
}

// TAL is the main framework
type TAL struct {
	configPath string
}

// New creates a new TAL object
func New(configPath string) TAL {
	return TAL{configPath: configPath}
}

// GetDocType returns the doctype required by this device
// according to the page strategy in config c.
func (t TAL) GetDocType(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "doctype")
}

// GetMimeType returns the HTTP mimetype required by this device
// according to the page strategy in config c.
func (t TAL) GetMimeType(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "mimetype")
}

// GetRootHTMLTag returns the root HTML tag required by this device
// according to the page strategy in config c.
func (t TAL) GetRootHTMLTag(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "rootelement")
}

// GetDeviceHeaders returns any extra HTML content to be placed in the HTML <head> required by this device
// according to the page strategy in config c.
func (t TAL) GetDeviceHeaders(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "header")
}

// GetDeviceBody returns any extra HTML content to be placed in the HTML <body> required by this device
// according to the page strategy in config c.
func (t TAL) GetDeviceBody(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "body")
}

// GetConfigurationFromFilesystem returns a JSON formatted device configuration from the file system.
//
// Example:
//  framework.GetConfigurationFromFilesystem("default-webkit-default", "/devices")
//
// 'key' is a string representing the unique device identifier, typically "brand-model".
//
// 'subDir' is the sub-directory where the device configuration is located.
func (t TAL) GetConfigurationFromFilesystem(key string, subDir string) (string, error) {
	raw, err := ioutil.ReadFile(t.configPath + subDir + "/" + key + ".json")
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// NormaliseKeyNames replaces whitespace with underscores and lowercases all uppercase characters.
// Used to compare strings where capitalization is not guaranteed.
func NormaliseKeyNames(value string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9]")
	return strings.ToLower(re.ReplaceAllString(value, "_"))
}

// getPageStrategyElement returns an element (property) of the page strategy or the provided default value.
// The page strategy is used to determine the value of various properties to be used in the HTTP response depending
// on the class of device. This is required to support multiple specification standards such as HTML5, PlayStation 3,
// HBBTV and Maple.

// The page strategy elements are contained a separate repository referenced by a Node module, tal-page-strategies.

// Typical page strategy elements include: the HTTP header mimetype property, HTML doctype, HTML head, HTML root
// element & HTML body. For example the HTML <head> and <body> may need to contain vendor specific code/markup.

// @param string pageStrategy The page strategy used by this device.
// @param string element The page strategy property to return (Sub-directory of {{_configPath}}/config/pagestrategy/{{pageStrategy}}/{{element}}
// directory).
// @return string An element (property) of the page strategy or the default value.
func getPageStrategyElement(pageStrategy string, element string) string {
	el, err := tps.GetPageStrategyElement(pageStrategy, element)
	if err != nil {
		// after the first failed attempt to load the element, try to load the default
		_el, _err := tps.GetPageStrategyElement("default", element)
		if _err != nil {
			// TODO: throw error?
			fmt.Println(_err.Error())
		}
		return _el
	}

	return el
}
