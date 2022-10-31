CREATE database testdbrun;

CREATE USER IF NOT EXISTS 'common_user'@'localhost' IDENTIFIED BY "common_pass"; --Check is user exists and create if not

--Granting to all db is a HUGE MISTAKE, [TODO] Needs fix:w
GRANT ALL PRIVILEGES ON *.* TO 'common_user'@localhost IDENTIFIED BY 'common_pass';

-- USE testdbrun;
--
-- CREATE TABLE admins(
--   adm_id varchar(64),
--   adm_hash_pass varchar(64), --Needs better way of storing password with hash and salting
--   PRIMARY KEY (adm_id)
-- );
--
-- CREATE TABLE doc(
--   doc_id varchar(64), --64 char doc id remain as default doc_id, custom configs can change
--   def_permbit binary(8), --8 bit permbits is default [Only 5 is used], custom configs can change
--   adm_id varchar(64),
--   PRIMARY KEY (doc_id),
--   FOREIGN KEY (adm_id) REFERENCES admins(adm_id)
-- );
--
-- CREATE TABLE user_perms(
--   doc_id varchar(64), 
--   user_id varchar(64),
--   PRIMARY KEY (user_id),
--   FOREIGN KEY(doc_id) REFERENCES doc(doc_id)
--   nd_permbit binary(8), --Same structure as def_permbit
-- );
--
-- CREATE TABLE auth_code(
--   adm_id varchar(64),
--   FOREIGN KEY (adm_id) REFERENCES admins(adm_id),
--   auth_code varchar(64)
-- );
