package feature

type Flags struct {
	Debug bool `short:"D" help:"Use debug features"`
}

var features *Flags

func Features() Flags {
	if features == nil {
		panic("Must set features before using the global access")
	}
	return *features
}

func SetFeatures(flags Flags) {
	// TODO: Fix this to not be as reliant on globals
	if features != nil {
		panic("can not set feature flags twice in one process")
	}
	features = &flags
}
