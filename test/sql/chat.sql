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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息-用户'
;