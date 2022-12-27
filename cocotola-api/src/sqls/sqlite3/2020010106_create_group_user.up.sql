create table `group_user` (
 `created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`app_user_group_id` int not null
,`app_user_id` int not null
,primary key(`organization_id`, `app_user_group_id`, `app_user_id`)
,unique(`app_user_group_id`, `app_user_id`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
);
