package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(to, seller_name, product_name, product_id string) error {
	from := "samuelbonux10@gmail.com"
	password := "sbaa fgpv zdcs wucw"

	subject := "Alert: Product Stock Near Minimum Threshold"
	body := fmt.Sprintf(
		"Dear %s,\n\n"+
			"We hope this message finds you well. This is a notification regarding your product:\n\n"+
			"  - Product Name: %s\n"+
			"  - Product ID: %s\n\n"+
			"The stock for this product is nearing its minimum threshold. To avoid potential stock-outs, we recommend restocking it as soon as possible.\n\n"+
			"Thank you for your attention to this matter.\n\n"+
			"Best regards,\n"+
			"The Inventory Management Team",
		seller_name, product_name, product_id,
	)

	smtpServer := "smtp.gmail.com"
	auth := smtp.PlainAuth("", from, password, smtpServer)

	msg := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	conn, err := tls.Dial("tcp", smtpServer+":465", tlsconfig)
	if err != nil {
		log.Printf("tls connection error: %v", err)
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		log.Printf("smtp client error: %v", err)
		return err
	}

	if err = client.Auth(auth); err != nil {
		log.Printf("auth error: %v", err)
		return err
	}

	if err = client.Mail(from); err != nil {
		log.Printf("mail error: %v", err)
		return err
	}

	if err = client.Rcpt(to); err != nil {
		log.Printf("recipient error: %v", err)
		return err
	}

	wc, err := client.Data()
	if err != nil {
		log.Printf("data error: %v", err)
		return err
	}

	_, err = wc.Write(msg)
	if err != nil {
		log.Printf("write error: %v", err)
		return err
	}

	err = wc.Close()
	if err != nil {
		log.Printf("close error: %v", err)
		return err
	}

	client.Quit()
	log.Println("email sent successfully")
	return nil
}
