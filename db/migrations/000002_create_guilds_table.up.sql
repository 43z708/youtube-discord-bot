CREATE TABLE IF NOT EXISTS guilds (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bot_id VARCHAR(255) NULL,
    category_id VARCHAR(255) NULL,
    admin_channel_id VARCHAR(255) NULL,
    youtube_api_key VARCHAR(255) UNIQUE NULL,
    CONSTRAINT guild_bot_id_fk FOREIGN KEY (bot_id) REFERENCES bots(id) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
