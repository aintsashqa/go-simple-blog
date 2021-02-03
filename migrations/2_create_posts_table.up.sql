create table if not exists `posts` (
    `id` varchar(36) not null primary key,
    `title` varchar(255) not null,
    `slug` varchar(255) not null unique,
    `content` text not null,
    `user_id` varchar(36) not null references `users` (`id`) on delete cascade,
    `created_at` timestamp null default null,
    `updated_at` timestamp null default null,
    `published_at` timestamp null default null,
    `deleted_at` timestamp null default null
);
