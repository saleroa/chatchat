CREATE TABLE message (
                         mid int comment'消息的id',
                         fromid int unsigned NOT NULL comment'发送者的id',
                         toid  int unsigned NOT NULL comment'接收者的id',
                         content TEXT NOT NULL comment'消息的内容',
                         time TIMESTAMP NOT NULL,
                         sendtype int unsigned NOT NULL,
                         messagetype int unsigned NOT NULL
);
