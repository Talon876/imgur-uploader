/* sed 's/##MYSQL_USER/username/g' setup.sql > setup_real.sql */
/* sed 's/##MYSQL_PASSWORD/password/g' setup.sql > setup_real.sql */

DROP USER '##MYSQL_USER'@'localhost';
CREATE USER '##MYSQL_USER'@'localhost' IDENTIFIED BY '##MYSQL_PASSWORD';
GRANT ALL PRIVILEGES ON stpp . * TO '##MYSQL_USER'@'localhost';
FLUSH PRIVILEGES;

