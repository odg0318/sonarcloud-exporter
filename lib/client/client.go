package client

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/sonarcloud-exporter/internal"
	sonar "github.com/whyeasy/sonarcloud-exporter/lib/sonar"
)

// Stats struct is the list of expected to results to export.
type Stats struct {
	Projects     *[]ProjectStats
	Measurements *[]MeasurementsStats
}

// ProjectStats is the struct for SonarCloud projects data we want.
type ProjectStats struct {
	Organization string
	Key          string
	Name         string
	Qualifier    string
	LastAnalysis *time.Time
}

// MeasurementsStats is the struct for SonarCloud measurements we want.
type MeasurementsStats struct {
	Key       string
	Metric    string
	Value     string
	BestValue string
}

// ExporterClient contains SonarCloud information for connecting
type ExporterClient struct {
	sqc *sonar.Client
}

// New returns a new Client connection to SonarCloud
func New(c internal.Config) *ExporterClient {
	return &ExporterClient{
		sqc: sonar.NewClient(c.Token, c.Organization, c.Metrics),
	}
}

// GetStats retrieves data from API to create metrics from.
func (c *ExporterClient) GetStats() (*Stats, error) {
	projects, err := getProjects(c.sqc)
	if err != nil {
		return nil, err
	}

	measurements, err := getMeasurements(c.sqc, projects)
	if err != nil {
		return nil, err
	}

	return &Stats{
		Projects:     projects,
		Measurements: measurements,
	}, nil
}

func getProjects(c *sonar.Client) (*[]ProjectStats, error) {
	var result []ProjectStats

	page := 1

	for {
		projects, err := c.ListProjects(&sonar.ListOptions{
			Page: page,
		})
		if err != nil {
			return nil, err
		}

		for _, project := range projects.Components {
			var lastAnalysis *time.Time

			if len(project.LastAnalysisDate) > 0 {
				parsedTime, err := time.Parse("2006-01-02T15:04:05-0700", project.LastAnalysisDate)
				if err != nil {
					return nil, err
				}

				lastAnalysis = &parsedTime
			}

			result = append(result, ProjectStats{
				Name:         project.Name,
				Qualifier:    project.Qualifier,
				Key:          project.Key,
				Organization: project.Organization,
				LastAnalysis: lastAnalysis,
			})
		}

		if len(projects.Components) == 0 {
			break
		}

		page++
	}

	log.Info("Found a total of: ", len(result), " projects")

	return &result, nil
}

func getMeasurements(c *sonar.Client, projects *[]ProjectStats) (*[]MeasurementsStats, error) {
	var result []MeasurementsStats

	for _, project := range *projects {
		data, err := c.ProjectMeasurements(project.Key)
		if err != nil {
			return nil, err
		}
		for _, measurement := range data.Component.Measures {
			result = append(result, MeasurementsStats{
				Key:       data.Component.Key,
				BestValue: strconv.FormatBool(measurement.BestValue),
				Metric:    measurement.Metric,
				Value:     measurement.Value,
			})
		}
	}

	return &result, nil
}
