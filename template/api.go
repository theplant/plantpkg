package template

type HelloService interface {
	SayHello(input *SayHelloInput) (r *SayHelloResult, err error)
}
