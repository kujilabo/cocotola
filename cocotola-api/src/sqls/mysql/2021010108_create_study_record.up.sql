create table `study_record` (
 `id` varchar(26) not null
,`organization_id` int not null
,`app_user_id` int not null
,`workbook_id` int not null
,`problem_type_id` int not null
,`study_type_id` int not null
,`problem_id` int not null
,`record_date` datetime not null
,`mastered` tinyint
,primary key (`id`)
,unique(`app_user_id`, `workbook_id`, `problem_type_id`, `study_type_id`, `problem_id`, `record_date`)
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
,foreign key(`problem_type_id`) references `problem_type`(`id`) on delete cascade
,foreign key(`study_type_id`) references `study_type`(`id`) on delete cascade
,foreign key(`workbook_id`) references `workbook`(`id`) on delete cascade
);
