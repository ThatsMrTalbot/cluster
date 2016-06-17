package client

import (
	"testing"

	"crypto/rand"
	"crypto/rsa"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptionParse(t *testing.T) {
	Convey("Given options", t, func() {
		publicKeyHMAC := PublicKeyHMAC([]byte("Secret"))
		pk, err := rsa.GenerateKey(rand.Reader, 512)
		So(err, ShouldBeNil)

		publicKeyRSA := PublicKeyRSA(&pk.PublicKey)

		Convey("When options are parsed", func() {
			opts1 := parse(
				publicKeyHMAC,
			)

			opts2 := parse(
				publicKeyRSA,
			)

			Convey("Then the options should be valid", func() {
				So(opts1.PublicKey, ShouldNotBeNil)
				So(opts2.PublicKey, ShouldNotBeNil)
			})
		})
	})
}
