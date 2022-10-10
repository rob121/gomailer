# Go mailer

Ergonomic Mail Sending package using gomail.

## Basic Usage

```
 gomailer.Config(config) //set the overall config, from, host, etc this is a map, map
 gomailer.Template(tmpl) //set a overall *template.Template for each message, the data to execute is "
 mm := gomailer.New()
 mm.To("dave@example.com")
 mm.Subject("Test Mail")
 mm.Send("msgstring")
```

## Modes

You can set a message body in a few ways:

* Set a email template and then message within the template, for nice formatted emails
* Set just a body string for brief messages, untemplated
* Set a data source to pull a message from and use the message key to send a message, wrapping a template

### Set a email template to wrap your message

```
tm,_ := template.ParseFiles("email_template.html")
gomailer.Template(tm)

//create data to replace in the template/body
mp := make(mp[string]interface{})
mp["Extra"] = "the end"

mm.Data(mp) //"Extra" will also be available to the template above
mm.Send("This is also a html/template template that will get parsed {{.Extra}}")

```
### Send just a body, either html or plain text
```
mm.IsPlain = true
mm.Send("Hi, This is my whole message")
```

### Set a data source for use later, set this function when to fetch a message, for example a cms
```
gomailer.DataSource(func(key interface{}) (string){

   qry := `SELECT from Content WHERE k=key`
   ....

})


mp := make(mp[string]interface{})
mm.Data(mp)
mm.Content("article-lookup-key") using datasource method
mm.Send()

```
## Creating a template

The only thing that is mandatory to create a template is the value for the message within the tempalte is "MessageContent" refrenced via {{.MessageContent}}


## Configuration

Two ways to set configuration, set a map, or use the in built Conf type

```
	gomailer.Config(conf) //set the overall config, from, host, etc this is a map[string]interface{}
```

Set individual settings in the config

```
		gomailer.Conf.Set(k, v)
```

### Config Validation

We require the following "Keys" in the config

"Host", "Port", "User", "Password", "From", "Reply"

### Underlying Gomail Object

The system uses gomail at a low level, access other low level methods under m.Msg

#### Attach a file

m.Msg.Attach("path/to/file.html")

