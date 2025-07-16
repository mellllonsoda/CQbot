import os
import json
import random
from dotenv import load_dotenv
import discord
from discord import app_commands
import logging

# .envファイルからトークンを読み込み
load_dotenv()
TOKEN = os.getenv("DISCORD_TOKEN")

# 必要なIntentの設定
intents = discord.Intents.default()
intents.message_content = True

# クライアント作成
bot = discord.Client(intents=intents)
tree = app_commands.CommandTree(bot)

# 語録とキーワードの読み込み
with open("quotes.json", encoding="utf-8") as f:
    quotes = json.load(f)
with open("keywords.json", encoding="utf-8") as f:
    keywords = json.load(f)    

@tree.command(name="random_quote", description="ランダムに☭革命的☭な名言を出す")
async def test_command(interaction: discord.Interaction):
    random_id = random.choice(list(quotes.keys()))
    quote = quotes[random_id]
    embed = discord.Embed(description=quote, color=0x1abc9c)

    # まず応答を遅延(defer)させる
    await interaction.response.defer()

    # フォローアップメッセージとして送信し、メッセージオブジェクトを取得
    message = await interaction.followup.send(embed=embed)

@bot.event
async def on_ready():
    print(f"Logged in as {bot.user} (ID: {bot.user.id})")
    await tree.sync()
    print(f"スラッシュコマンドを同期しました")

@bot.event
async def on_message(msg):
    # Bot自身の発言には❌リアクションをつける
    if msg.author == bot.user:
        await msg.add_reaction("❌")
        return  # 自分の発言はそれだけで終了

    matched_ids = []

    # ユーザーのメッセージからキーワードを探す
    for kw, ids in keywords.items():
        if kw in msg.content:
            matched_ids.extend(ids)

    # 一致があれば10%の確率で反応
    if matched_ids and random.random() < 0.1:
        selected_id = random.choice(matched_ids)
        quote = quotes.get(selected_id, "（該当する語録が見つかりませんでした）")
        embed = discord.Embed(description=quote, color=0x1abc9c)
        await msg.channel.send(embed=embed)
        
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