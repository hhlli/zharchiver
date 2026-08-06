import re

with open('zharchiver/services/telegram.go', 'r') as f:
    content = f.read()

# Add db *sql.DB to functions
content = content.replace('func sendMessageToTelegram(token', 'func sendMessageToTelegram(db *sql.DB, token')
content = content.replace('func sendMessageToTelegramWithMarkup(token', 'func sendMessageToTelegramWithMarkup(db *sql.DB, token')
content = content.replace('func editMessageText(token', 'func editMessageText(db *sql.DB, token')
content = content.replace('func answerCallbackQuery(token', 'func answerCallbackQuery(db *sql.DB, token')
content = content.replace('func downloadTelegramFile(token', 'func downloadTelegramFile(db *sql.DB, token')

# Update function calls inside StartTelegramBotListener
content = content.replace('sendMessageToTelegram(token', 'sendMessageToTelegram(db, token')
content = content.replace('sendMessageToTelegramWithMarkup(token', 'sendMessageToTelegramWithMarkup(db, token')
content = content.replace('editMessageText(token', 'editMessageText(db, token')
content = content.replace('answerCallbackQuery(token', 'answerCallbackQuery(db, token')
content = content.replace('downloadTelegramFile(token', 'downloadTelegramFile(db, token')

# Replace api.telegram.org in fmt.Sprintf with dynamic endpoint
content = content.replace('"https://api.telegram.org/bot%s/sendMessage", token', '"%s/bot%s/sendMessage", models.GetTelegramAPIEndpoint(db), token')
content = content.replace('"https://api.telegram.org/bot%s/editMessageText", token', '"%s/bot%s/editMessageText", models.GetTelegramAPIEndpoint(db), token')
content = content.replace('"https://api.telegram.org/bot%s/answerCallbackQuery", token', '"%s/bot%s/answerCallbackQuery", models.GetTelegramAPIEndpoint(db), token')
content = content.replace('"https://api.telegram.org/bot%s/getFile?file_id=%s", token, fileID', '"%s/bot%s/getFile?file_id=%s", models.GetTelegramAPIEndpoint(db), token, fileID')
content = content.replace('"https://api.telegram.org/file/bot%s/%s", token, fileResp.Result.FilePath', '"%s/file/bot%s/%s", models.GetTelegramAPIEndpoint(db), token, fileResp.Result.FilePath')
content = content.replace('"https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=30", token, lastUpdateID+1', '"%s/bot%s/getUpdates?offset=%d&timeout=30", models.GetTelegramAPIEndpoint(db), token, lastUpdateID+1')

with open('zharchiver/services/telegram.go', 'w') as f:
    f.write(content)

print("Done")
