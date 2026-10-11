// Package emailtemplate renders the emails Nokori sends. The wording lives in
// templates/ as plain files and is embedded into the binary at build time, so
// nothing is read from disk at runtime. Kept outside infrastructure so the
// domain layer can render emails without importing infrastructure.
package emailtemplate

import (
	"bytes"
	"embed"
	htmltemplate "html/template"
	texttemplate "text/template"
)

//go:embed templates/*
var files embed.FS

// Parsed once at startup; a broken template panics immediately instead of
// failing on the first signup.
var (
	htmlTemplates = htmltemplate.Must(htmltemplate.ParseFS(files, "templates/*.html"))
	textTemplates = texttemplate.Must(texttemplate.ParseFS(files, "templates/*.txt"))
)

type Email struct {
	Subject string
	HTML    string
	Text    string
}

type VerificationCodeData struct {
	Code         string
	ValidMinutes int
}

func VerificationCode(data VerificationCodeData) (Email, error) {
	return render("verification_code", data)
}

// render fills <name>.subject.txt, <name>.html and <name>.txt with data.
func render(name string, data any) (email Email, err error) {
	var buf bytes.Buffer

	if err = textTemplates.ExecuteTemplate(&buf, name+".subject.txt", data); err != nil {
		return
	}
	email.Subject = buf.String()

	buf.Reset()
	if err = htmlTemplates.ExecuteTemplate(&buf, name+".html", data); err != nil {
		return
	}
	email.HTML = buf.String()

	buf.Reset()
	if err = textTemplates.ExecuteTemplate(&buf, name+".txt", data); err != nil {
		return
	}
	email.Text = buf.String()

	return
}
