

CREATE TABLE authors(
id INT AUTO_INCREMENT NOT NULL,
name VARCHAR(255) NOT NULL,
PRIMARY KEY(id)
);

CREATE TABLE albumns(
id INT AUTO_INCREMENT NOT NULL,
name VARCHAR(255) NOT NULL,
author_id INT NOT NULL,
albumn_type VARCHAR(36),
 albumn_date VARCHAR(200) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id)

);


CREATE TABLE songs(
id INT AUTO_INCREMENT NOT NULL,
name VARCHAR(255) NOT NULL,
author_id INT NOT NULL,
albumns_id INT NOT NULL,
song_index INT NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);


CREATE TABLE genres(
name VARCHAR(100) NOT NULL,
PRIMARY KEY(name)
);


CREATE TABLE tags(
name VARCHAR(100) NOT NULL,
PRIMARY KEY(name)
);

CREATE TABLE similar_tags(
id  VARCHAR(36) NOT NULL,
first_tag VARCHAR(36) NOT NULL,
second_tag VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(first_tag) REFERENCES tags(name),
FOREIGN KEY(second_tag) REFERENCES tags(name)
);

CREATE TABLE similar_genres(
id  VARCHAR(36) NOT NULL,
first_genre VARCHAR(36) NOT NULL,
second_genre VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(first_genre) REFERENCES genres(name),
FOREIGN KEY(second_genre) REFERENCES genres(name)
);



CREATE TABLE authors_to_songs(
id  VARCHAR(36) NOT NULL,
author_id VARCHAR(36) NOT NULL,
songs_id  VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id),
FOREIGN KEY(songs_id) REFERENCES songs(id)
);


CREATE TABLE tags_to_songs(
id VARCHAR(36) NOT NULL,
tag_name VARCHAR(36) NOT NULL,
song_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(tag_name) REFERENCES tags(name),
FOREIGN KEY(song_id) REFERENCES songs(id)
);




CREATE TABLE genres_to_songs(
id VARCHAR(36) NOT NULL,
genre_name VARCHAR(255) NOT NULL,
song_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(genre_name) REFERENCES genres(name),
FOREIGN KEY(song_id) REFERENCES songs(id)
);

CREATE TABLE songs_views(
id VARCHAR(36) NOT NULL,
song_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(song_id) REFERENCES songs(id)
);

CREATE TABLE albumns_views(
id VARCHAR(36) NOT NULL,
albumn_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);

CREATE TABLE authors_views(
id VARCHAR(36) NOT NULL,
author_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id)
);

