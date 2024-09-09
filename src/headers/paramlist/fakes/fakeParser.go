package fakes

type FakeParser struct {
	callCount        int
	parseListReturns struct {
		parsed map[string]string
		err    error
	}
	parseListArgsForCall []struct {
		auth   string
		prefix string
	}
}

func (p *FakeParser) ParseListCallCount() int {
	return p.callCount
}

func (p *FakeParser) ParseListReturns(parsed map[string]string, err error) {
	p.parseListReturns = struct {
		parsed map[string]string
		err    error
	}{
		parsed,
		err,
	}
}

func (p *FakeParser) ParseListArgsForCall(i int) (string, string) {
	args := p.parseListArgsForCall[i]
	return args.auth, args.prefix
}

func (p *FakeParser) ParseList(auth string, prefix string) (map[string]string, error) {
	p.callCount++
	p.parseListArgsForCall = append(p.parseListArgsForCall, struct {
		auth   string
		prefix string
	}{auth, prefix})
	return p.parseListReturns.parsed, p.parseListReturns.err
}
