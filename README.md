# CQbot

A simple Discord bot that replies with preset quotes when specific keywords are detected in messages. This project is implemented in JavaScript (Node.js) using the discord.js library.

## Features

-   Detects keywords in messages and replies with a corresponding quote from `quotes.json`.
-   The probability of a response can be configured.
-   Slash commands for random quotes (`/random_quote`), text transformation (`/revolutionized`), and health check (`/ping`).
-   Bot responses can be deleted by reacting with '❌'.
-   Keywords and quotes are easily managed via `keywords.json` and `quotes.json`.

## Prerequisites

-   [Node.js](https://nodejs.org/) (v16.9.0 or higher)
-   [npm](https://www.npmjs.com/)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/mellllonsoda/CQbot.git
    cd CQbot
    ```

2.  **Install dependencies:**
    ```bash
    npm install
    ```

3.  **Create a configuration file:**
    Create a file named `.env` in the root of the project and add the following, replacing the placeholder values:
    ```
    DISCORD_TOKEN=YourDiscordBotToken
    CLIENT_ID=YourBotApplicationClientID
    DEV_GUILD_ID=YourDevelopmentServerID
    ```
    -   `DISCORD_TOKEN`: Your Discord bot's token.
    -   `CLIENT_ID`: The application's client ID.
    -   `DEV_GUILD_ID`: The ID of the server where you want to instantly register the slash commands for testing.

4.  **Register Slash Commands:**
    Run the deploy script to register the slash commands with your server. You only need to do this once, or whenever you change the commands.
    ```bash
    npm run deploy
    ```

## Usage

**Start the bot:**
```bash
npm run start
```
The bot should now be online and responding to commands and keywords in your server.

## License

This project is licensed under the ISC License.
