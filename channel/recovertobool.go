package channel

func recoverToBool(b *bool, d bool) {
	if recover() != nil {
		*b = d
	}
}
