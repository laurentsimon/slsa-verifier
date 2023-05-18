package static

import (
	"strings"
	"testing"
)

func Test_validVerifierIDs(t *testing.T) {
	for i := range metadatas {
		for j := 0; j < i; j++ {
			if strings.HasPrefix(metadatas[i].id, metadatas[j].id) ||
				strings.HasPrefix(metadatas[j].id, metadatas[i].id) {
				t.Errorf("%v , %v", metadatas[j].id, metadatas[i].id)
			}
		}
	}
}

func Test_differentKeyIDs(t *testing.T) {
	for i := range metadatas {
		for j := 0; j < i; j++ {
			if metadatas[i].keyid == metadatas[j].keyid {
				t.Errorf("%v == %v", metadatas[j].keyid, metadatas[i].keyid)
			}
		}
	}
}

func Test_differentKeys(t *testing.T) {
	for i := range metadatas {
		for j := 0; j < i; j++ {
			if metadatas[i].keyfile == metadatas[j].keyfile {
				t.Errorf("%v == %v", metadatas[j].keyfile, metadatas[i].keyfile)
			}
		}
	}
}
