package matrix

import ()

func (matrix *matrix) Scroll(bitmask [][]bool, color [3]byte, y int, x int, offset int, length int) {
	s := len(bitmask[0])
	if s < length {
		matrix.DrawBitmask(bitmask, color, y, x)
		return
	}
	// Seamlessly loop the bitmap around
	offset = offset % s
	size := s - offset
	if size >= length {
		matrix.DrawBitmask(matrix.OffsetBitmask(bitmask, offset, length), color, y, x)
	} else {
		matrix.DrawBitmask(matrix.OffsetBitmask(bitmask, offset, length), color, y, x)
		matrix.DrawBitmask(matrix.OffsetBitmask(bitmask, 0, length-size), color, y, x+size)
	}
}
