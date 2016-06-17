package service

import (
	"testing"
	"time"

	"crypto/rand"
	"crypto/rsa"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptionParse(t *testing.T) {
	Convey("Given options", t, func() {
		publicKeyHMAC := PublicKeyHMAC([]byte("Secret"))
		privateKeyHMAC := PrivateKeyHMAC([]byte("Secret"))
		tokenExpiry := TokenExpiry(time.Second)
		refreshExpiry := RefreshTokenExpiry(time.Second)

		// --

		pk, err := rsa.GenerateKey(rand.Reader, 512)
		So(err, ShouldBeNil)

		publicKeyRSA := PublicKeyRSA(&pk.PublicKey)
		privateKeyRSA := PrivateKeyRSA(pk)

		// --

		tokenPublicKeyHMAC := TokenPublicKeyHMAC([]byte("Secret"))
		tokenPrivateKeyHMAC := TokenPrivateKeyHMAC([]byte("Secret"))
		refreshPublicKeyHMAC := RefreshTokenPublicKeyHMAC([]byte("Secret"))
		refreshPrivateKeyHMAC := RefreshTokenPrivateKeyHMAC([]byte("Secret"))

		// --

		tokenPublicKeyRSA := TokenPublicKeyRSA(&pk.PublicKey)
		tokenPrivateKeyRSA := TokenPrivateKeyRSA(pk)
		refreshPublicKeyRSA := RefreshTokenPublicKeyRSA(&pk.PublicKey)
		refreshPrivateKeyRSA := RefreshTokenPrivateKeyRSA(pk)

		Convey("When options are parsed", func() {
			opts1 := parse(
				publicKeyHMAC,
				privateKeyHMAC,
				tokenExpiry,
				refreshExpiry,
			)

			opts2 := parse(
				publicKeyRSA,
				privateKeyRSA,
			)

			opts3 := parse(
				tokenPublicKeyHMAC,
				tokenPrivateKeyHMAC,
				refreshPublicKeyHMAC,
				refreshPrivateKeyHMAC,
			)

			opts4 := parse(
				tokenPublicKeyRSA,
				tokenPrivateKeyRSA,
				refreshPublicKeyRSA,
				refreshPrivateKeyRSA,
			)

			Convey("Then the options should be valid", func() {
				So(opts1.TokenPublicKey, ShouldNotBeNil)
				So(opts1.TokenPrivateKey, ShouldNotBeNil)

				So(opts1.TokenPublicKey, ShouldEqual, opts1.RefreshPublicKey)
				So(opts1.TokenPrivateKey, ShouldEqual, opts1.RefreshPrivateKey)
				So(opts1.TokenExpiry, ShouldEqual, time.Second)
				So(opts1.RefreshExpiry, ShouldEqual, time.Second)

				//--

				So(opts2.TokenPublicKey, ShouldNotBeNil)
				So(opts2.TokenPrivateKey, ShouldNotBeNil)

				So(opts2.TokenPublicKey, ShouldEqual, opts2.RefreshPublicKey)
				So(opts2.TokenPrivateKey, ShouldEqual, opts2.RefreshPrivateKey)

				//--

				So(opts3.TokenPublicKey, ShouldNotBeNil)
				So(opts3.TokenPrivateKey, ShouldNotBeNil)
				So(opts3.RefreshPublicKey, ShouldNotBeNil)
				So(opts3.RefreshPrivateKey, ShouldNotBeNil)

				So(opts3.TokenPublicKey, ShouldNotEqual, opts3.RefreshPublicKey)
				So(opts3.TokenPrivateKey, ShouldNotEqual, opts3.RefreshPrivateKey)

				//--

				So(opts4.TokenPublicKey, ShouldNotBeNil)
				So(opts4.TokenPrivateKey, ShouldNotBeNil)
				So(opts4.RefreshPublicKey, ShouldNotBeNil)
				So(opts4.RefreshPrivateKey, ShouldNotBeNil)

				So(opts4.TokenPublicKey, ShouldNotEqual, opts4.RefreshPublicKey)
				So(opts4.TokenPrivateKey, ShouldNotEqual, opts4.RefreshPrivateKey)

			})
		})
	})
}
