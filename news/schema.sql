DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
                      id SERIAL PRIMARY KEY,
                      title TEXT,
                      content TEXT NOT NULL UNIQUE,
                      pubTime BIGINT DEFAULT 0,
                      link TEXT
);
