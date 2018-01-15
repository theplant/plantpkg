This `plantpkg` is a command line tool to generate an opinioned Go package structure that solves questions like: __How should I organize my Go code?__, __How do I manage dependencies?__

## How to use

Install the package with

```
$ go get github.com/theplant/plantpkg
```

Run this command in terminal and follow promote to generate a new package

```
$ plantpkg
```

Outputs looks like:

```
$ plantpkg
Your GOPATH: /Users/sunfmin/gopkg
✔ Generate go package: github.com/theplant/helloworld
✔ Service Name: Helloworld
Package "github.com/theplant/helloworld" generated.
cd $GOPATH/github.com/theplant/helloworld && modd; for go to the project directory and run tests.
```

The package depends on `github.com/cortesi/modd` to automatically generate protobuf go structs, and mock package. In order for `modd` command to run correctly, You will need to install:

```
$ brew install protobuf
$ go get -v github.com/cortesi/modd
$ go get -v github.com/golang/mock/mockgen
$ go get -v ./...
```

The command will generate these files inside the Go package.

```
.
├── api.go
├── config.go
├── errors.go
├── factory
│   ├── api_test.go
│   └── new.go
├── internal
│   └── impl.go
├── mock
│   └── mock.go
├── modd.conf
├── spec.pb.go
├── spec.proto
└── utils.go
```

## API functions definition

`api.go` define the outer facing API interface the package is exposing. which looks like:

```go
type CheckoutService interface {
	GiftCardApply(checkoutId string, input *GiftCardInput) (r *GiftCardResult, err error)
	ShippingAddressUpdate(checkoutId string, input *AddressInput) (r *Address, err error)
	...
}
```

The params and results are either Go primitive data types, Or protobuf defined data structs.

We intentionly limiting the package exposing API as Go interface for these reasons:

- __Easy to read__: Limit the places people come to understand the package. People can come to read your `api.go` to know what features the package provide. and don't need to read implementation details in other files of different directories.
- __Easy to switch__: Other packages who use the package can easily change to a different implementation, By switching a different `New` to construct the interface's instance.
- __Easy to extend__: We apply [Decoration](https://martinfowler.com/bliki/DecoratedCommand.html) pattern to the service to wrap more features to a basic implementation.
- __Easy to mock__: Other packages can easily pass in `mock` package instance of the package, when they don't want to test your packages implentation.
- __Easy to test__: write tests for all functions defined in `api.go`

Say with the above `CheckoutService` for example, we want to validate the Address before call `ShippingAddressUpdate` and after call it we send an email to the user to notify the change, We can do:

```go
type ValidateAndNotifyCheckoutService struct {
	basicCheckout CheckoutService
}

func (ch *ValidateAndNotifyCheckoutService) ShippingAddressUpdate(checkoutId string, input *AddressInput) (r *Address, err error) {
	err = validateAddress(input)
	if err != nil {
		return
	}

	r, err = ch.basicCheckout.ShippingAddressUpdate(checkoutId, input)
	if err != nil {
		return
	}

	err = sendNotifyEmail(input.Email)
	if err != nil {
		return
	}
	return
}
```

## API data structs definition

`spec.proto` is the place for the package to define any outer facing data structs that other package depend on this package. It is defined as [Google Protocol Buffers](https://developers.google.com/protocol-buffers/), The reason for this is:

- Can easily write an wrapper to expose the package through TCP or HTTP with preferable serialization built-in
- Can embed them into your application API protobuf structs to be part of bigger API definiation
- It has pretty good default json format generation by default
- Can still be used as standard Go structs
- All other benefits protobuf provides

## Construct new instance: factory package

`factory/new.go` is for you to construct the new instance of the interface defined in `api.go`. Where normally you pass in foreigh dependencies the package depends, like database connection, configurations, or other `plantpkg` created package instances (services).

For example:
```go
func New(db *gorm.DB, cfg *helloworld.Config, emailService email.EmailService, validationService validation.ValidationService) (service *internal.HelloworldImpl) {
	...
	return
}
```

The above `email.EmailService`, `validation.ValidationService` could be yet another `plantpkg` generated Go packages interface instance. by formalizing all the packages in `plantpkg` style. which gains:

- Easy (Unified) dependencies management
- More clarification of what package do and expose
- Easy to replace to a different implementation
- Easy to extends current implementation with wrapper, and pass them into other packages
- Unified style that the whole team would spend less time to learn new packages

In the above example, say `emailService` and `validationService` is all optional services, that with them our without them, the package always work. Then it can be implemented this way:

```go
serv := factory.New(db, cfg).EmailService(eserv).ValidationService(vserv)
```

Which is not passing optional dependencies in `New` method, But instead create a separate setter for those optional dependencies and let use set them if they need them.

## Claim errors the package will raise explicitly

`errors.go` is the place for you to define all the errors that the package can raise. It makes it clear that API could return, So that users of the package could go through them and think about how to handle those errors without looking into implementation details.

## The utility that only depends on protobuf structs

`utils.go` is for the utility methods related to this package that only depends on the protobuf data structs, which is useful when the data structs are quite complex, and you want people to easy combine, find, calulate values based on variaty of API functions returned data structs.

## Hide internal implementation

`internal/*.go` will be all your internal implementation details source code that you don't need to expose to the people who use your package. You can change the source code inside `internal` package from the bottom up, That won't effect other packages depends on the package.

## Mock support built-in

`mock` folder is a automatically generated package that mocks `api.go` definiation, and provide you a mock package by using `github.com/golang/mock/gomock`.

## A config package ties together all plantpkg packages

In your main application that depends on many `plantpkg` packages, and each of them also depends on their own `plantpkg` packages. You can setup them in a `config` package like this:

```go

func MustGetEmailService() email.EmailService {
    return emailfactory.New(...)
}

func MustGetValidationService() validation.ValdationService {
    return validationfactory.New(...)
}

func MustGetCheckoutService() checkout.CheckoutService {
    return checkoutfactory.New(...)
}

func MustGetValidateAndNotifyCheckoutService() checkout.CheckoutService {
    vs := MustGetValidationService()
    es := MustGetEmailService()
    checkout := MustGetCheckoutService()
    return vncheckoutfactory.New(..., checkout, vs, es)
}
```
