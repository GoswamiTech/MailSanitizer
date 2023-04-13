# MailSanitizer
service which senitize email by deleting spams

Refer [Using OAuth 2.0 to Access Google APIs](https://developers.google.com/identity/protocols/oauth2) and [Oauth Scopes](https://developers.google.com/gmail/api/reference/rest/v1/users.messages/delete) to setup google API authentication

```bash
./bin/MailSanitizer --h
 usage:
	./bin/MailSanitizer [-credential <credential file path>] [-label <label>] [-spammers <spammers file path>]

options:
	-credential     Optional 	oauth credential filen downloaded from google console(default value is credentials.json)
	-label		    Optional 	labelIds of labels/folder, separated by ","
	-token 			Optiobal	path of token file (default value is token.json)
	-spammers		Optional    path of spamers file. If not provided, all mails will be deleted into label
	-user			Optional 	emailId, default value is me
```