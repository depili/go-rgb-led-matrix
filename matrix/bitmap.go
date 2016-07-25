package matrix

import ()

func (matrix *matrix) DrawBitmap(bitmap [][]bool, color [3]byte, r int, c int) {
	for y, row := range bitmap {
		for x, b := range row {
			if b {
				matrix.SetPixel(r+y, c+x, color)
			}
		}
	}
}

func (matrix *matrix) ColorBitmap(bits [][]bool, color [3]byte) [][][3]byte {
	bitmap := make([][][3]byte, len(bits))
	for r, row := range bits {
		bitmap[r] = make([][3]byte, len(row))
		for c, b := range row {
			if b {
				bitmap[r][c] = color
			} else {
				bitmap[r][c] = [3]byte{0, 0, 0}
			}
		}
	}
	return bitmap
}

func (matrix *matrix) OffsetBitmap(bitmap [][]bool, offset int, length int) [][]bool {
	ret := make([][]bool, len(bitmap))
	for r, row := range bitmap {
		end := offset + length
		if end > len(row) {
			end = len(row)
		}
		ret[r] = make([]bool, end-offset)
		for c, b := range row[offset:end] {
			ret[r][c] = b
		}
	}
	return ret
}
