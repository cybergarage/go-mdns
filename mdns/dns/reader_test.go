package dns

import (
	"errors"
	"strings"
	"testing"
)

func TestReadNameWireSizeLimit(t *testing.T) {
	label := strings.Repeat("a", 63)
	valid := strings.Join([]string{label, label, label, label[:61]}, ".")
	if name, err := NewReaderWithBytes(nameToBytes(valid)).ReadName(); err != nil || name != valid {
		t.Fatalf("ReadName at the wire limit = %q, %v", name, err)
	}
	if _, err := NewReaderWithBytes(nameToBytes(valid + "a")).ReadName(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("ReadName above the wire limit = %v, want ErrInvalid", err)
	}
}
