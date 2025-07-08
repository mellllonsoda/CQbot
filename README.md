# CQbot

CQbot is a simple Discord bot that replies with preset quotes when specific keywords are detected in messages.

## Features
- Detects keywords in messages and replies with a corresponding quote with a 10% probability
- Quotes and keywords are managed via JSON files
- Replies are sent as embedded messages

## Usage
- When the bot is running in your server, if a message contains a keyword registered in `keywords.json`, the bot will reply with a corresponding quote with a 10% chance.

## Notes
- The bot ignores its own messages.
- You can freely add or edit quotes and keywords in the JSON files.

## License
This project is licensed under the MIT License.
