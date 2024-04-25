DROP TABLE `user`;
DROP TABLE message;
DROP TABLE message_user;
DROP TABLE user_contact;
DROP TABLE `group`;
DROP TABLE `group_user`;

CREATE TABLE `user` (
                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增',
                        `phone` varchar(32) NOT NULL DEFAULT '' COMMENT '用户手机号',
                        `user_name` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名称',
                        `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '用户昵称',
                        `wx_openid` varchar(128) NOT NULL DEFAULT '' COMMENT '微信openid',
                        `wx_unionid` varchar(128) NOT NULL DEFAULT '' COMMENT '微信unionid',
                        `wx_session_key` varchar(128) NOT NULL DEFAULT '' COMMENT '微信会话密钥session_key',
                        `gender` tinyint(2) NOT NULL DEFAULT '-1' COMMENT '性别',
                        `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
                        `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
                        `last_login_ip` varchar(64) DEFAULT '' COMMENT '最后登录IP',
                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                        PRIMARY KEY (`id`),
                        KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息';

CREATE TABLE `message` (
                           `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增',
                           `sender` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息发送人',
                           `type` varchar(32) NOT NULL DEFAULT '' COMMENT '消息类型',
                           `content` text NOT NULL COMMENT '消息内容',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           KEY `idx_sender` (`sender`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息';

CREATE TABLE `message_user` (
                                `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增',
                                `message_id` bigint(20) unsigned NOT NULL COMMENT '消息ID',
                                `receiver` bigint(20) unsigned NOT NULL COMMENT '消息接收人',
                                `is_read` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '是否已读：0-未读；1-已读；',
                                `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                PRIMARY KEY (`id`),
                                KEY `idx_message_id` (`message_id`),
                                KEY `idx_sender` (`receiver`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息-用户';

CREATE TABLE `user_contact` (
                                `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增',
                                `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户1-申请人',
                                `friend_user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户2-被申请添加好友的人',
                                `status` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '1-申请中；2-好友；3-拒绝；',
                                `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                PRIMARY KEY (`id`),
                                KEY `idx_user` (`user_id`,`friend_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户联系人';

CREATE TABLE `group` (
                         `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增',
                         `title` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '群聊名称',
                         `user_count` INT(11) NOT NULL DEFAULT '0' COMMENT '群聊用户总数',
                         `created_user_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '群创建的用户ID',
                         `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
                         PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='聊天群';

CREATE TABLE `group_user` (
                              `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增',
                              `user_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                              `group_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '聊天ID',
                              `custom_group_title` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '自定义群 标题',
                              `custom_nickname` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '自定义个人在群的昵称',
                              `is_message_remind` TINYINT(2) NOT NULL DEFAULT 0 COMMENT '是否消息提醒：1-提醒；2-免打扰；',
                              `is_top` TINYINT(2) NOT NULL DEFAULT 0 COMMENT '是否消息提醒：1-提醒；2-免打扰；',
                              `is_top_time` DATETIME DEFAULT NULL COMMENT '设置置顶时间',
                              `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
                              PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='聊天群用户';
