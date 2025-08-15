

DELIMITER //



CREATE FUNCTION are_author_similar(first_author_id VARCHAR(36), second_author_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH author_cte AS (
    SELECT   DISTINCT authors.id AS authors_id, collaboration_tags.tag_name AS collaboration_tag, tags_to_albumns.tag_name AS tag  FROM authors
INNER JOIN albumns  ON authors.id = albumns.author_id
INNER JOIN authors_to_albumns ON authors.id = authors_to_albumns.author_id
INNER JOIN albumns AS collaboration_albumns ON   authors_to_albumns.albumn_id = collaboration_albumns.id
INNER JOIN tags_to_albumns  ON  albumns.id= tags_to_albumns.albumn_id
INNER JOIN tags_to_albumns AS collaboration_tags ON collaboration_albumns.id = collaboration_tags.albumn_id)
SELECT COUNT(tags.name) INTO count_result FROM tags 
INNER JOIN author_cte AS first_author ON tags.name = first_author.tag OR tags.name = first_author.collaboration_tag
  INNER JOIN author_cte AS second_author ON tags.name = second_author.tag OR tags.name = second_author.collaboration_tag
WHERE first_author.authors_id = first_author_id  AND  second_author.authors_id = second_author_id
GROUP BY tags.name;
RETURN count_result;
END//

CREATE FUNCTION are__authors_similar_genre(first_author_id VARCHAR(36), second_author_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH author_cte AS (
    SELECT   DISTINCT authors.id AS authors_id, collaboration_genres.genre_name AS collaboration_genre, genres.name AS genre  FROM authors
INNER JOIN albumns  ON author.id = albumn.authors_id
INNER JOIN authors_to_albumns ON authors.id = authors_to_albumns.author_id
INNER JOIN albumns AS collaboration_albumns ON   authors_to_albumns.albumn_id = collaboration_albumns.id
INNER JOIN genres_to_albumns AS genres ON  albumns.id= genres.albumn_id
INNER JOIN genres_to_albumns AS collaboration_genres ON collaboration_albumns.id = collaboration_genres.albumn_id)
SELECT COUNT(genres.name) INTO count_result FROM genres
INNER JOIN author_cte AS first_author ON genres.name = first_author.genre OR genres.name = first_author.collaboration_genre
  INNER JOIN author_cte AS second_author ON genres.name = second_author.genre OR genres.name = second_author.collaboration_genre
WHERE first_author.authors_id = first_author_id  AND  second_author.authors_id = second_author_id
GROUP BY genres.name;
RETURN count_result;
END//


CREATE FUNCTION are__authors_similar_ws_genres(first_author_id VARCHAR(36), second_author_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH author_cte AS (
    SELECT   DISTINCT authors.id AS authors_id, collaboration_genres.genre_name AS collaboration_genre, genres.name AS genre,
    first_genres.second_genre AS first_similar, second_genres.first_genre AS second_similar, first_genres_collaboration.second_genre  AS first_similar_collaboration,
 second_genres_collaboration.first_genre AS second_similar_collaboration    FROM authors
INNER JOIN albumns  ON author.id = albumns.authors_id
INNER JOIN authors_to_albumns ON authors.id = authors_to_albumns.author_id
INNER JOIN albumns AS collaboration_albumns ON   authors_to_albumns.albumn_id = collaboration_albumns.id
INNER JOIN genres_to_albumns AS genres ON  albumns.id= genres.albumn_id
INNER JOIN genres_to_albumns AS collaboration_genres ON collaboration_albumns.id = collaboration_genres.albumn_id
INNER JOIN similar_genres as first_genres ON genres.genre_name = first_genres.first_genre
INNER JOIN similar_genres as second_genres ON genres.genre_name = second_genres.second_genre
INNER JOIN similar_genres as first_genres_collaboration ON collaboration_genres.genre_name = first_genres_collaboration.first_genre 
INNER JOIN similar_genres as second_genres_collaboration ON collaboration_genres.genre_name = second_genres_collaboration.second_genre
)
SELECT COUNT(genres.name) INTO count_result FROM genres
INNER JOIN author_cte AS first_author ON genres.name = first_author.genre OR genres.name = first_author.collaboration_genre OR genres.name=author_cte.first_similar
OR genres.name = author_cte.second_similar OR genres.name = author_cte.first_similar_collaboration  OR genres.name = author_cte.second_similar_collaboration 
  INNER JOIN author_cte AS second_author ON genres.name = second_author.tag OR genres.name = second_author.collaboration_tag OR genres.name=author_cte.first_similar
 OR genres.name = author_cte.second_similar OR genres.name = author_cte.first_similar_collaboration  OR genres.name = author_cte.second_similar_collaboration 
WHERE first_author.authors_id = first_author_id  AND  second_author.authors_id = second_author_id
GROUP BY genres.name;
RETURN count_result;
END//


CREATE FUNCTION are__authors_similar_ws(first_author_id VARCHAR(36), second_author_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH author_cte AS (
    SELECT   DISTINCT authors.id AS authors_id, collaboration_genres.genre_name AS collaboration_genre, genres.name AS genre,
    first_genres.second_genre AS first_similar, second_genres.first_genre AS second_similar, first_genres_collaboration.second_genre  AS first_similar_collaboration,
 second_genres_collaboration.first_genre AS second_similar_collaboration    FROM authors
INNER JOIN albumns  ON author.id = albumns.authors_id
INNER JOIN authors_to_albumns ON authors.id = authors_to_albumns.author_id
INNER JOIN albumns AS collaboration_albumns ON   authors_to_albumns.albumn_id = collaboration_albumns.id
INNER JOIN tags_to_albumns AS tags ON  albumns.id= tags.albumn_id
INNER JOIN tags_to_albumns AS collaboration_tags ON collaboration_albumns.id = collaboration_genres.albumn_id
INNER JOIN similar_tags as first_tags ON tags.tag_name = first_tags.first_tag
INNER JOIN similar_tags as second_tags ON tags.genre_name = second_tags.second_genre
INNER JOIN similar_tags as first_tags_collaboration ON collaboration_tags.tag_name = first_tags_collaboration.first_tag  
INNER JOIN similar_genres as second_tags_collaboration ON collaboration_tags.tag_name = second_tags_collaboration.second_tag
)
SELECT COUNT(tags.name) INTO count_result FROM tags
INNER JOIN author_cte AS first_author ON tags.name = first_author.tag OR tags.name = first_author.collaboration_tag OR tags.name=author_cte.first_similar
OR tags.name = author_cte.second_similar OR tags.name = author_cte.first_similar_collaboration  OR tags.name = author_cte.second_similar_collaboration 
  INNER JOIN author_cte AS second_author ON tags.name = second_author.tag OR tags.name = second_author.collaboration_tag OR tags.name=author_cte.first_similar
 OR tags.name = author_cte.second_similar OR tags.name = author_cte.first_similar_collaboration  OR tags.name = author_cte.second_similar_collaboration 
WHERE first_author.authors_id = first_author_id  AND  second_author.authors_id = second_author_id
GROUP BY tags.name;
RETURN count_result;
END//


CREATE FUNCTION are_songs_albumns_tags_similar(first_song_id VARCHAR(36), second_song_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH songs_tags AS(
SELECT DISTINCT songs.id AS songs_id,  tags_to_albumns.tag_name  AS albumn_tag FROM songs
INNER JOIN albumns ON songs.albumn_id = songs.id 
INNER JOIN tags_to_albumns ON albumns.id = tags_to_albumns.albumn_id 
)
SELECT COUNT(tags_to_songs.tag_name) INTO count_result FROM tags_to_songs
INNER JOIN song_tags AS first_song_tags ON tags_to_songs.tag_name =  first_song_tags.albumn_tag
INNER JOIN  songs_tags AS second_song_tags ON tags_to_songs.tag_name = second_song_tags.albumn_tag
WHERE first_song_tags.songs_id = first_song_id AND second_song_tags = second_song_id
GROUP BY tags_to_songs.name;
RETURN count_result;

END//



CREATE FUNCTION are__songs_albumns_genres_similar(first_song_id VARCHAR(36), second_song_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH songs_genres AS(
SELECT DISTINCT songs.id AS songs_id,  genres_to_albumns.genre_name  AS albumn_genre FROM songs
INNER JOIN albumns ON songs.albumn_id = songs.id 
INNER JOIN genres_to_albumns ON albumns.id = genres_to_albumns.albumn_id 
)
SELECT COUNT(genres_to_songs.tag_name) INTO count_result FROM genres_to_songs
INNER JOIN songs_genres AS first_song_genres ON genres_to_songs.genre_name =  first_song_genres.albumn_genre
INNER JOIN  songs_genres AS second_song_genres ON genres_to_songs.genre_name = second_song_tags.albumn_genre
WHERE first_song_genres.songs_id = first_song_id AND second_song_genres.song_id = second_song_id
GROUP BY genres_to_songs.name;
RETURN count_result;

END//

CREATE FUNCTION are_songs_albumns_tags_similar_ws(first_song_id VARCHAR(36), second_song_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH songs_tags AS(
SELECT DISTINCT songs.id AS songs_id,  tags_to_albumns.tag_name  AS albumn_tag,  first_similar.second_tag AS second_similar_tag, second_similar.first_tag AS first_similar_tag FROM songs
INNER JOIN albumns ON songs.albumn_id = songs.id 
INNER JOIN tags_to_albumns ON albumns.id = tags_to_albumns.albumn_id
INNER JOIN similar_tags AS first_similar ON tags_to_albumns.tag_name = first_similar.first_tag
INNER JOIN similar_tags AS second_similar ON tags_to_albumns.tag_name = second_similar.secodn_tag 

)
SELECT COUNT(tags_to_songs.tag_name) INTO count_result FROM tags_to_songs
INNER JOIN song_tags AS first_song_tags ON tags_to_songs.tag_name =  first_song_tags.albumn_tag OR tags_to_songs.tag_name = first_song_tags.first_similar
 OR tags_to_songs.tag_name= first_song_tags.second_similar
INNER JOIN  songs_tags AS second_song_tags ON tags_to_songs.tag_name = second_song_tags.albumn_tag OR  tags_to_songs.tag_name = second_song_tags.first_similar
 OR tags_to_songs.tag_name= second_song_tags.second_similar
WHERE first_song_tags.songs_id = first_song_id AND second_song_tags = second_song_id
GROUP BY tags_to_songs.name;
RETURN count_result;

END//

CREATE FUNCTION are_songs_albumns_genres_similar_ws(first_song_id VARCHAR(36), second_song_id VARCHAR(36) ) RETURNS INT
READS SQL DATA 
DETERMINISTIC
BEGIN
DECLARE count_result INT;
WITH songs_genres AS(
SELECT DISTINCT songs.id AS songs_id,  albumn_genres.tag_name  AS albumn_genre,  first_similar.second_genre AS second_similar_genre, second_similar.first_genre AS first_similar_genre FROM songs
INNER JOIN albumns ON songs.albumn_id = songs.id 
INNER JOIN genres_to_albumns ON albumns.id = genres_to_albumns.albumn_id
INNER JOIN similar_tags AS first_similar ON genres_to_albumns.genre_name = first_similar.first_genre
INNER JOIN similar_tags AS second_similar ON genres_to_albumns.genre_name = second_similar.second_genre 

)
SELECT COUNT(genre_to_songs.genre_name) INTO count_result FROM genres_to_songs
INNER JOIN song_genres AS first_song_genres ON genres_to_songs.genre_name =  first_song_genres.albumn_genre OR genres_to_songs.genre_name = first_song_genres.first_similar
 OR genres_to_songs.genre_name= first_song_genres.second_similar
INNER JOIN  songs_genres AS second_song_genres ON genres_to_songs.genre_name = second_song_genres.albumn_genre OR  tags_to_songs.genre_name = second_song_genres.first_similar
 OR genres_to_songs.genre_name= second_song_genres.second_similar
WHERE first_song_genres.songs_id = first_song_id AND second_song_genres = second_song_id
GROUP BY genres_to_songs.name;
RETURN count_result;

END//

DELIMITER ;