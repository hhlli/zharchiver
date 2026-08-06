import re

with open('zharchiver/services/telegram_sender.go', 'r') as f:
    content = f.read()

# Replace signatures
content = content.replace('func sendMediaGroup(token, chatID string', 'func sendMediaGroup(apiEndpoint, token, chatID string')
content = content.replace('func sendText(token, chatID string', 'func sendText(apiEndpoint, token, chatID string')
content = content.replace('func sendSingleMedia(token, chatID string', 'func sendSingleMedia(apiEndpoint, token, chatID string')
content = content.replace('func ShareAnswerToTelegram(token, chatID string', 'func ShareAnswerToTelegram(db *sql.DB, token, chatID string')

# Replace api.telegram.org in urls
content = content.replace('url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMediaGroup", token)', 'url := fmt.Sprintf("%s/bot%s/sendMediaGroup", apiEndpoint, token)')
content = content.replace('url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)', 'url := fmt.Sprintf("%s/bot%s/sendMessage", apiEndpoint, token)')
content = content.replace('url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)', 'url := fmt.Sprintf("%s/bot%s/%s", apiEndpoint, token, method)')

# Add apiEndpoint generation inside ShareAnswerToTelegram
old_share = 'func ShareAnswerToTelegram(db *sql.DB, token, chatID string, data *models.AnswerData) error {'
new_share = old_share + '\n\tapiEndpoint := models.GetTelegramAPIEndpoint(db)'
content = content.replace(old_share, new_share)

# Update callers inside ShareAnswerToTelegram
content = content.replace('sendMediaGroup(token, chatID', 'sendMediaGroup(apiEndpoint, token, chatID')
content = content.replace('sendSingleMedia(token, chatID', 'sendSingleMedia(apiEndpoint, token, chatID')
content = content.replace('sendText(token, chatID', 'sendText(apiEndpoint, token, chatID')

# Inside AutoPushToTelegram, pass db
content = content.replace('err := ShareAnswerToTelegram(token, chatID, data)', 'err := ShareAnswerToTelegram(db, token, chatID, data)')

with open('zharchiver/services/telegram_sender.go', 'w') as f:
    f.write(content)

with open('zharchiver/handlers/telegram_share.go', 'r') as f:
    handler_content = f.read()

handler_content = handler_content.replace('err := services.ShareAnswerToTelegram(token, chatID, data)', 'err := services.ShareAnswerToTelegram(env.db, token, chatID, data)')

with open('zharchiver/handlers/telegram_share.go', 'w') as f:
    f.write(handler_content)

print("Done sender")
