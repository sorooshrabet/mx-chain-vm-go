package vmhooksgenerate

import (
	"fmt"
)

func autoGeneratedHeader(out *eiGenWriter) {
	out.WriteString(`// Code generated by vmhooks generator. DO NOT EDIT.

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
`)
}

func autoGeneratedGoHeader(out *eiGenWriter, packageName string) {
	out.WriteString(fmt.Sprintf(`package %s

`, packageName))
	autoGeneratedHeader(out)
}
