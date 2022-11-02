module github.com/stefanoconti/rrc7100

go 1.12

require (
	github.com/dchote/gpio v0.0.0-20160912012454-03d78156ad1a
	github.com/kennygrant/sanitize v1.2.4
	github.com/stianeikeland/go-rpio v4.2.0+incompatible
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.1.0 // indirect
)

replace golang.org/x/net v0.1.0 => github.com/golang/net v0.1.0
