package paramlist_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestParamlist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Paramlist Suite")
}
