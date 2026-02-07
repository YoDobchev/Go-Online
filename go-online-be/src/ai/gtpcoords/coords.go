package gtpcoords

import (
	"fmt"
	"strconv"
	"strings"
)

const letters = "ABCDEFGHJKLMNOPQRST"

func ToGTP(size, row, col int) (string, error) {
	if row < 0 || col < 0 || row >= size || col >= size {
		return "", fmt.Errorf("out of bounds")
	}
	return fmt.Sprintf("%c%d", letters[col], size-row), nil
}

func FromGTP(size int, s string) (row, col int, pass bool, err error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "PASS" {
		return 0, 0, true, nil
	}
	if s == "RESIGN" {
		return 0, 0, false, fmt.Errorf("resign")
	}

	if len(s) < 2 {
		return 0, 0, false, fmt.Errorf("bad coord")
	}
	c := s[0]
	n, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, 0, false, err
	}

	col = strings.IndexByte(letters, c)
	if col < 0 {
		return 0, 0, false, fmt.Errorf("bad col")
	}

	row = size - n
	if row < 0 || row >= size {
		return 0, 0, false, fmt.Errorf("bad row")
	}

	return row, col, false, nil
}
