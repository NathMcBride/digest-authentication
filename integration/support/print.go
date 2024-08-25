package support

import "fmt"

type Color int

const (
	Red Color = iota
	Green
	Blue
	Cyan
	Orange
	Yellow
)

const csi = "\033["
const rgbForeground = "38;2;"
const sgrCode = csi + rgbForeground
const reset = "\033[0m"

var rgbMap = map[Color]string{
	Red:    foreGround(255, 0, 28),
	Green:  foreGround(0, 255, 49),
	Blue:   foreGround(26, 0, 248),
	Cyan:   foreGround(0, 255, 253),
	Orange: foreGround(255, 167, 43),
	Yellow: foreGround(253, 255, 57),
}

func foreGround(r int, g int, b int) string {
	return fmt.Sprintf("%s%d;%d;%dm", sgrCode, r, g, b)
}

func CSprintf(c Color, s string, a ...interface{}) string {
	colorCode, exists := rgbMap[c]
	if !exists {
		return fmt.Sprintf(s, a...)
	}

	return fmt.Sprintf(colorCode+s+reset, a...)
}
