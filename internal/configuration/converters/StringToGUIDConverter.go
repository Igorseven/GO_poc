package converters

import (
	configIO "fmt"
	"syscall"
)

func BytesToString(bytes []byte) string {
	formattedString := configIO.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16])

	return formattedString
}

func GuidToStringForSql(id syscall.GUID) (*string, error) {
	guidString := configIO.Sprintf("%08x-%04x-%04x-%04x-%012x",
		id.Data1, id.Data2, id.Data3,
		id.Data4[:2], id.Data4[2:])

	return &guidString, nil
}
