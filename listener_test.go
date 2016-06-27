package godux_test

import (
	. "github.com/zpencerq/godux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListenerSet", func() {
	It("Signals all listeners", func() {
		l := NewListenerSet()

		var (
			listener1, listener2, listener3 Listener
		)

		listener1 = func() {}
		s1 := MakeSpy(listener1, &listener1)

		listener2 = func() {}
		s2 := MakeSpy(listener2, &listener2)

		listener3 = func() {}
		s3 := MakeSpy(listener3, &listener3)

		l.Add(&listener1)
		l.Add(&listener2)

		l.Signal()

		Expect(s1.Calls).To(HaveLen(1))
		Expect(s2.Calls).To(HaveLen(1))
		Expect(s3.Calls).To(HaveLen(0))

		l.Add(&listener3)
		l.Signal()

		Expect(s1.Calls).To(HaveLen(2))
		Expect(s2.Calls).To(HaveLen(2))
		Expect(s3.Calls).To(HaveLen(1))
	})

	It("Provides ordered equality semantics", func() {
		l := NewListenerSet()
		o := NewListenerSet()

		var (
			listener1, listener2 Listener
		)

		listener1 = func() {}
		listener2 = func() {}

		l.Add(&listener1)
		l.Add(&listener2)

		Expect(l.Equal(o)).To(BeFalse())

		o.Add(&listener1)
		o.Add(&listener2)

		Expect(l.Equal(o)).To(BeTrue())

		o = NewListenerSet()
		o.Add(&listener2)
		o.Add(&listener1)

		Expect(l.Equal(o)).To(BeFalse())

		o.Delete(&listener2)
		o.Add(&listener2)

		Expect(l.Equal(o)).To(BeTrue())
	})

	It("Delete operation returns nil if not present", func() {
		l := NewListenerSet()

		var listener1 Listener = func() {}
		l.Add(&listener1)

		Expect(l.Delete(&listener1)).To(Equal(&listener1))

		Expect(l.Delete(&listener1)).To(BeNil())
	})
})
