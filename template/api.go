package template

type TemplateService interface {
	SayHello(input *SayHelloInput) (r *SayHelloResult, err error)
}
