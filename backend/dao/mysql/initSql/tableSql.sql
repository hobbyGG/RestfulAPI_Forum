create table `user` (
    `id` bigint(20) not null auto_increment,
    `user_id` bigint(20) not null,
    `username` varchar(64) not null,
    `password` varchar(64) not null,
    `email` varchar(64) not null,
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update current_timestamp,
    primary key (`id`),
    UNIQUE KEY `idx_email` (`email`),
    unique key `idx_username` (`username`),
    unique key `idx_user_id` (`user_id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

-- 插入测试用户数据
INSERT INTO `user` (`user_id`, `username`, `password`, `email`)
VALUES 
(-1, 'test1', '123', 'test1@example.com'),
(-2, 'test2', '123', 'test2@example.com'),
(-3, 'test3', '123', 'test3@example.com');

CREATE Table `post` (
    `postID` BIGINT(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `authorUID` BIGINT(20) NOT NULL,
    `score` BIGINT(20) NOT NULL,
    `status` TINYINT(4) NOT NULL,
    `commID` TINYINT(4) NOT NULL,
    `title` VARCHAR(128) NOT NULL COLLATE utf8mb4_general_ci ,
    `content` VARCHAR(8192) NOT NULL COLLATE utf8mb4_general_ci,
    `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 建议使用下面的约束
    UNIQUE KEY `idx_postID` (`postID`),
    key `id_authorUID` (`authorUID`),
    KEY `idx_commID` (`commID`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;