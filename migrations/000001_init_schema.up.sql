CREATE TABLE hello_table (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message text
);