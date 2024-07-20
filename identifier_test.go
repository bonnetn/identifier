package id

import (
	"github.com/google/uuid"
	"testing"
)

var stubUUID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var stubUUIDGeneration = func() (uuid.UUID, error) { return stubUUID, nil }

const seed = "seed"

func Test(t *testing.T) {
	randomUUID = stubUUIDGeneration

	g := NewGenerator(seed)
	i, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v; want nil", err)
	}

	if g.Verify(i) != true {
		t.Errorf("Generate().Verify() = %v; want true", g.Verify(i))
	}

	if i.Version() != 1 {
		t.Errorf("Generate().Version() = %d; want 1", i.version)
	}

	wantUUID := [16]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}
	if i.ID() != wantUUID {
		t.Errorf("Generate().ID() = %v; want %v", i.id, wantUUID)
	}

	wantHMAC := [32]byte{209, 97, 188, 90, 60, 178, 13, 69, 185, 104, 232, 137, 18, 201, 35, 66, 123, 23, 148, 225, 1, 117, 236, 35, 67, 137, 187, 195, 118, 233, 145, 122}
	if i.HMAC() != wantHMAC {
		t.Errorf("Generate().HMAC() = %v; want %v", i.hmac, wantHMAC)
	}

	b := i.Bytes()
	wantBytes := [...]byte{
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // UUID
		1, 0, 0, 0, 0, 0, 0, 0, // Version
		209, 97, 188, 90, 60, 178, 13, 69, 185, 104, 232, 137, 18, 201, 35, 66, 123, 23, 148, 225, 1, 117, 236, 35, 67, 137, 187, 195, 118, 233, 145, 122, // HMAC
	}
	if b != wantBytes {
		t.Errorf("ID.Bytes() = %v; want %v", b, wantBytes)
	}

	const wantStr = "248H248H248H248H248H248H240G00000000006HC6Y5MF5J1N2VJT78H49CJ8T2FCBS9R81EQP26GW9QF1QDTCHFT"
	if i.String() != wantStr {
		t.Errorf("Generate().String() = %q; want %q", i.String(), wantStr)
	}

	id2 := g.GenerateNewVersion(i, 42)
	if id2.version != 42 {
		t.Errorf("GenerateNewVersion().version = %d; want 42", id2.version)
	}

	if g.Verify(id2) != true {
		t.Errorf("Generate().Verify() = %v; want true", g.Verify(id2))
	}

	if id2.id != wantUUID {
		t.Errorf("GenerateNewVersion().id = %v; want %v", id2.id, wantUUID)
	}

	wantHMAC2 := [32]byte{19, 255, 202, 14, 243, 102, 35, 125, 133, 116, 107, 201, 41, 242, 149, 3, 246, 220, 166, 145, 4, 192, 188, 222, 221, 97, 251, 192, 33, 36, 219, 184}
	if id2.hmac != wantHMAC2 {
		t.Errorf("GenerateNewVersion().hmac = %v; want %v", id2.hmac, wantHMAC2)
	}

	const wantStr2 = "248H248H248H248H248H248H24N000000000000KZZ50XWV64DYRAX3BS4MZ5583YVEAD484R2YDXQB1ZF02296VQR"
	if id2.String() != wantStr2 {
		t.Errorf("GenerateNewVersion().String() = %v; want %q", id2.String(), wantStr2)
	}
}

func TestGenerator_Verify(t *testing.T) {
	randomUUID = stubUUIDGeneration
	g := NewGenerator(seed)
	id, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v; want nil", err)
	}

	if g.Verify(id) != true {
		t.Errorf("Verify() = %v; want true", g.Verify(id))
	}

	id.hmac[0]++
	if g.Verify(id) != false {
		t.Errorf("Verify() = %v; want false", g.Verify(id))
	}
}

func BenchmarkNewGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewGenerator(seed)
	}
}

func BenchmarkGenerate(b *testing.B) {
	randomUUID = stubUUIDGeneration
	g := NewGenerator(seed)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := g.Generate()
		if err != nil {
			b.Fatalf("Generate() error = %v; want nil", err)
		}
	}
}
func BenchmarkGenerateNewVersion(b *testing.B) {
	id := setup(b)
	g := NewGenerator(seed)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = g.GenerateNewVersion(id, 42)
	}
}

func BenchmarkID(b *testing.B) {
	id := setup(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.ID()
	}
}

func BenchmarkVersion(b *testing.B) {
	id := setup(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.Version()
	}
}

func BenchmarkHMAC(b *testing.B) {
	id := setup(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.HMAC()
	}
}

func BenchmarkString(b *testing.B) {
	id := setup(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.String()
	}
}

func BenchmarkBytes(b *testing.B) {
	id := setup(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.Bytes()
	}
}

func setup(b *testing.B) ID {
	b.Helper()

	randomUUID = stubUUIDGeneration
	g := NewGenerator(seed)
	id, err := g.Generate()
	if err != nil {
		b.Fatalf("Generate() error = %v; want nil", err)
	}

	return id
}
