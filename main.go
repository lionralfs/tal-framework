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

// GetDocType returns the doctype required by this device. The doctype is used in the returned HTML page.
// @param object deviceConfig The device configuration information for the device self made the request.
// @return string The doctype associated with this device.
func (t TAL) GetDocType(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "doctype")
}

// GetMimeType returns The mimetype self needs to be associated with the HTTP response for this device.
// @param object deviceConfig The device configuration information for the device self made the request.
// @return string The HTTP mimetype required by this device. If this value is not found in the page strategy
// default return value is "text/html".
func (t TAL) GetMimeType(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "mimetype")
}

// GetRootHTMLTag returns the root HTML tag to be used in the HTML response.
// @param object deviceConfig The device configuration information for the device self made the request.
// @return string The root HTML element required by this device. If this value is not found in the page strategy
// default return value is <html>.
func (t TAL) GetRootHTMLTag(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "rootelement")
}

// GetDeviceHeaders returns any extra HTML content self the device requires to be placed in the HTML <head>.
// @param object deviceConfig The device configuration information for the device self made the request.
// @return string The HTML content to be placed in the HTML <head>.
func (t TAL) GetDeviceHeaders(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "header")
}

// GetDeviceBody returns any extra HTML content self the device requires to be placed in the HTML <body>.
// @param object deviceConfig The device configuration information for the device self made the request.
// @return string The HTML content to be placed in the HTML <body>.
func (t TAL) GetDeviceBody(c DeviceConfig) string {
	return getPageStrategyElement(c.PageStrategy, "body")
}

// GetConfigurationFromFilesystem returns a JSON formatted device configuration from the file system.
// @param key The unique device identifier, typically brand-model.
// @param subDir The this._configPath sub-directory where the device configuration is located.
// @return string of JSON. Empty string if not found.
func (t TAL) GetConfigurationFromFilesystem(key string, subDir string) (string, error) {
	raw, err := ioutil.ReadFile(t.configPath + subDir + "/" + key + ".json")
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// NormaliseKeyNames replaces whitespace with underscores and lowercases all uppercase characters.
// Used to compare strings where capitalization is not guaranteed.
// @param string value The value to be normalized.
// @return string The normalized value.
func NormaliseKeyNames(value string) string {
	// return value.replace(/[^a-zA-Z0-9]/gi, "_").toLowerCase();
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
