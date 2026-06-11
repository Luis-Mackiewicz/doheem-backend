package http

import "strings"

func onlyDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func maskCPF(s string) string {
	d := onlyDigits(s)
	if len(d) != 11 {
		return s
	}
	return d[:3] + "." + d[3:6] + "." + d[6:9] + "-" + d[9:]
}

func maskCNPJ(s string) string {
	d := onlyDigits(s)
	if len(d) != 14 {
		return s
	}
	return d[:2] + "." + d[2:5] + "." + d[5:8] + "/" + d[8:12] + "-" + d[12:]
}

func maskDocument(s string) string {
	d := onlyDigits(s)
	switch len(d) {
	case 11:
		return maskCPF(d)
	case 14:
		return maskCNPJ(d)
	default:
		return s
	}
}

func maskPhone(s string) string {
	d := onlyDigits(s)
	if len(d) == 10 {
		return "(" + d[:2] + ") " + d[2:6] + "-" + d[6:]
	}
	if len(d) == 11 {
		return "(" + d[:2] + ") " + d[2:7] + "-" + d[7:]
	}
	return s
}

func maskCEP(s string) string {
	d := onlyDigits(s)
	if len(d) != 8 {
		return s
	}
	return d[:5] + "-" + d[5:]
}
