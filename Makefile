.PHONY: run
run:
	$(info Running...)
	go run ./cmd/a2tonium -config=local
generatePublicKey:
	$(info Public Key Generation...)
	go run ./cmd/a2tonium -config=local --generatePublicKey
