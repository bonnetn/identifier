package id

func toBase32(b [fullIDLength]byte) [encodedIDLength]byte {
	const mask5 = 0b11111

	var (
		result [encodedIDLength]byte
		q      uint16
		l      int
		j      int
	)
	for i := 0; i < fullIDLength; i++ {
		q <<= 8
		q |= uint16(b[i])
		l += 8

		for l >= 5 {
			l -= 5
			c := (q >> l) & mask5

			result[j] = mapping[c]
			j++
		}
	}

	if l > 0 {
		l -= 5
		if l < 0 {
			l = 0
		}
		c := (q >> l) & mask5

		result[j] = mapping[c]
		j++
	}
	return result
}
func toBase32_temp(b [fullIDLength]byte) [encodedIDLength]byte {
	const mask5 = 0b11111

	var (
		result [encodedIDLength]byte
		q      uint16
		l      int
		j      = len(result) - 1
	)
	for i := fullIDLength - 1; i >= 0; i-- {
		q |= uint16(b[i]) << l
		l += 8

		for l >= 5 {
			l -= 5
			c := q & mask5
			q >>= 5

			result[j] = mapping[c]
			j--
		}
	}

	if l > 0 {
		c := q & mask5

		result[j] = mapping[c]
		j--
	}
	return result
}

func toBase32String(b [fullIDLength]byte) string {
	result := toBase32(b)
	return string(result[:])
}
