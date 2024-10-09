create table if not exists teacher (
    id uuid primary key,
    first_name varchar not null,
    last_name varchar not null,
    middle_name varchar not null,
    report_email varchar not null unique,
    username varchar not null unique,

    check(length(username) >= 5 and length(username) <= 20)
);
