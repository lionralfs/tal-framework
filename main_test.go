package tal

import (
	j "encoding/json"
	"strings"
	"testing"
)

func getGenericDevice1Config() DeviceConfig {
	framework := New("./fixtures/")
	json, _ := framework.GetConfigurationFromFilesystem("generic-tv1", "deviceconfig")

	var deviceConfigParsed DeviceConfig
	errr := j.Unmarshal([]byte(json), &deviceConfigParsed)
	if errr != nil {
		panic(errr)
	}
	return deviceConfigParsed
}

// Generic TV1 Device has no Headers
func TestGetDeviceHeaders(t *testing.T) {
	framework := New("./fixtures/")
	headers := strings.TrimSpace(framework.GetDeviceHeaders(getGenericDevice1Config()))

	if headers != "" {
		t.Error("The device headers are not empty. It contains: " + headers)
	}
}

// Generic TV1 Device has no body
func TestGetDeviceBody(t *testing.T) {
	framework := New("./fixtures/")
	body := strings.TrimSpace(framework.GetDeviceBody(getGenericDevice1Config()))

	if body != "" {
		t.Error("The device body is not empty. It contains: " + body)
	}
}

// Generic TV1 Device has default Mime type
func TestGetMimeType(t *testing.T) {
	framework := New("./fixtures/")
	mimeType := strings.TrimSpace(framework.GetMimeType(getGenericDevice1Config()))

	if mimeType != "text/html" {
		t.Error("The mime type is not text/html. The value was " + mimeType)
	}
}

// Generic TV1 Device has default Root element
func TestGetRootHTMLTag(t *testing.T) {
	framework := New("./fixtures/")
	rootElement := strings.TrimSpace(framework.GetRootHTMLTag(getGenericDevice1Config()))

	if rootElement != "<html>" {
		t.Error("The root element is not '<html>'. The value was " + rootElement)
	}
}

// Generic TV1 Device has default Doc type
func TestGetDocType(t *testing.T) {
	framework := New("./fixtures/")
	docType := strings.TrimSpace(framework.GetDocType(getGenericDevice1Config()))

	if docType != "<!DOCTYPE html>" {
		t.Error("The device does not have the default doc type (<!DOCTYPE html>). The value was " + docType)
	}
}

// Normalise key names replaces special characters with underscores
func TestNormaliseKeyNames1(t *testing.T) {
	v := NormaliseKeyNames("one$two(three")
	e := "one_two_three"

	if v != e {
		t.Error("Expected "+e+", got ", v)
	}
}

// Normalise key names replaces upper case to lower case
func TestNormaliseKeyNames2(t *testing.T) {
	v := NormaliseKeyNames("one_TWO_Three")
	e := "one_two_three"

	if v != e {
		t.Error("Expected "+e+", got ", v)
	}
}
