package catalog

type Experiment struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Target          string   `json:"target"`
	Description     string   `json:"description"`
	BlastRadius     string   `json:"blastRadius"`
	DefaultDuration string   `json:"defaultDuration"`
	SafetyNotes     []string `json:"safetyNotes"`
	RunbookID       string   `json:"runbookId"`
	LastOutcome     string   `json:"lastOutcome"`
	Tags            []string `json:"tags"`
}

type Service struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Owner                 string   `json:"owner"`
	Team                  string   `json:"team"`
	Tier                  string   `json:"tier"`
	Criticality           string   `json:"criticality"`
	RecentErrorRate       float64  `json:"recentErrorRate"`
	OwnershipCompleteness int      `json:"ownershipCompleteness"`
	RunbookIDs            []string `json:"runbookIds"`
	ObservabilityMaturity int      `json:"observabilityMaturity"`
	LastExperimentOutcome string   `json:"lastExperimentOutcome"`
	SecuritySensitivity   string   `json:"securitySensitivity"`
	SLO                   string   `json:"slo"`
	ExperimentIDs         []string `json:"experimentIds"`
	Repository            string   `json:"repository"`
	DashboardURL          string   `json:"dashboardUrl"`
}

type Runbook struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	AppliesTo []string `json:"appliesTo"`
	Severity  string   `json:"severity"`
	Steps     []string `json:"steps"`
	Checks    []string `json:"checks"`
	Rollback  []string `json:"rollback"`
	Links     []Link   `json:"links"`
}

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type IntegrationStatus struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Category    string   `json:"category"`
	Details     string   `json:"details"`
	Signals     []string `json:"signals"`
	DocsURL     string   `json:"docsUrl"`
	LastChecked string   `json:"lastChecked"`
}
