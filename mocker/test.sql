CREATE TABLE IF NOT EXISTS `tbl_user` (
    `id` bigint(20) not null AUTO_INCREMENT,
    `name` varchar(64) not null default '' comment '姓名',
    `sex` tinyint(4) not null default 0 comment '性别: 1-男 2-女',
    `age` int(8) not null default 0 comment '年龄',
    primary key (`id`)
);

insert into `tbl_user`(`name`, `sex`, `age`) values('lp', 1, 18);