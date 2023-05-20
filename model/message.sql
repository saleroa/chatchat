CREATE TABLE message (
                         mid int
                         fromid int unsigned NOT NULL ,
                         toid  int unsigned NOT NULL,
                         content TEXT NOT NULL,
                         time TIMESTAMP NOT NULL,
                         sendtype int unsigned NOT NULL,
                         messagetype int unsigned NOT NULL
);
