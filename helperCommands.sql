-- //Adding a new movie
-- //title: Forrest Gump,
-- //year: 1994,
-- //release date: 1994-07-06,
-- //mpaa_rating: PG-13,
-- //runtime: 142,
-- //created_at: "now()",
-- //updated_at: "now()"
-- //rating (1-5): 4,
-- //description (optional): Forrest Gump is a 1994 American romantic comedy film directed by Robert Zemeckis. The film stars Tom Hanks in the title role, and Robin Wright in the supporting role of Winston Groom. The film is based on the 1986 novel of the same name by Walter Isaacson. The film is the first in the Hanks film series, and the first in the Hanks film series to be nominated for an Academy Award for Best Pictures.
-- INSERT INTO movies (title, year, release_date, mpaa_rating, runtime,rating, description) VALUES ('Forrest Gump', 1994, '1994-07-06', 'PG13', 142, 4, 'Forrest Gump is a 1994 American romantic comedy film directed by Robert Zemeckis. The film stars Tom Hanks in the title role, and Robin Wright in the supporting role of Winston Groom. The film is based on the 1986 novel of the same name by Walter Isaacson. The film is the first in the Hanks film series, and the first in the Hanks film series to be nominated for an Academy Award for Best Pictures.');
-- DELETE FROM movies WHERE id = 6;
-- SELECT * FROM movies;
-- //alter table movies
--     alter column updated_at ON UPDATE CURRENT_TIMESTAMP
-- alter table movies
--     alter column updated_at set default now();

-- create a new table for users: id, username, password, created_at, updated_at
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);
