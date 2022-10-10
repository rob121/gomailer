package gomailer

import (
	"bytes"
	"errors"
	"gopkg.in/gomail.v2"
	"html/template"
	"regexp"
	"sort"
	"strings"
)

var temp *template.Template
var Conf config
var funcmap template.FuncMap
var datasource func(interface{}) string //datasource is a function to query a data source by a string, and get content returned

func init() {

	Conf = make(config)

}

type Mailer struct {
	Msg         *gomail.Message
	Body        string
	TmplData    map[string]interface{}
	Errors      []error
	IsPlain     bool
	ContentBody string
	Config      config
}

func New(fconf ...config) *Mailer {

	m := gomail.NewMessage()

	am := new(Mailer)

	am.Msg = m

	if len(fconf) > 0 {
		//pass in a custom config for this message
		am.Config = fconf[0]
	} else {
		//use the default set in gomailer.Config
		am.Config = Conf
	}

	am.IsPlain = false

	am.Msg.SetHeader("From", am.Config.GetString("From"))

	if len(am.Config.GetString("Reply")) > 0 {
		am.Msg.SetHeader("Reply-To", am.Config.GetString("Reply"))
	}

	return am

}

func (m *Mailer) Plain() *Mailer {

	m.IsPlain = true
	return m

}

func (m *Mailer) To(addresses ...string) *Mailer {

	m.Msg.SetHeader("To", addresses...)
	return m

}

func (m *Mailer) Subject(subj string) *Mailer {

	m.Msg.SetHeader("Subject", subj)
	return m

}

func (m *Mailer) CC(addresses ...string) *Mailer {

	m.Msg.SetHeader("Cc", addresses...)
	return m

}

func (m *Mailer) BCC(addresses ...string) *Mailer {

	m.Msg.SetHeader("Bcc", addresses...)
	return m

}

//Sets a custom template for the mailer"

func (m *Mailer) Data(data map[string]interface{}) *Mailer {

	m.TmplData = data
	return m
}

func (m *Mailer) template() *Mailer {

	var b bytes.Buffer

	if m.TmplData == nil {

		m.TmplData = make(map[string]interface{})
	}

	m.TmplData["MessageContent"] = m.ContentBody

	e := temp.Execute(&b, m.TmplData)

	if e != nil {
		m.Errors = append(m.Errors, e)
	}

	m.Body = b.String()
	return m
}

func (m *Mailer) Attach(dir string) *Mailer {

	m.Msg.Attach(dir)
	return m

}

func (m *Mailer) Content(key string) *Mailer {

	if datasource == nil {
		m.Errors = append(m.Errors, errors.New("Datasource not defined, use gomailer.DataSource to set one"))
	}

	content := datasource(key)

	tm, tmerr := template.New("bodymessage").Funcs(funcmap).Parse(content)

	if tmerr != nil {
		m.Errors = append(m.Errors, tmerr)
	}

	var b bytes.Buffer

	err := tm.Execute(&b, m.TmplData)

	if err != nil {
		m.Errors = append(m.Errors, err)
		return m
	}

	m.ContentBody = b.String()

	return m

}

func (m *Mailer) Send(body_raw ...string) error {

	//embed logo

	if !m.Config.validate() {
		return errors.New("Invalid Config")
	}

	var body string
	if len(body_raw) > 0 {
		body = body_raw[0]
	}

	if m.IsPlain == true {

		m.Msg.SetBody("text/plain", body)

	} else {

		//if we have a template, use that
		if temp != nil && len(m.ContentBody) > 0 { //we got set via Content()
			m.template()
		} else if temp != nil { //use the template we have and wrap it around our message string
			m.ContentBody = body
			m.template()
		} else { //use the body passed as the entire message
			m.Body = body
		}

		if len(m.Errors) > 0 {
			//just return the first
			return m.Errors[0]
		}

		m.Msg.SetBody("text/html", m.Body)

	}

	d := gomail.NewDialer(m.Config.GetString("Host"), m.Config.GetInt("Port"), m.Config.GetString("User"), m.Config.GetString("Password"))

	// Send the email
	if err := d.DialAndSend(m.Msg); err != nil {

		return err
	}

	if len(m.Errors) > 0 {

		return m.Errors[0]
	}

	return nil

}

func Config(c map[string]interface{}) {
	Conf = make(config)
	Conf = c
}

func DataSource(fn func(interface{}) string) {
	datasource = fn
}

func Template(t *template.Template) {
	temp = t
}

func FuncMap(f template.FuncMap) {
	funcmap = f
}

func removeHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}
