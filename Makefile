proto-gen:
	rm -rf ./gen
	buf lint && buf generate
