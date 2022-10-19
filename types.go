package goc

type FileData struct {
	CalendarId  string
	CurrentTask DataTask
	TaskAlias   map[string]string
}

type DataTask struct {
	Name  string
	Start string
}

func (f *DataTask) Reset() {
	f.Name = ""
	f.Start = ""
}
