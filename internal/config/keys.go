package config

type Keys struct {
	ForceQuit              []string `toml:"force_quit"`
	Exit                   []string `toml:"exit"`
	ShowMatrix             []string `toml:"show_matrix"`
	Help                   []string `toml:"help"`
	Submit                 []string `toml:"submit"`
	Up                     []string `toml:"up"`
	Down                   []string `toml:"down"`
	Left                   []string `toml:"left"`
	Right                  []string `toml:"right"`
	RotateCounterClockwise []string `toml:"rotate_counter_clockwise"`
	RotateClockwise        []string `toml:"rotate_clockwise"`
}

func DefaultKeys() *Keys {
	return &Keys{
		ForceQuit:              []string{"ctrl+c"},
		Exit:                   []string{"esc"},
		ShowMatrix:             []string{"m"},
		Help:                   []string{"?"},
		Submit:                 []string{" ", "enter"},
		Up:                     []string{"z"},
		Down:                   []string{"s"},
		Left:                   []string{"q"},
		Right:                  []string{"d"},
		RotateCounterClockwise: []string{"a"},
		RotateClockwise:        []string{"e"},
	}
}
