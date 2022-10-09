package goc

type FileData struct {
	CalendarId  string
	CurrentTask DataTask
}

type DataTask struct {
	Name  string
	Start string
}
