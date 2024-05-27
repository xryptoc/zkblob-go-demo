run:
	go run .

## gen ZkBlob.go
abigen: sol
	abigen --bin bin/ZkBlob.bin --abi bin/ZkBlob.abi --pkg=contract --type=ZkBlob --out=contract/ZkBlob.go
