CREATE
EXTENSION vector;

create table documents
(
    id        bigserial primary key,
    content   text,
    embedding vector(1536)
);


create index on documents
    using ivfflat (embedding vector_cosine_ops)
    with
    (lists = 100);