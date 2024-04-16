package entity

type Args struct {
	SchemaName         string
	TableName          string
	TableNameHump      string
	TableNameHumpUpper string

	PrimaryKey          string
	PrimaryKeyHump      string
	PrimaryKeyHumpUpper string

	ForeignKey     string
	ForeignKeyHump string

	FieldsKey          []string
	FieldsKeyHump      []string
	FieldsKeyHumpUpper []string

	FieldsQuestion []string
	FieldsModify   []string
}
