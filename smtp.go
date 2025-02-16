package main

type SmtpServer struct {
	Server 	string
	Port 	int
	Auth 	bool
	User 	string
	Pass 	string
	Tls		string
}

type Mail struct {
	From 	string
	To 		[]string
	Subject string
	Body 	string
}