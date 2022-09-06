DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          news_id INT,
                          parent_comment_id INT DEFAULT NULL,
                          FOREIGN KEY (parent_comment_id) REFERENCES comments (id) ON DELETE CASCADE ,
                          content TEXT NOT NULL DEFAULT 'empty',
                          pubtime BIGINT NOT NULL DEFAULT extract (epoch from now())


);

INSERT INTO comments(news_id,content)  VALUES (100,'Предзаведенный тестовый комментарий')