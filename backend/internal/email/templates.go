package email

import "fmt"

// PasswordResetTemplate generates the HTML email body for password reset
func PasswordResetTemplate(code string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset Code</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px 40px; text-align: center;">
                            <h1 style="margin: 0; color: #333333; font-size: 24px;">Password Reset Request</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 20px 40px;">
                            <p style="margin: 0 0 20px 0; color: #666666; font-size: 16px; line-height: 24px;">
                                You requested to reset your password. Use the code below to proceed:
                            </p>
                            
                            <!-- Code Box -->
                            <div style="background-color: #f8f9fa; border: 2px solid #e9ecef; border-radius: 8px; padding: 30px; text-align: center; margin: 30px 0;">
                                <p style="margin: 0 0 10px 0; color: #666666; font-size: 14px; text-transform: uppercase; letter-spacing: 1px;">Your Reset Code</p>
                                <p style="margin: 0; color: #2563eb; font-size: 32px; font-weight: bold; letter-spacing: 8px; font-family: 'Courier New', monospace;">%s</p>
                            </div>
                            
                            <p style="margin: 20px 0 0 0; color: #666666; font-size: 14px; line-height: 20px;">
                                This code will expire in <strong>10 minutes</strong>.
                            </p>
                            
                            <p style="margin: 20px 0 0 0; color: #999999; font-size: 13px; line-height: 18px;">
                                If you didn't request this password reset, you can safely ignore this email. Your password will remain unchanged.
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f8f9fa; border-radius: 0 0 8px 8px;">
                            <p style="margin: 0; color: #999999; font-size: 12px; text-align: center;">
                                © 2025 Aroma Sense. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, code)
}

// WelcomeEmailTemplate generates the HTML email body for welcome emails
func WelcomeEmailTemplate(name string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Aroma Sense</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px;">
                    <tr>
                        <td style="padding: 40px; text-align: center;">
                            <h1 style="color: #2563eb;">Welcome to Aroma Sense, %s!</h1>
                            <p style="color: #666666; font-size: 16px; line-height: 24px;">
                                We're excited to have you on board. Discover our exclusive collection of fragrances.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, name)
}

// OrderConfirmationTemplate generates the HTML email body for order confirmation
func OrderConfirmationTemplate(orderID string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Confirmation</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px;">
                    <tr>
                        <td style="padding: 40px; text-align: center;">
                            <h1 style="color: #2563eb;">Order Confirmed!</h1>
                            <p style="color: #666666; font-size: 16px;">
                                Thank you for your order. Your order ID is: <strong>%s</strong>
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, orderID)
}
