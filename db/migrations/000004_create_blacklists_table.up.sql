CREATE TABLE IF NOT EXISTS blacklists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    distributor VARCHAR(255),
    bot_id VARCHAR(255),
    guild_id VARCHAR(255),
    CONSTRAINT blacklist_guild_id_fk FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE,
    CONSTRAINT blacklist_bot_id_fk FOREIGN KEY (bot_id) REFERENCES bots(id) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
