create table `user_preference` (
 `id` int auto_increment
,`version` int not null
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp on update current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`app_user_id` int not null
,`preference` json not null
,primary key(`id`)
,unique(`organization_id`, `app_user_id`)
,foreign key(`created_by`) references `app_user`(`id`) on delete cascade
,foreign key(`updated_by`) references `app_user`(`id`) on delete cascade
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
);
