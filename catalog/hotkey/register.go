package hotkey

import "github.com/lobehub/lobe-mygo/catalog/registry"

func init() {
	registry.Add("hotkey/basic", Basic)
	registry.Add("hotkey/variants", Variants)
	registry.Add("hotkey/sizes", Sizes)
}
