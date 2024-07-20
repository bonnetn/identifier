package id

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"github.com/google/uuid"
)

const (
	uuidLength        = 16
	versionLength     = 8
	canonicalIDLength = uuidLength + versionLength
	fullIDLength      = canonicalIDLength + sha256.Size
	encodedIDLength   = (fullIDLength*8 + 5 - 1) / 5
)

// For testing purposes
var randomUUID = uuid.NewV7

type Generator struct {
	seedHash [sha256.Size]byte
}

func NewGenerator(seed string) *Generator {
	return &Generator{
		seedHash: sha256.Sum256([]byte(seed)),
	}
}

func (g *Generator) Generate() (ID, error) {
	if g.seedHash == [32]byte{} {
		panic("seed is empty")
	}

	id, err := randomUUID()
	if err != nil {
		return ID{}, err
	}

	return ID{
		id:      id,
		version: 1,
		hmac:    g.hmac(canonical(id, 1)),
	}, nil
}

func (g *Generator) Verify(id ID) bool {
	want := g.hmac(canonical(id.id, id.version))
	return subtle.ConstantTimeCompare(
		id.hmac[:],
		want[:],
	) == 1
}

func (g *Generator) GenerateNewVersion(id ID, version uint64) ID {
	return ID{
		id:      id.id,
		version: version,
		hmac:    g.hmac(canonical(id.id, version)),
	}
}

func (g *Generator) hmac(canonical [canonicalIDLength]byte) [sha256.Size]byte {
	var key [sha256.Size + canonicalIDLength]byte
	copy(key[:], g.seedHash[:])
	copy(key[sha256.Size:], canonical[:])
	return sha256.Sum256(key[:])
}

type ID struct {
	id      uuid.UUID
	version uint64
	hmac    [sha256.Size]byte
}

func (i *ID) ID() uuid.UUID {
	return i.id
}

func (i *ID) Version() uint64 {
	return i.version
}

func (i *ID) HMAC() [sha256.Size]byte {
	return i.hmac
}

func (i *ID) String() string {
	b := i.Bytes()
	return toBase32String(b)
}

func (i *ID) Bytes() [fullIDLength]byte {
	var result [fullIDLength]byte
	copy(result[:], i.id[:])
	binary.LittleEndian.PutUint64(result[uuidLength:], i.version)
	copy(result[canonicalIDLength:], i.hmac[:])
	return result
}

func canonical(id uuid.UUID, version uint64) [canonicalIDLength]byte {
	var result [16 + 8]byte
	copy(result[:], id[:])
	binary.LittleEndian.PutUint64(result[16:], version)
	return result
}

var mapping = []byte("0123456789ABCDEFGHJKMNPQRSTVWXYZ")
