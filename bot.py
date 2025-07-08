import os
import json
import random
from dotenv import load_dotenv
import discord
import logging

# .envファイルからトークンを読み込み
load_dotenv()
TOKEN = os.getenv("DISCORD_TOKEN")

# 必要なIntentの設定
intents = discord.Intents.default()
intents.message_content = True

# クライアント作成
bot = discord.Client(intents=intents)

# 語録とキーワードの読み込み
with open("quotes.json", encoding="utf-8") as f:
    quotes = json.load(f)
with open("keywords.json", encoding="utf-8") as f:
    keywords = json.load(f)

@bot.event
async def on_ready():
    print(f"Logged in as {bot.user} (ID: {bot.user.id})")

@bot.event
async def on_message(msg):
    # Botの発言は無視
    if msg.author.bot:
        return

    matched_ids = []

    # メッセージに含まれるキーワードに対応する語録IDを集める
    for kw, ids in keywords.items():
        if kw in msg.content:
            matched_ids.extend(ids)

    # 一致があれば10%の確率で反応して返信
    if matched_ids and random.random() < 1:  # 10% の確率で反応
        selected_id = random.choice(matched_ids)
        quote = quotes.get(selected_id, "（該当する語録が見つかりませんでした）")
        embed = discord.Embed(description=quote, color=0x1abc9c)
        reply_msg = await msg.reply(embed=embed)
        await reply_msg.add_reaction("❌")

@bot.event
async def on_reaction_add(reaction, user):
    if user.bot:
        return

    if str(reaction.emoji) == "❌":
        try:
            await reaction.message.delete()
        except Exception as e:
            print(f"Failed to remove reaction: {e}")

# Bot起動
bot.run(TOKEN)
