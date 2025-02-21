package main

type Sender struct {
	Initiator string
	Email     string
}

type Recipient struct {
	Email string
}

type Mail struct {
	From    *Sender
	To      []*Recipient
	Subject string
	Body    string
}