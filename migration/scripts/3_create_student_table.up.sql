create table if not exists student (
    id uuid primary key,
    first_name varchar not null,
    last_name varchar not null,
    middle_name varchar not null,
    educational_email varchar not null unique,
    group_id uuid not null,
    username varchar not null unique,

    CONSTRAINT fk_student_group_id FOREIGN KEY (group_id) REFERENCES "group" (id)
);
