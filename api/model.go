package api

type CreateSessionRequest struct {
	Account  string `json:"account"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type IdGenerator string

var (
	IdGeneratorUUID      IdGenerator = "uuid"
	IdGeneratorSnowflake IdGenerator = "snowflake"
)

type UpsertOrganizationRequest struct {
	ID          *string      `json:"id" binding:"omitempty"`
	IdGenerator *IdGenerator `json:"idGenerator" binding:"omitempty"`
	Name        string       `json:"name"`
}
