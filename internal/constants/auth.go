package constants

const (
	// Token types
	TokenTypeDefault      = "DEFAULT"
	TokenTypeNutritionist = "NUTRITIONIST"

	// Permissions
	PermissionListDiet   = "list_diet"
	PermissionCreateDiet = "create_diet"
	PermissionUpdateDiet = "update_diet"
	PermissionUploadFile = "upload_file"
)

// GetPermissionsByUserType returns the permissions for a given user type
func GetPermissionsByUserType(userType string) []string {
	switch userType {
	case TokenTypeDefault:
		return []string{PermissionListDiet}
	case TokenTypeNutritionist:
		return []string{
			PermissionListDiet,
			PermissionCreateDiet,
			PermissionUpdateDiet,
			PermissionUploadFile,
		}
	default:
		return []string{}
	}
}
