module github.com/lambda-platform/puzzle

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/labstack/echo/v4 v4.1.17
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/lambda-platform/agent v0.1.9
	github.com/lambda-platform/lambda v0.1.8
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	github.com/lambda-platform/dataform v0.1.1
	github.com/lambda-platform/datagrid v0.1.1
	github.com/lambda-platform/krud v0.1.0
)

replace github.com/lambda-platform/lambda v0.1.8 => ../lambda
replace github.com/lambda-platform/agent v0.1.9 => ../agent
replace github.com/lambda-platform/dataform v0.1.1 => ../dataform
replace github.com/lambda-platform/datagrid v0.1.1 => ../datagrid
replace github.com/lambda-platform/krud v0.1.0 => ../krud
