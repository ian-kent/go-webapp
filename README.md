go-webapp
=========

Demonstration Go web application, built on [this boilerplate](https://github.com/ian-kent/go-angularjs-jquery-bootstrap-boilerplate)

* [nosurf](https://github.com/justinas/nosurf) CSRF protection
* [gorilla/schema](https://github.com/gorilla/schema) form decoding
* [validator.v5](https://github.com/bluesuncorp/validator) form validation
* [htmlform](https://github.com/ian-kent/htmlform) form rendering
* [gorilla/sessions](https://github.com/gorilla/sessions) session handling
* [gofigure](https://github.com/ian-kent/gofigure) configuration
* Example data storage
* Registration, login and logout pages
* Request timeout handling
* Per-request logging

## Building go-webapp

```bash
go generate ./... && go build . && ./go-webapp
```

## License

Copyright ©‎ 2015, Ian Kent (http://iankent.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
