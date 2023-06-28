DROP TABLE IF EXISTS ban;

CREATE TABLE IF NOT EXISTS ban (
                                    id SERIAL PRIMARY KEY,
                                    ban_words TEXT
);


INSERT INTO ban (stop_list) VALUES ('qwerty');

INSERT INTO ban (stop_list) VALUES ('йцукен');

INSERT INTO ban (stop_list) VALUES ('zxvbnm');