use dbnews;

create table articles(
  ArticleId   bigint auto_increment primary key,
  SourceName  mediumtext not null,
  Title       mediumtext not null,
  Link        mediumtext not null,
  Description mediumtext,
  PubDate     datetime not null,
  Category    mediumtext,
  PictureUrl  mediumtext
);
create unique index index_articles on articles(SourceName(100), Link(500));

create table logs(
  Severity mediumtext not null,
  UserId mediumtext,
  UserIP mediumtext,
  Timestamp datetime not null,
  Tag mediumtext not null,
  Message mediumtext not null
);
create index index_logs_datetime on logs(Timestamp);
create index index_logs_severity on logs(Severity(50));

create table users(
  UserId bigint auto_increment primary key,
  UserName mediumtext,
  UserPhone mediumtext
);
create unique index index_users_UserName on users(UserName(100));

create table userDevices(
  UserId bigint not null,
  DeviceId mediumtext not null,
  FirebaseId mediumtext not null,
  FOREIGN KEY (UserId) REFERENCES users(UserId) ON DELETE CASCADE
);
create unique index index_userDevices on userDevices(UserId, DeviceId(100));

create table userDeviceTokens(
  UserId bigint not null,
  DeviceId mediumtext not null,
  Timestamp datetime not null,
  Token mediumtext not null,
  FOREIGN KEY (UserId) REFERENCES users(UserId) ON DELETE CASCADE
);
create unique index index_userDeviceTokens on userDeviceTokens(UserId, DeviceId(100), Timestamp);

create table articleLikes(
  UserId bigint not null,
  ArticleId   bigint not null,
  Dislike bool not null,
  Timestamp datetime not null,
  FOREIGN KEY (UserId) REFERENCES users(UserId) ON DELETE CASCADE,
  FOREIGN KEY (ArticleId) REFERENCES articles(ArticleId) ON DELETE CASCADE
);
create unique index index_articleLikes on articleLikes(UserId, ArticleId);
create index index_articleLikes_timestamp on articleLikes(Timestamp);

create table articleComments(
  CommentId  bigint auto_increment primary key,
  UserId bigint not null,
  ArticleId   bigint not null,
  Timestamp datetime not null,
  Comment  mediumtext not null,
  Status int not null,
  FOREIGN KEY (UserId) REFERENCES users(UserId) ON DELETE CASCADE,
  FOREIGN KEY (ArticleId) REFERENCES articles(ArticleId) ON DELETE CASCADE
);
create unique index index_articleComments on articleComments(UserId, ArticleId, Timestamp);
create index index_articleComments_timestamp on articleComments(Timestamp);

create table articleCommentLikes(
  UserId bigint not null,
  CommentId   bigint not null,
  Dislike bool not null,
  Timestamp datetime not null,
  FOREIGN KEY (UserId) REFERENCES users(UserId) ON DELETE CASCADE,
  FOREIGN KEY (CommentId) REFERENCES articleComments(CommentId) ON DELETE CASCADE
);
create unique index index_articleCommentLikes on articleCommentLikes(UserId, CommentId);
create index index_arfticleCommentLikes_timestamp on articleCommentLikes(Timestamp);
