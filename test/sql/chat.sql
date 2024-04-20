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