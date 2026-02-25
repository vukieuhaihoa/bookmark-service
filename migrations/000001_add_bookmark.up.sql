CREATE TABLE bookmarks 
(
  id varchar(36),
  description varchar(256),
  url varchar(2048) NOT NULL,
  user_id varchar(36) NOT NULL,
  code_shorten BIGSERIAL,
  code_shorten_encoded VARCHAR(25),

  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT bookmarks_pk PRIMARY KEY (id),
  CONSTRAINT bookmarks_code_shorten_unique UNIQUE (code_shorten),
  CONSTRAINT bookmarks_code_shorten_encoded_unique UNIQUE (code_shorten_encoded)
);

