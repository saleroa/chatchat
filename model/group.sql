CREATE TABLE `groups` (
                        gid int unsigned auto_increment primary key ,
                        group_name VARCHAR(255) NOT NULL unique,
                        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members (
                               group_id INT UNSIGNED NOT NULL,
                               user_id INT UNSIGNED NOT NULL,
                               identity int unsigned  not null  default 0,
                               PRIMARY KEY (group_id, user_id),
                               FOREIGN KEY (group_id) REFERENCES `groups`(gid) ON DELETE CASCADE,
                               FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

