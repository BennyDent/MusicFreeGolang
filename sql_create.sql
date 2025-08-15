

CREATE TABLE authors(
id VARCHAR(36) NOT NULL,
name VARCHAR(255) NOT NULL,
PRIMARY KEY(id)
);

CREATE TABLE albumns(
id VARCHAR(36) NOT NULL,
name VARCHAR(255) NOT NULL,
author_id INT NOT NULL,
albumn_type VARCHAR(36),
 albumn_date VARCHAR(200) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id)

);


CREATE TABLE songs(
id VARCHAR(36) NOT NULL,
name VARCHAR(255) NOT NULL,
author_id INT NOT NULL,
albumn_id INT NOT NULL,
song_index INT NOT NULL,
cover_filename VARCHAR(36) DEFAULT NULL,
song_filename VARCHAR(36) DEFAULT NULL,
albumn_index INT NOT NULL,
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


CREATE TABLE radio_listened (
  id varchar(36) NOT NULL,
  user_id varchar(50) NOT NULL,
  song_id varchar(36) DEFAULT NULL,
  last_listened datetime DEFAULT CURRENT_TIMESTAMP,
  expires datetime NOT NULL,
  PRIMARY KEY (id),
FOREIGN KEY (user_id) REFERENCES users (user_id),
  FOREIGN KEY (song_id) REFERENCES songs (id) );

CREATE TABLE authors_to_albumns(
id  VARCHAR(36) NOT NULL,
author_id VARCHAR(36) NOT NULL,
albumn_id  VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);

CREATE TABLE tags_to_songs(
id VARCHAR(36) NOT NULL,
tag_name VARCHAR(36) NOT NULL,
song_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(tag_name) REFERENCES tags(name),
FOREIGN KEY(song_id) REFERENCES songs(id)
);


CREATE TEMPORARY TABLE  guest_session(
    session_id VARHCAR(100) NOT NULL

    PRIMARY KEY(session_id);

)

CREATE TEMPORARY TABLE guest_session_view(
id VARCHAR(36) NOT NULL,
session_id VARCHAR(100) NOT NULL,
song_id VARCHAR(36) NOT NULL,
last_listened DATETIME DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY(id),
FOREIGN KEY(session_id) REFERENCES guest_session(session_id),
FOREIGN KEY (song_id) REFERENCES songs(id)
)

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
listened INT DEFAULT 0,
last_listened DATETIME DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY(id),
FOREIGN KEY(song_id) REFERENCES songs(id)
);

CREATE TABLE albumns_views(
id VARCHAR(36) NOT NULL,
albumn_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,
listened INT DEFAULT 0,
last_listened DATETIME DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY(id),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);

CREATE TABLE authors_views(
id VARCHAR(36) NOT NULL,
author_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,
listened INT DEFAULT 0,
last_listened DATETIME DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id)
);

CREATE TABLE songs_views(
id VARCHAR(36) NOT NULL,
song_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,
listened INT DEFAULT 0,
last_listened DATETIME DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY(id),
FOREIGN KEY(song_id) REFERENCES songs(id)
);

CREATE TABLE tags_to_albumns(
id VARCHAR(36) NOT NULL,
tag_name VARCHAR(36) NOT NULL,
albumn_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(tag_name) REFERENCES tags(name),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);



CREATE TABLE genress_to_albumns(
id VARCHAR(36) NOT NULL,
genre_name VARCHAR(36) NOT NULL,
albumn_id  VARCHAR(36) NOT NULL,

PRIMARY KEY(id),
FOREIGN KEY(genre_name) REFERENCES genres(name),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);
CREATE TABLE users (
  user_id varchar(36) NOT NULL,
  username varchar(150) NOT NULL,
  email varchar(150) NOT NULL,
  PRIMARY KEY (user_id)
)

 CREATE TABLE playlist(
    id VARCHAR(36) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description VARCHAR(350) DEFAULT "",
    user_id VARCHAR(36) NOT NULL,
    date_creation DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE liked_song(
   id VARCHAR(36) NOT NULL,
song_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(song_id) REFERENCES songs(id));



CREATE TABLE liked_albumns(
    id VARCHAR(36) NOT NULL,
albumn_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(albumn_id) REFERENCES albumns(id)
);


CREATE TABLE liked_authors(
    id VARCHAR(36) NOT NULL,
author_id VARCHAR(255) NOT NULL,
user_id  VARCHAR(36) NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(author_id) REFERENCES authors(id)
)

CREATE TABLE playlist_songs(
       id VARCHAR(36) NOT NULL,
       playlist_id VARCHAR(36) NOT NULL,
       song_id VARCHAR(36) NOT NULL,
       playlist_index INT, NOT NULL,
       user_id VARCHAR(36) DEFAULT NULL,
       PRIMARY KEY(id),
          FOREIGN KEY playlist_id REFERENCES playlist(id),
       FOREIGN KEY  song_id  REFERENCES songs(id),
       FOREIGN KEY user_id REFERENCES users(id) ON DELETE SET NULL
       );
       