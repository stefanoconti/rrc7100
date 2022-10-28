module github.com/stefanoconti/rrc7100

go 1.11

require (
	github.com/dchote/gpio v0.0.0-20160912012454-03d78156ad1a
	github.com/dchote/gumble v0.0.0-20171217174924-fb6d7385a5cd
	github.com/kennygrant/sanitize v1.2.4
	github.com/stianeikeland/go-rpio v4.2.0+incompatible
)

require (
	github.com/dchote/go-openal v0.0.0-20171116030048-f4a9a141d372 // indirect
	github.com/dchote/gopus v0.0.0-20171117015032-b7e16762b096 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.1.0 // indirect
)

replace golang.org/x/net v0.1.0 => github.com/golang/net v0.1.0
