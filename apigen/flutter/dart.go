package flutter

import (
	"reflect"

	"github.com/rubikorg/blocks/apigen/sdkdist"
	"github.com/rubikorg/rubik"
)

// DartGenPlugin generates dart client code for Flutter development
// from rubik routes
type DartGenPlugin struct{}

// OnPlug implements rubik plugin interface
func (DartGenPlugin) OnPlug(app *rubik.App) error {
	var _ = sdkdist.TransformTree(app.RouteTree, getDartTypeEquivalent)
	return nil
}

// Name that is displayed while running
func (DartGenPlugin) Name() string {
	return "Flutter/Dart Client SDK Generator"
}

// RunID specifies the running identifier for this plugin
func (DartGenPlugin) RunID() string {
	return "bind_flutter"
}

// TODO: need to handle nested struct inside the type
func getDartTypeEquivalent(goType string, value reflect.StructField) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float", "float32", "float64",
		"complex32", "complex64":
		return "number"
	case "bool":
		return "boolean"
	default:
		// ignore!
		return "-1"
	}
}

func init() {
	rubik.Plug(DartGenPlugin{})
}
