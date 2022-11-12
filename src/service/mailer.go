package service

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"sync"
	"time"

	"github.com/jedzeins/go-mailer/src/config"
	"github.com/vanng822/go-premailer/premailer"

	mail "github.com/xhit/go-simple-mail/v2"
)

type MailService interface {
	BuildHTMLMessage(msg Message) (string, error)
	BuildPlainTextMessage(msg Message) (string, error)
	GetEncryption(e string) mail.Encryption
	InlineCSS(s string) (string, error)
	SendMail(msg Message, errorChan chan error)
}

type MailServiceImpl struct {
	Domain      string
	Host        string
	Port        string
	Username    string
	Password    string
	DefaultFrom string
	Encryption  string
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

type Message struct {
	From        string `json:"from,omitempty"`
	FromName    string `json:"fromName,omitempty"`
	To          string `json:"to,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Attachments []string
	Data        MessageBody `json:"messageBody,omitempty"`
	Template    string
}

type MessageBody struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
	Url          string `json:"url,omitempty"`
}

func NewMailService(config *config.Config, wg *sync.WaitGroup) *MailServiceImpl {
	// create channels
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	return &MailServiceImpl{
		Domain:      config.SMTP_HOST,
		Host:        config.SMTP_HOST,
		Port:        config.SMTP_PORT,
		Encryption:  config.EMAIL_ENCRYPTION_TYPE,
		Username:    config.SMTP_USERNAME,
		Password:    config.SMTP_PASSWORD,
		DefaultFrom: config.EMAIL_FROM,
		Wait:        wg,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
	}
}

// a function to listen for messages on the MailerChan

func (m *MailServiceImpl) SendMail(msg Message, errorChan chan error) {
	defer m.Wait.Done()

	if msg.Template == "" {
		msg.Template = "mail"
	}

	if msg.From == "" {
		msg.From = m.DefaultFrom
	}

	if msg.FromName == "" {
		msg.FromName = m.DefaultFrom
	}

	// build html mail
	formattedMessage, err := m.BuildHTMLMessage(msg)
	if err != nil {
		fmt.Println("err 1: ", err.Error())
		errorChan <- err
	}

	// build plain text mail
	plainMessage, err := m.BuildPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}

	port, err := strconv.Atoi(m.Port)
	if err != nil {
		errorChan <- err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.GetEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}

}

func (m *MailServiceImpl) BuildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./src/templates/%s.html.gohtml", msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.InlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *MailServiceImpl) BuildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./src/templates/%s.plain.gohtml", msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		fmt.Println("err 4: ", err.Error())
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *MailServiceImpl) InlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *MailServiceImpl) GetEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
