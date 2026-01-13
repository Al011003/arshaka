package usecase

func isEditable(flow string) bool {
	return flow == "REQUESTED" || flow == "REVISION"
}
