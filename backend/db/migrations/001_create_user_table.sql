CREATE TABLE IF NOT EXISTS users (
  "id"            uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
  "completename"  varchar(255)    NOT NULL,
  "cpf"           varchar(11),
  "phone"         varchar(50),
  "email"         varchar(50)     NOT NULL,
  "password"      varchar(100)    NOT NULL,
  "country"       smallint        NOT NULL,
  "createat"      timestamp       NOT NULL,
  "lastlogin"     timestamp       NOT NULL,
  "usertype"      smallint        NOT NULL DEFAULT 1,
  UNIQUE(email)
)

-- undo
DROP TABLE IF EXISTS users