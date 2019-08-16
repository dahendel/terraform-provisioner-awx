package types

type Launch struct {
	CanStartWithoutUserInput bool            `json:"can_start_without_user_input"`
	PasswordsNeededToStart   []interface{}   `json:"passwords_needed_to_start"`
	AskVariablesOnLaunch     bool            `json:"ask_variables_on_launch"`
	AskTagsOnLaunch          bool            `json:"ask_tags_on_launch"`
	AskDiffModeOnLaunch      bool            `json:"ask_diff_mode_on_launch"`
	AskSkipTagsOnLaunch      bool            `json:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch       bool            `json:"ask_job_type_on_launch"`
	AskLimitOnLaunch         bool            `json:"ask_limit_on_launch"`
	AskVerbosityOnLaunch     bool            `json:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch     bool            `json:"ask_inventory_on_launch"`
	AskCredentialOnLaunch    bool            `json:"ask_credential_on_launch"`
	SurveyEnabled            bool            `json:"survey_enabled"`
	VariablesNeededToStart   []string        `json:"variables_needed_to_start"`
	CredentialNeededToStart  bool            `json:"credential_needed_to_start"`
	InventoryNeededToStart   bool            `json:"inventory_needed_to_start"`
	JobTemplateData          JobTemplateData `json:"job_template_data"`
	Defaults                 Defaults        `json:"defaults"`
}
type JobTemplateData struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Description string `json:"description"`
}
type LaunchInventory struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
type LaunchCredentials struct {
	ID              int           `json:"id"`
	Name            string        `json:"name"`
	CredentialType  int           `json:"credential_type"`
	PasswordsNeeded []interface{} `json:"passwords_needed"`
}
type Defaults struct {
	ExtraVars   string        `json:"extra_vars"`
	DiffMode    bool          `json:"diff_mode"`
	Limit       string        `json:"limit"`
	JobTags     string        `json:"job_tags"`
	SkipTags    string        `json:"skip_tags"`
	JobType     string        `json:"job_type"`
	Verbosity   int           `json:"verbosity"`
	Inventory   LaunchInventory     `json:"inventory"`
	Credentials []LaunchCredentials `json:"credentials"`
}

type LaunchRequestBody struct {
	ExtraVars map[string]interface{} `json:"extra_vars"`
}

