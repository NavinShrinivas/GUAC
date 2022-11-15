CREATE database testdbrun;

CREATE USER IF NOT EXISTS 'common_user'@'localhost' IDENTIFIED BY "common_pass"; --Check is user exists and create if not

--Granting to all db is a HUGE MISTAKE, [TODO] Needs fix
GRANT ALL PRIVILEGES ON *.* TO 'common_user'@localhost IDENTIFIED BY 'common_pass';

USE testdbrun;

CREATE TABLE admins(
  adm_id varchar(64),
  adm_hash_pass varchar(64),
  PRIMARY KEY (adm_id)
);

CREATE TABLE doc(
  doc_id varchar(64), 
  def_permbit int(1), 
  adm_id varchar(64),
  PRIMARY KEY (doc_id),
  FOREIGN KEY (adm_id) REFERENCES admins(adm_id)
);

CREATE TABLE user_perms(
  doc_id varchar(64), 
  user_id varchar(64),
  nd_permbit int(1), 
  FOREIGN KEY(doc_id) REFERENCES doc(doc_id)
);

CREATE TABLE auth_code(
  adm_id varchar(64),
  auth_code varchar(64),
  FOREIGN KEY (adm_id) REFERENCES admins(adm_id)
);


INSERT INTO admins VALUES ("adm1","pass1");
INSERT INTO admins VALUES ("adm2","pass2");
INSERT INTO admins VALUES ("adm3","pass3");
INSERT INTO admins VALUES ("adm4","pass3");
INSERT INTO admins VALUES ("adm5","pass5");

INSERT INTO doc VALUES ("doc1",8,"adm1");
INSERT INTO doc VALUES ("doc2",0,"adm2");
INSERT INTO doc VALUES ("doc3",11,"adm2");
INSERT INTO doc VALUES ("doc4",8,"adm3");
INSERT INTO doc VALUES ("doc5",8,"adm4");

INSERT INTO user_perms VALUES ("doc1","user1",10);
INSERT INTO user_perms VALUES ("doc2","user2",8);
INSERT INTO user_perms VALUES ("doc3","user3",8);
INSERT INTO user_perms VALUES ("doc4","user3",8);
INSERT INTO user_perms VALUES ("doc5","user4",8);

DELIMITER $$

CREATE TRIGGER before_delete_doc BEFORE DELETE ON doc FOR EACH ROW BEGIN DECLARE error_msg varchar(256);DECLARE good_msg varchar(256);DECLARE cnt int;
set
  error_msg = "No user records deleted for this document";
set
  good_msg = "Delete all user record for that document";
SET
  cnt = (
    SELECT
      count(*)
    FROM
      doc
    where
      doc_id = old.doc_id
  );IF cnt >= 1 THEN
DELETE FROM
  user_perms
WHERE
  doc_id = old.doc_id;END IF;END;$$

DELIMITER ;
