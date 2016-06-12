package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtil(t *testing.T) {
	Convey("Given a unix path", t, func() {
		path := "/home/thatsmrtalbot/test"

		Convey("When converted to a mount string", func() {
			mount := mountString(path, "/test")

			Convey("Then the mount string should be valid", func() {
				So(mount, ShouldEqual, "/home/thatsmrtalbot/test:/test")
			})
		})
	})

	Convey("Given a windows path", t, func() {
		path := "C:\\Users\\ThatsMrTalbot\\test"

		Convey("When converted to a mount string", func() {
			mount := mountString(path, "/test")

			Convey("Then the mount string should be valid", func() {
				So(mount, ShouldEqual, "//c/Users/ThatsMrTalbot/test:/test")
			})
		})
	})

	Convey("Given a container id", t, func() {
		original := createID()

		Convey("When another id is genenerated", func() {
			new := createID()

			Convey("Then the ids should not match", func() {
				So(original, ShouldNotEqual, new)
			})
		})

		Reset(func() {
			counter = 0
		})
	})
}
