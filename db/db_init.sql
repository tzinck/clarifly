DO
$body$
BEGIN
  IF NOT EXISTS (
      SELECT                       -- SELECT list can stay empty for this
      FROM   pg_catalog.pg_user
      WHERE  usename = 'navi') THEN

    CREATE ROLE navi LOGIN PASSWORD 'heylisten';
  END IF;
END
$body$;

ALTER USER navi WITH SUPERUSER;
GRANT ALL PRIVILEGES ON DATABASE clarifly TO navi;

\connect clarifly


CREATE TABLE IF NOT EXISTS rooms (
  room_code   VARCHAR(4) PRIMARY KEY,
  creator_id  INT,
  secret      VARCHAR(32),
  start_time  TIMESTAMP,
  end_time    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
  u_id        SERIAL PRIMARY KEY,
  ip_addr     VARCHAR(15),
  room_code   VARCHAR(4) REFERENCES rooms(room_code),
  join_time   TIMESTAMP,
  leave_time  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS questions (
  q_id        SERIAL PRIMARY KEY,
  u_id        INT REFERENCES users(u_id),
  room_code   VARCHAR(4) REFERENCES rooms(room_code),
  text        VARCHAR(140),
  votes       INT,
  reports     INT,
  hide        BOOLEAN,
  ask_time    TIMESTAMP
);
