package terraform

import "fmt"

// Error wraps an error returned by the Terraform executable. There are two
// errors contained. Go is the error returned by the go exec framework, while
// Terraform is an error made up of the stderr output of the Terraform command.
type Error struct {
	Command   string
	Go        error
	Terraform error
}

// Error makes our Error struct match the standard Go error interface.
func (e Error) Error() string {
	return fmt.Sprintf("terraform command (%s) failed (%s) (%s)", e.Command, e.Go.Error(), e.Terraform.Error())
}
