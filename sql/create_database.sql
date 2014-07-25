DROP TABLE IF EXISTS `ImgurLink`;
CREATE TABLE `ImgurLink` (
	`imgur_link_id` int NOT NULL AUTO_INCREMENT,

	/*
	 * The 7 character imgur id obtained via the imgur API.
	 */
	`imgur_id` char(7) NOT NULL,

	/*
	 * The alias that maps to the imgur id.
	 */
	`alias` varchar(12) NOT NULL,

	/*
	 * When this link was created.
	 */
	`create_date` datetime NOT NULL,

	PRIMARY KEY (`imgur_link_id`),

	/*
	 * Constraint that checks if the alias is in use.
	 */
	CONSTRAINT `alias_exists_uc` UNIQUE (`alias`)
) DEFAULT CHARSET=utf8;

DROP FUNCTION IF EXISTS `get_link_from_alias`;
DELIMITER $$
/*
 * This function returns a fully formed imgur link to the image
 * associated with the given alias.
 * If there is no imgur id associated with the alias then
 * the function returns error 45000
 */
CREATE FUNCTION `get_link_from_alias`(pAlias varchar(12))
RETURNS varchar(40)
BEGIN
	DECLARE link varchar(40);
	SELECT CONCAT('http://i.imgur.com/', imgur_id, '.png')
	INTO link
	FROM ImgurLink
	WHERE alias=pAlias;

	IF (link IS NULL) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT =
		'alias does not map to an imgur id';
	ELSE
		RETURN link;
	END IF;
END $$
DELIMITER ;

/*Add test data*/
INSERT INTO ImgurLink (imgur_id, alias, create_date)
VALUES ('JQJ8kEk', 'timebeing', now());
INSERT INTO ImgurLink (imgur_id, alias, create_date)
VALUES ('SkbPtsx', 'meta', now());
INSERT INTO ImgurLink (imgur_id, alias, create_date)
VALUES ('jQPouaZ', 'iliketurtles', now());

/*Test the function using test data*/
SELECT get_link_from_alias('timebeing') as 'link';
SELECT get_link_from_alias('doesnotexist') as 'link';

