package sanitizer

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"google.golang.org/api/gmail/v1"
)

var (
	spammers []string
)

func init() {
	if bytes, err := os.ReadFile("spammers.txt"); err == nil {
		content := string(bytes)
		lines := strings.Split(content, "\n")
		spammers = make([]string, len(lines))

		for i, spammer := range lines {
			if len(strings.Trim(spammer, " ")) == 0 {
				continue
			}
			spammers[i] = strings.Trim(spammer, " ")
		}
		// fmt.Printf("spammers: %v", spammers)
	}
}

func isSpammer(address string) bool {

	for _, spammer := range spammers {
		// spammer = fmt.Sprintf("@%s", spammer)
		if strings.Contains(address, spammer) {
			// fmt.Printf("%s -- %s= %v", address, spammer, strings.Contains(address, spammer))
			return true
		}
	}
	return false
}

func getSender(headers []*gmail.MessagePartHeader) string {
	for _, h := range headers {
		if h.Name == "From" {
			return h.Value
		}
	}
	return "Not Found"
}

func decodeBody(msg string) string {
	if msg == "" {
		return ""
	}
	fmt.Println("==========================body==================================")
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(msg)))

	n, _ := base64.StdEncoding.Decode(base64Text, []byte(msg))
	fmt.Println("base64Text:", string(base64Text[:n]))
	return string(base64Text[:n])
}
