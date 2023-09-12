# sonarcloud-exporter

A Prometheus Exporter for SonarCloud

Currently this exporter retrieves the following metrics:

- Project info within a given organization.
- Lines of Code within a project.
- Code Coverage of a project.
- Amount of bugs within a project.
- Amount of Code smells within a project.
- Amount of vulnerabilities within a project.

## Requirements

### Required

Provide your SonarCloud organization; `-organization <string>` or as env variable `SC_ORGANIZATION`.

Provide a SonarCloud Access Token to access the API; `-token <string>` or as env variable `SC_TOKEN`.

### Optional

Provide a list of metric names with comma separated; `-metrics <string>` or as env variable `SC_METRICS`.

Change listening port of the exporter; `-listenAddress <string>` or as env variable `SC_LISTEN_ADDRESS`. Default = `8080`

Change listening path of the exporter; `-listenPath <string>` or as env variable `SC_LISTEN_PATH`. Default = `/metrics`
