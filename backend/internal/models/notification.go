package models

type FormField struct {
	Id       string `json:"name"`
	Title    string `json:"display_name"`
	Name     string `json:"placeholder"`
	Type     string `json:"type"`
	Default  string `json:"default"`
	Optional bool   `json:"optional"`

	// display_name	String	Display name of the field shown to the user in the dialog. Maximum 24 characters.
	// name	String	Name of the field element used by the integration. Maximum 300 characters. You should use unique name fields in the same dialog.
	// type	String	Set this value to bool for a checkbox element.
	// default	String	(Optional) Set a default value for this form element. true or false.
}
