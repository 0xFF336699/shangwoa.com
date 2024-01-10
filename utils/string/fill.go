package string

func StringPrefixFillToLength(prefix, str string, l int) string {
	for len(str) < l{
		str = prefix + str
	}
	return str
}

func StringSuffixFillToLength(suffix, str string, l int) string {
	for len(str) < l{
		str = suffix + str
	}
	return str
}
