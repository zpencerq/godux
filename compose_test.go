package godux_test

import (
	. "github.com/zpencerq/godux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compose", func() {
	It("Composes from right to left", func() {
		double := func(x int) int {
			return 2 * x
		}
		square := func(x int) int {
			return x * x
		}

		Expect(Compose(square)(5)).To(Equal(25))
		Expect(Compose(square, double)(5)).To(Equal(100))
		Expect(Compose(double, square, double)(5)).To(Equal(200))
	})

	It("Composes functions from right to left", func() {
		a := func(next func(string) string) func(string) string {
			return func(x string) string {
				return next(x + "a")
			}
		}
		b := func(next func(string) string) func(string) string {
			return func(x string) string {
				return next(x + "b")
			}
		}
		c := func(next func(string) string) func(string) string {
			return func(x string) string {
				return next(x + "c")
			}
		}

		final := func(s string) string {
			return s
		}

		Expect(Compose(a, b, c)(final).(func(string) string)("")).To(Equal("abc"))
		Expect(Compose(b, c, a)(final).(func(string) string)("")).To(Equal("bca"))
		Expect(Compose(c, a, b)(final).(func(string) string)("")).To(Equal("cab"))
	})

	It("Can be seeded with multiple arguments", func() {
		square := func(x int) int {
			return x * x
		}

		add := func(x, y int) int {
			return x + y
		}

		Expect(Compose(square, add)(1, 2)).To(Equal(9))
	})

	It("Returns the first given argument if given no functions", func() {
		Expect(Compose()(1, 2)).To(Equal(1))
		Expect(Compose()(3)).To(Equal(3))
		Expect(Compose()(nil)).To(BeNil())
		Expect(func() {
			Compose()()
		}).To(Panic())
	})

	It("Returns the first function if given only one", func() {
		fn := func() int { return 3 }

		Expect(Compose(fn)()).To(Equal(fn()))
	})
})
