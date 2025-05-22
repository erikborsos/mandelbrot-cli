package kitty

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func Kitty(pngBytes []byte, cols, rows int) (string, error) {
	const chunkSize = 16384
	var sb strings.Builder
	data := pngBytes
	prefix := fmt.Sprintf("\x1b_Ga=T,f=100,m=1,c=%d,r=%d;", cols, rows)
	for len(data) > chunkSize {
		chunk := data[:chunkSize]
		data = data[chunkSize:]
		sb.WriteString(prefix)
		sb.WriteString(base64.StdEncoding.EncodeToString(chunk))
		sb.WriteString("\x1b\\")
		prefix = "\x1b_Ga=T,f=100,m=1;"
	}
	sb.WriteString("\x1b_Ga=T,f=100,m=0;")
	sb.WriteString(base64.StdEncoding.EncodeToString(data))
	sb.WriteString("\x1b\\")
	return sb.String(), nil
}

func KittyClearImages() {
	fmt.Print("\x1b_Ga=d,d=A,q=2;\x1b")
}
