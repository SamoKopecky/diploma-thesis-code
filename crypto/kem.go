package crypto

import (
	"github.com/SamoKopecky/pqcom/main/crypto/kem"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func init() {
	var ids []uint8
	for name, alg := range kems {
		if slices.Contains(ids, alg.Id()) {
			log.Fatal().
				Int("id", int(alg.Id())).
				Str("name", name).
				Msg("Kem algorithm Id conflict, change to a different id")
		}
		ids = append(ids, alg.Id())
	}
}

var kems = map[string]KemAlgorithm{
	"PqComKyber512": &kem.PqComKyber512{},
	"CirclKyber512": &kem.CirclKyber512{},
}

type KemAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Dec(c, sk []byte) (key []byte)
	Enc(pk []byte) (c, key []byte)
	EkLen() (ekLen int)
	Id() (id uint8)
}

type Kem struct {
	Id uint8
	F  KemAlgorithm
}

func GetKem(kemName string) Kem {
	functions := kems[kemName]
	return Kem{functions.Id(), functions}
}

func GetAllKems() []string {
	keys := make([]string, 0, len(kems))
	for k := range kems {
		keys = append(keys, k)
	}
	return keys
}
