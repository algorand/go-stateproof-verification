module github.com/algorand/go-stateproof-verification

go 1.17

require (
	github.com/algorand/go-algorand-sdk v1.17.0
	github.com/algorand/falcon v0.0.0-20220727072124-02a2a64c4414
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
)

replace (
	github.com/algorand/go-algorand-sdk v1.17.0 => github.com/almog-t/go-algorand-sdk 5abab22f5cc31e0bcd433615ff44cb2a50f36b13
)