create table if not exists `users` (
    `id` varchar(36) not null primary key,
    `email` varchar(255) not null unique,
    `username` varchar(255) not null,
    `encrypted_password` varchar(255) not null,
    `created_at` timestamp null default null,
    `updated_at` timestamp null default null,
    `deleted_at` timestamp null default null
);
