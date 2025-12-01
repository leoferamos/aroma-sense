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

// AccountDeactivatedTemplate generates the HTML body for account deactivation notification
func AccountDeactivatedTemplate(reason string, contestationDeadline string) string {
	return fmt.Sprintf(`
<h2>Conta Desativada</h2>
<p>Olá,</p>
<p>Informamos que sua conta no Aroma Sense foi desativada pelos seguintes motivos:</p>
<p><strong>%s</strong></p>
<p>Você tem até <strong>%s</strong> para apresentar contestação através do nosso suporte.</p>
<p>Para contestar, acesse sua conta ou entre em contato conosco.</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`, reason, contestationDeadline)
}

// ContestationReceivedTemplate generates the HTML body confirming receipt of contestation
func ContestationReceivedTemplate() string {
	return `
<h2>Contestação Recebida</h2>
<p>Olá,</p>
<p>Recebemos sua contestação sobre a desativação da conta.</p>
<p>Nossa equipe irá analisar o caso em até 5 dias úteis e entraremos em contato.</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`
}

// ContestationResultTemplate generates the HTML body for contestation review results
func ContestationResultTemplate(approved bool, reason string) string {
	status := "rejeitada"
	if approved {
		status = "aprovada"
	}
	return fmt.Sprintf(`
<h2>Resultado da Contestação</h2>
<p>Olá,</p>
<p>Sua contestação foi <strong>%s</strong>.</p>
<p><strong>Motivo:</strong> %s</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`, status, reason)
}

// DeletionRequestedTemplate informs the user their deletion request was received and how to cancel
func DeletionRequestedTemplate(name, requestedAt, cancelLink string) string {
	return fmt.Sprintf(`
<h2>Pedido de Exclusão Recebido</h2>
<p>Olá %s,</p>
<p>Recebemos seu pedido de exclusão em %s.</p>
<p>Você tem 7 dias para cancelar a solicitação. Para cancelar, acesse: <a href="%s">Cancelar exclusão</a></p>
<p>Se não solicitou esta ação, entre em contato com o suporte.</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`, name, requestedAt, cancelLink)
}

// DeletionAutoConfirmedTemplate notifies user their deletion was auto-confirmed
func DeletionAutoConfirmedTemplate(name, confirmedAt string) string {
	return fmt.Sprintf(`
<h2>Exclusão Confirmada</h2>
<p>Olá %s,</p>
<p>Seu pedido de exclusão foi confirmado em %s. Seus dados serão retidos por 2 anos antes de serem anonimizados.</p>
<p>Se você acredita que isto é um erro, entre em contato com o suporte.</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`, name, confirmedAt)
}

// DataAnonymizedTemplate notifies user their personal data was anonymized
func DataAnonymizedTemplate(anonymousDate string) string {
	return fmt.Sprintf(`
<h2>Dados Anonimizados</h2>
<p>Olá,</p>
<p>Conforme sua solicitação e nossa política, seus dados pessoais foram anonimizados em %s.</p>
<p>Se precisar de mais informações, contate suporte.</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`, anonymousDate)
}

// DeletionCancelledTemplate notifies the user their deletion request was cancelled
func DeletionCancelledTemplate(name, cancelledAt string) string {
    return fmt.Sprintf(`
<h2>Solicitação de Exclusão Cancelada</h2>
<p>Olá %s,</p>
<p>Sua solicitação de exclusão foi cancelada em %s. Sua conta permanece ativa.</p>
<p>Se precisar de ajuda, entre em contato com o suporte.</p>
<p>Atenciosamente,<br>Equipe Aroma Sense</p>
`, name, cancelledAt)
}
