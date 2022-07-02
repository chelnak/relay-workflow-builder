tag:
	@git tag -a $(version) -m "$(version)"
	@git push --follow-tags

lint:
	@golangci-lint run ./...

