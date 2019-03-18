create table article(
  ArticleId   bigint auto_increment primary key,
  SourceName  mediumtext not null,
  Title       mediumtext not null,
  Link        mediumtext not null,
  Description mediumtext,
  PubDate     datetime not null,
  Category    mediumtext,
  PictureUrl  mediumtext
);
create unique index index_article on article(SourceName(100), Link(500));

create table log(
  Severity mediumtext not null,
  UserId mediumtext,
  UserIP mediumtext,
  Timestamp datetime not null,
  Tag mediumtext not null,
  Message mediumtext not null
);
create index index_log_datetime on log(Timestamp);
create index index_log_severity on log(Severity(50));

create table userInfo(
  UserId bigint auto_increment primary key,
  UserName mediumtext,
  UserPhone mediumtext,
  UserPhoto blob
);

create table userDevices(
  UserId bigint not null,
  DeviceId mediumtext not null,
  FirebaseId mediumtext not null,
  FOREIGN KEY (UserId) REFERENCES userInfo(UserId)
);
create unique index index_userDevices on userDevices(UserId, DeviceId(100));

create table userDeviceToken(
  UserId bigint not null,
  DeviceId mediumtext not null,
  Timestamp datetime not null,
  Token mediumtext not null,
  FOREIGN KEY (UserId) REFERENCES userInfo(UserId)
);
create unique index index_userDeviceToken on userDeviceToken(UserId, DeviceId(100), Timestamp);

create table userArticleLikes(
  UserId bigint not null,
  ArticleId   bigint not null,
  Dislike bool not null,
  FOREIGN KEY (UserId) REFERENCES userInfo(UserId),
  FOREIGN KEY (ArticleId) REFERENCES article(ArticleId)
);
create unique index index_userArticleLikes on userArticleLikes(UserId, ArticleId);

create table userArticleComments(
  UserId bigint not null,
  ArticleId   bigint not null,
  Timestamp datetime not null,
  Comment  mediumtext not null,
  Status int not null,
  FOREIGN KEY (UserId) REFERENCES userInfo(UserId),
  FOREIGN KEY (ArticleId) REFERENCES article(ArticleId)
);
create unique index index_userArticleComments on userArticleComments(UserId, ArticleId, Timestamp);
