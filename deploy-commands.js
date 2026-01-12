
const { REST, Routes } = require('discord.js');
require('dotenv').config();

const commands = [
  {
    name: 'random_quote',
    description: 'ランダムに名言を出す',
  },
  {
    name: 'revolutionized',
    description: '入力を変換する',
    options: [
      {
        type: 3, // STRING
        name: 'message',
        description: '変換するメッセージ',
        required: true,
      },
    ],
  },
  {
    name: 'ping',
    description: 'Bot の応答速度を測定',
  },
];

const rest = new REST({ version: '10' }).setToken(process.env.DISCORD_TOKEN);

(async () => {
  try {
    console.log('Started refreshing application (/) commands.');

    if (process.env.DEV_GUILD_ID) {
      await rest.put(
        Routes.applicationGuildCommands(process.env.CLIENT_ID, process.env.DEV_GUILD_ID),
        { body: commands },
      );
      console.log('Successfully reloaded application (/) commands for guild.');
    } else {
      await rest.put(
        Routes.applicationCommands(process.env.CLIENT_ID),
        { body: commands },
      );
      console.log('Successfully reloaded application (/) commands globally.');
    }
  } catch (error) {
    console.error(error);
  }
})();
