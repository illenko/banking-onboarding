package util

import "strings"

func MaskPan(pan string) string {
	if len(pan) <= 10 {
		return pan
	}
	return pan[0:6] + strings.Repeat("*", len(pan)-10) + pan[len(pan)-4:]
}
