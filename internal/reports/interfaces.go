package reports

type Reporter interface {
	Print(reports []FileReport)
}
