package service

const (
	prevBody = "prev"
	nextBody = "next"
)

// Определяет предыдущий и следующий оффсеты. В случае, если таковых нет, возвращается -1
func getPrevAndNextOffset(total, limit, offset int) (prevOffset, nextOffset int) {
	if nextOffset = offset + limit; nextOffset >= total {
		nextOffset = -1
	}

	if offset < 0 {
		prevOffset = -1
	} else {
		if prevOffset = offset - limit; prevOffset < 0 {
			if offset == 0 {
				prevOffset = -1
			} else {
				prevOffset = 0
			}
		}
	}

	return
}
