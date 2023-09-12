build:
	go build -o bin/sonarcloud-exporter ./cmd/sonarcloud-exporter

deps:
	go mod verify
	go mod tidy -v

tag:
	git fetch --tags
	git tag $(TAG)
	git push origin $(TAG)

untag:
	git fetch --tags
	git tag -d $(TAG)
	git push origin :refs/tags/$(TAG)
	curl --request DELETE --header "Authorization: token ${GITHUB_TOKEN}" "https://api.github.com/repos/whyeasy/sonarcloud-exporter/releases/:release_id/$(TAG)"
