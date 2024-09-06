package main

import (
	connector2 "browser-connector-client/connector"
	"log"
)

func main() {
	connector := connector2.NewBrowserConnector("http://localhost:8081/")
	session, err := connector.CreateSession()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	err = session.Goto("https://google.com")
	if err != nil {
		log.Printf(err.Error())
		return
	}

	err = session.ImplicitWait(10)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	content, err := session.GetPageContent()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Printf(content)

}
