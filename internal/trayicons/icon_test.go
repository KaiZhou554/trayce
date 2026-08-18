package trayicons

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func makePNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func TestEncodeIconBase64(t *testing.T) {
	valid := makePNG(t)
	got := EncodeIconBase64(valid)
	if got == "" {
		t.Fatal("want non-empty base64 for valid PNG")
	}
	if !bytes.HasPrefix(valid, pngMagic) {
		t.Fatal("test png should start with png magic")
	}
}

func TestEncodeIconBase64Invalid(t *testing.T) {
	if got := EncodeIconBase64([]byte("not a png at all")); got != "" {
		t.Errorf("non-png data: got %q, want empty", got)
	}
	if got := EncodeIconBase64(nil); got != "" {
		t.Errorf("nil data: got %q, want empty", got)
	}
	if got := EncodeIconBase64([]byte{0x89, 0x50, 0x4E}); got != "" { // 截断的 PNG 头
		t.Errorf("truncated png: got %q, want empty", got)
	}
}
