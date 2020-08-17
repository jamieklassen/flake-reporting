package flake_reporting

import (
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlakeReporting", func() {
	It("flakes", func() {
		rand.Seed(time.Now().UnixNano())
		Expect(rand.Intn(3)).To(Equal(2))
	})
	It("also flakes", func() {
		rand.Seed(time.Now().UnixNano())
		Expect(rand.Intn(3)).To(Equal(2))
	})
})
