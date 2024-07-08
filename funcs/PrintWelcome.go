package netcat

import "strings"

func PrintWelcome() string {

	lines := []string{
		"Welcome to TCP-Chat!",
		"         _nnnn_",
		"        dGGGGMMb",
		"       @p~qp~~qMb",
		"       M|@||@) M|",
		"       @,----.JM|",
		"      JS^\\__/  qKL",
		"     dZP        qKRb",
		"    dZP          qKKb",
		"   fZP            SMMb",
		"   HZM            MMMM",
		"   FqM            MMMM",
		" __| \".        |\\dS\"qML",
		" |    `.       | `' \\Zq",
		"_)      \\.___.,|     .'",
		"\\____   )MMMMMP|   .'",
		"     `-'       `--'",
	}

	horizontalLines := make([]string, len(lines))

	startLine := 0
	endLine := len(lines) - 1
	for i := startLine; i <= endLine; i++ {
		if i >= 0 && i < len(lines) {
			horizontalLines[i-startLine] += lines[i]
		}
	}

	output := strings.Join(horizontalLines, "\n")
	return output
}
