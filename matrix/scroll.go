package matrix

import ()

func (matrix *matrix) Scroll(bitmap [][]bool, color [3]byte, y int, x int, offset int, length int) {
	s := len(bitmap[0])
	if s < length {
		matrix.DrawBitmap(bitmap, color, y, x)
		return
	}
	// Seamlessly loop the bitmap around
	offset = offset % s
	size := s - offset
	if size >= length {
		matrix.DrawBitmap(matrix.OffsetBitmap(bitmap, offset, length), color, y, x)
	} else {
		matrix.DrawBitmap(matrix.OffsetBitmap(bitmap, offset, length), color, y, x)
		matrix.DrawBitmap(matrix.OffsetBitmap(bitmap, 0, length-size), color, y, x+size)
	}
}
