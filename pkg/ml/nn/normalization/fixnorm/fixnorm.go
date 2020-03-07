// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fixnorm

import (
	"io"
	"log"
	"saientist.dev/spago/pkg/ml/ag"
	"saientist.dev/spago/pkg/ml/nn"
)

var _ nn.Model = &Model{}

// Reference: "Improving Lexical Choice in Neural Machine Translation" by Toan Q. Nguyen and David Chiang (2018)
// (https://arxiv.org/pdf/1710.01329.pdf)
type Model struct{}

func New() *Model {
	return &Model{}
}

func (m *Model) ForEachParam(callback func(param *nn.Param)) {
	nn.ForEachParam(m, callback)
}

func (m *Model) Serialize(w io.Writer) (int, error) {
	return nn.Serialize(m, w)
}

func (m *Model) Deserialize(r io.Reader) (int, error) {
	return nn.Deserialize(m, r)
}

type Processor struct {
	opt   []interface{}
	model *Model
	g     *ag.Graph
}

func (m *Model) NewProc(g *ag.Graph, opt ...interface{}) nn.Processor {
	p := &Processor{
		model: m,
		opt:   opt,
		g:     g,
	}
	p.init(opt)
	return p
}

func (p *Processor) init(opt []interface{}) {
	if len(opt) > 0 {
		log.Fatal("fixnorm: invalid init options")
	}
}

func (p *Processor) Model() nn.Model       { return p.model }
func (p *Processor) Graph() *ag.Graph      { return p.g }
func (p *Processor) RequiresFullSeq() bool { return true }
func (p *Processor) Reset()                { p.init(p.opt) }

func (p *Processor) Forward(xs ...ag.Node) []ag.Node {
	ys := make([]ag.Node, len(xs))
	eps := p.g.NewScalar(1e-10)
	for i, x := range xs {
		norm := p.g.Sqrt(p.g.ReduceSum(p.g.Square(x)))
		ys[i] = p.g.DivScalar(x, p.g.AddScalar(norm, eps))
	}
	return ys
}
