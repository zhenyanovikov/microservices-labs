BEGIN;

CREATE TABLE persons
(
    id   SERIAL       NOT NULL,
    name VARCHAR(100) NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO persons (name) VALUES ('John');
INSERT INTO persons (name) VALUES ('Jane');
INSERT INTO persons (name) VALUES ('Joe');

COMMIT;
