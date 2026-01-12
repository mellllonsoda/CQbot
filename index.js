
const fs = require('fs');
const { Client, GatewayIntentBits, EmbedBuilder, Collection, REST, Routes } = require('discord.js');
require('dotenv').config();

// Botの定数
const embedColor = 0x1abc9c;
const removeReaction = '❌';
const responseChance = 0.1;

// Botの構成
const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.GuildMessageReactions,
    GatewayIntentBits.MessageContent,
  ],
});

client.commands = new Collection();
const quotes = JSON.parse(fs.readFileSync('quotes.json', 'utf8'));
const keywords = JSON.parse(fs.readFileSync('keywords.json', 'utf8'));

// イベントハンドラ
client.once('ready', () => {
  console.log(`Logged in as ${client.user.tag}`);
});

client.on('messageCreate', (message) => {
  // 自分のメッセージは無視
  if (message.author.bot) {
    return;
  }

  // キーワードにマッチする語録IDを収集
  const matchedIDs = [];
  for (const kw in keywords) {
    if (message.content.includes(kw)) {
      matchedIDs.push(...keywords[kw]);
    }
  }

  // マッチした語録があり、一定の確率を満たした場合に返信
  if (matchedIDs.length > 0 && Math.random() < responseChance) {
    // 重複を除外したユニークなIDリストを作成
    const uniqueIDs = [...new Set(matchedIDs)];

    // ユニークIDからランダムに1つ選択
    const selectedID = uniqueIDs[Math.floor(Math.random() * uniqueIDs.length)];
    const quote = quotes[selectedID];

    if (!quote) {
      console.log(`Quote not found for ID: ${selectedID}`);
      return; // 見つからなければ何もしない
    }

    const embed = new EmbedBuilder()
      .setDescription(quote)
      .setColor(embedColor);

    message.channel.send({ embeds: [embed] }).catch(console.error);
  }
});

client.on('messageReactionAdd', async (reaction, user) => {
  // Bot自身のリアクションは無視
  if (user.bot) {
    return;
  }

  // 指定されたリアクションでない場合は無視
  if (reaction.emoji.name !== removeReaction) {
    return;
  }

  // リアクションがつけられたメッセージをフェッチ
  if (reaction.message.partial) {
    try {
      await reaction.message.fetch();
    } catch (error) {
      console.error('Failed to fetch message:', error);
      return;
    }
  }

  // メッセージの投稿者がBot自身である場合のみ削除
  if (reaction.message.author.id === client.user.id) {
    reaction.message.delete().catch(console.error);
  }
});

client.on('interactionCreate', async (interaction) => {
  if (!interaction.isChatInputCommand()) return;

  const command = client.commands.get(interaction.commandName);

  if (!command) {
    console.error(`No command matching ${interaction.commandName} was found.`);
    return;
  }

  try {
    await command.execute(interaction);
  } catch (error) {
    console.error(error);
    if (interaction.replied || interaction.deferred) {
      await interaction.followUp({ content: 'There was an error while executing this command!', ephemeral: true });
    } else {
      await interaction.reply({ content: 'There was an error while executing this command!', ephemeral: true });
    }
  }
});


// スラッシュコマンドのハンドラ
const commands = [
    {
        name: 'random_quote',
        description: 'ランダムに名言を出す',
        execute: async (interaction) => {
            if (Object.keys(quotes).length === 0) {
                await interaction.reply({ content: 'No quotes available.', ephemeral: true });
                return;
            }

            const keys = Object.keys(quotes);
            const randomID = keys[Math.floor(Math.random() * keys.length)];
            const quote = quotes[randomID];

            const embed = new EmbedBuilder()
                .setDescription(quote)
                .setColor(embedColor);

            await interaction.reply({ embeds: [embed] });
        }
    },
    {
        name: 'revolutionized',
        description: '入力を変換する',
        options: [{
            type: 3, // STRING
            name: 'message',
            description: '変換するメッセージ',
            required: true,
        }],
        execute: async (interaction) => {
            const message = interaction.options.getString('message');
            const transformed = message.split('').join('☆');
            await interaction.reply({ content: `🔴☭${transformed}☭🔴`, ephemeral: true });
        }
    },
    {
        name: 'ping',
        description: 'Bot の応答速度を測定',
        execute: async (interaction) => {
            await interaction.deferReply();
            const reply = await interaction.fetchReply();
            const latency = reply.createdTimestamp - interaction.createdTimestamp;
            await interaction.editReply(`Pong! Latency: ${latency}ms`);
        }
    }
];

commands.forEach(command => {
    client.commands.set(command.name, command);
});


// ログイン
client.login(process.env.DISCORD_TOKEN);
