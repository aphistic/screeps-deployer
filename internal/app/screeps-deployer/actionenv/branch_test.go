package actionenv

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type branchSuite struct{}

func (s *branchSuite) TestParseHead(t sweet.T) {
	Expect(ParseBranch("refs/heads/testbranch")).To(Equal("testbranch"))
	Expect(ParseBranch("refs/heads/master")).To(Equal("master"))
}
