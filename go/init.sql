create schema game;
-- 用户投票记录
CREATE TABLE game.user_votes
(
    id          SERIAL PRIMARY KEY,
    user_uuid   VARCHAR(36) NOT NULL,                                     -- 浏览器用户UUID
    vote_option INT         NOT NULL CHECK (vote_option BETWEEN 1 AND 4), -- 选项1-4
    vote_time   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                      -- 投票时间
    CONSTRAINT unique_user_vote UNIQUE (user_uuid)                        -- 约束
);

-- 投票统计
CREATE TABLE game.vote_statistics
(
    option_id          INT PRIMARY KEY CHECK (option_id BETWEEN 1 AND 4), -- 选项1-4
    option_description VARCHAR(255),                                      -- 选项描述
    vote_count         INT DEFAULT 0 NOT NULL                             -- 该选项获得的票数
);

-- 初始化
INSERT INTO game.vote_statistics (option_id, option_description, vote_count)
VALUES (1, '文字游戏', 11),
       (2, '动作游戏', 22),
       (3, '冒险游戏', 33),
       (4, '射击游戏', 15);