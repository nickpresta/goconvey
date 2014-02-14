package parser

import (
	"log"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type Parser struct {
	parser func(*contract.PackageResult, string)
}

func (self *Parser) Parse(packages []*contract.Package) {
	for _, p := range packages {
		if p.Active {
			self.parser(p.Result, p.Output)

			// TODO: if contract.BuildFailure
			//           rewrite build output

			// TODO: if contract.Panicked
			//           rewrite all story stack traces
			//           (p.Result.TestResults[0].Stories[0].Assertions[0].StackTrace)

		} else {
			p.Result.Outcome = contract.Ignored
		}
		log.Printf("[%s]: %s\n", p.Result.Outcome, p.Name)
	}
}

func NewParser(helper func(*contract.PackageResult, string)) *Parser {
	self := new(Parser)
	self.parser = helper
	return self
}
