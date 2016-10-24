go-blog
===
go-blog is a blog program. It is easy to deploy.

## how to write a blog post
To write a program, you need to write a md file under the post directory. FYI, example.md is provided.

Before writing your post, some info should be provided in the md file. the info includs:

|info name | description |
|---|---|
|@title| the title of the blog post|
|@tags| the tags of the blog post|
|@summary| the summary of the blog post|
|@author| the author of the blog post|
|@date| date of the blog post|

After providing the infos above, the body starts with a line of x's. Then in the next line, the body of the post starts.

## how to deploy the post
To deploy the post, you just need to put the post under the ``` /post ``` directory, and restart the program by running ```./myblog``` command, then the program will re-scan the post directory and sync the data into the database.

By default, if the post is kept unchanged, the program will check the file's md5 and find it unchanged, it just ignored it.
