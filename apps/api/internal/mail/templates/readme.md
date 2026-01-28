How to use this template :

```
// Reset Password Email
data := map[string]interface{}{
    "Name":     "John Doe",
    "Email":    "john@example.com",
    "ResetURL": "https://yourapp.com/reset?token=abc123",
    "Token":    "abc123xyz789",
}
html, _ := templateEngine.Render("reset_password.html", data)
mailer.SendTo("john@example.com", "Reset Your Password", html)

// Welcome Email
data := map[string]interface{}{
    "Name":      "Jane Doe",
    "Email":     "jane@example.com",
    "CreatedAt": "January 6, 2026",
    "LoginURL":  "https://yourapp.com/login",
}
html, _ := templateEngine.Render("welcome.html", data)
mailer.SendTo("jane@example.com", "Welcome to Blog API", html)

```
