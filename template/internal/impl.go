package internal

import (
	"github.com/jinzhu/gorm"
	"github.com/theplant/plantpkg/template"
)

type TemplateImpl struct {
	db             *gorm.DB
	cfg            *template.Config
	optionalParam1 string
	optionalParam2 string
}

func New(db *gorm.DB, cfg *template.Config) (r *TemplateImpl, err error) {
	r = &TemplateImpl{
		db:  db,
		cfg: cfg,
	}
	return
}

func (impl *TemplateImpl) OptionalParam1(param1 string) *TemplateImpl {
	impl.optionalParam1 = param1
	return impl
}

func (impl *TemplateImpl) OptionalParam2(param2 string) *TemplateImpl {
	impl.optionalParam2 = param2
	return impl
}

func (impl *TemplateImpl) SayHello(input *template.SayHelloInput) (r *template.SayHelloResult, err error) {
	r = &template.SayHelloResult{Result: "result"}
	return
}
