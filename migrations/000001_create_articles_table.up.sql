CREATE TABLE articles
(
    id SERIAL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    created_on TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    updated_on TIMESTAMP(0) WITHOUT TIME ZONE,
    deleted_on  TIMESTAMP(0) WITHOUT TIME ZONE,

    PRIMARY KEY (id)
);

-- create unique index for title up and down