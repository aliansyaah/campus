package domain

type Response struct {
	Status bool `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

var (
	SuccDataFound = "Data found"
	SuccCreateData = "Data successfully created"
	SuccUpdateData = "Data successfully updated"
	SuccDeleteData = "Data successfully deleted"
)