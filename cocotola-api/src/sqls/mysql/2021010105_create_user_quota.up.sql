create table `user_quota` (
 `id` varchar(26) not null
,`organization_id` int not null
,`app_user_id` int not null
,`date` datetime not null
,`name` varchar(32) not null
,`unit` varchar(16) not null
,`count` int not null
,primary key(`id`)
,unique(`organization_id`, `app_user_id`, `date`, `name`)
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
);
