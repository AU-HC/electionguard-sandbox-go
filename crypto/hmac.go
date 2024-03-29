package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"electionguard-sandbox-go/models"
	mod "electionguard-sandbox-go/modular_arithmetic"
	"encoding/binary"
	"reflect"
)

var nilType = reflect.TypeOf(nil)
var stringType = reflect.TypeOf("")
var intType = reflect.TypeOf(1)
var bigIntType = reflect.TypeOf(models.BigInt{})
var bigIntPointerType = reflect.TypeOf(&models.BigInt{})
var fileType = reflect.TypeOf(([]byte)(nil))

func HMAC(key models.BigInt, domainSeparator byte, a ...interface{}) *models.BigInt {
	mac := hmac.New(sha256.New, key.Bytes())

	// Add the domain separator first
	mac.Write([]byte{domainSeparator})

	// Then append the message (i.e. what is to be hashed)
	for _, x := range a {
		var toBeHashed []byte

		switch reflect.TypeOf(x) {
		case intType:
			// Type cast and create byte array which the number is to be stored in
			xInt, _ := x.(int)

			// We know that all small integers in ElectionGuard is smaller than 2^31, therefore we can typecast to uint32
			toBeHashed = make([]byte, 4)
			binary.BigEndian.PutUint32(toBeHashed, uint32(xInt))

		case stringType:
			// Type cast (strings are already utf8-encoded in Golang)
			xString, _ := x.(string)

			// Pad the string as a byte with four empty bytes
			pad := make([]byte, 4)
			binary.BigEndian.PutUint32(pad, uint32(len(xString)))
			toBeHashed = append(pad, []byte(xString)...)

		case bigIntType:
			bigInt := x.(models.BigInt)
			toBeHashed = bigInt.Bytes()

		case bigIntPointerType:
			bigIntPointer := x.(*models.BigInt)
			toBeHashed = bigIntPointer.Bytes()

		case fileType:
			file, _ := x.([]byte)

			pad := make([]byte, 4)
			binary.BigEndian.PutUint32(pad, uint32(len(file)))
			toBeHashed = append(pad, file...)

		default:
			panic("unknown type for hmac")
		}

		mac.Write(toBeHashed)
	}

	hash := models.MakeBigIntFromByteArray(mac.Sum(nil))
	return mod.ModQ(hash)
}
