package mailer

import (
	"errors"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	r "github.com/rubikorg/rubik"
)

// BlockName is the name of this block
const BlockName = "mailer"

// BlockMailer is a rubik block with simplified implementation
// to send mails from your rubik project
type BlockMailer struct {
	config mailerConfig
}

// Details are values required to send a mail
type Details struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	isHTML      bool
	UseAuth     bool
	Port        int
	Attachments []string
}

type mailerConfig struct {
	ShouldAuth bool `toml:"auth"`
	Host       string
	Username   string
	Password   string
}

// OnAttach is the implementation of rubik.Block for this package
func (m BlockMailer) OnAttach(app *r.App) error {
	var mconf mailerConfig
	err := app.Decode(BlockName, &mconf)
	if err != nil {
		return err
	}

	if mconf.ShouldAuth {
		if mconf.Host == "" || m.config.Username == "" || mconf.Password == "" {
			msg := "mailer requires username, host and password, one or more of which is missing"
			return errors.New(msg)
		}
	}
	return nil
}

// Send email with the provided details
func (m BlockMailer) Send(d Details) error {
	if d.Subject == "" {
		return errors.New("subject is empty")
	} else if d.Body == "" {
		return errors.New("body is empty")
	}

	e := email.NewEmail()
	e.From = d.From
	e.To = d.To
	e.Cc = d.Cc
	e.Bcc = d.Bcc
	e.Subject = d.Subject
	if d.isHTML {
		e.HTML = []byte(d.Body)
	} else {
		e.Text = []byte(d.Body)
	}

	port := 587
	if d.Port != 0 {
		port = d.Port
	}
	host := fmt.Sprintf("%s:%d", m.config.Host, port)

	var auth smtp.Auth
	if m.config.ShouldAuth || d.UseAuth {
		auth = smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	}

	if len(d.Attachments) > 0 {
		for _, a := range d.Attachments {
			_, err := e.AttachFile(a)
			if err != nil {
				return err
			}
		}
	}

	return e.Send(host, auth)
}

func init() {
	r.Attach(BlockName, BlockMailer{})
}
