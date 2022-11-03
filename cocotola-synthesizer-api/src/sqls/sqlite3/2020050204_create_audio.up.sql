create table `audio` (
 `id` integer primary key autoincrement
,`lang5` varchar(5)
,`text` varchar(100) not null
,`audio_content` text not null
,unique(`lang5`, `text`)
);
