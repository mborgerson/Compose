Compose
=======
**Compose** is a really simple blogging platform that I built over the course of a couple weekends to write and serve content on my [website](https://mborgerson.com).

![](https://github.com/mborgerson/Compose/raw/master/screenshots/admin_edit_post_body.png)
![](https://github.com/mborgerson/Compose/raw/master/screenshots/view_post.png)

For more screenshots, [click here](https://github.com/mborgerson/Compose/blob/master/SCREENSHOTS.md).

I decided to use the [Go Programming Language](https://golang.org/) over other more popular languages largely because it's a language that I had been wanting to learn. As someone who loves C and Python, I think Go is pretty awesome. I highly recommend spending some time with it or trying it out for your next project. [MongoDB](https://www.mongodb.org/) is used to store all the data.

I had fun building and using Compose, and I think others might too. It is still in pre-alpha stages and as such it is subject to some changes, but it is usable. If you're interested in contributing, feel free to send me feedback or, better yet, a pull request.

Current Features
----------------
* Simple
* REST API
* Admin Interface
* File Uploads
* Markdown
* Custom Themes

Installing
----------

### Install Required Dependencies

Install Go

On Ubuntu

    $ sudo apt-get install golang

On Mac OS X

    $ brew install go

#### Install MongoDB

On Ubuntu

    $ sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
    $ echo "deb http://repo.mongodb.org/apt/ubuntu "$(lsb_release -sc)"/mongodb-org/3.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.0.list
    $ sudo apt-get update
    $ sudo apt-get install -y mongodb-org

On Mac OS X

    $ brew install mongodb

### Setup Go Workspace

For those unfamiliar with Go, a workspace is needed when developing or installing Go packages. To create a workspace, simply assign the `GOPATH` environment variable.

    $ mkdir $HOME/go
    $ export GOPATH=$HOME/go
    $ export PATH=$PATH:$GOPATH/bin

### Install Compose

    $ go install github.com/mborgerson/Compose/compose

Building the Themes
-------------------
For running Compose as-is, you do not need to build the themes. If you are developing Compose however, you will probably want to modify the themes/templates.

### Install Required Dependencies

To build the themes, Node.js and NPM are required.

On Ubuntu
    
    $ sudo apt-get install nodejs-legacy npm

On Mac OS X

    $ brew install node npm

Install Gulp and Bower globally.

    $ sudo npm install -g bower gulp

Now, move to the theme_admin directory.

    $ cd theme_admin

Install Node.js packages and client-side dependencies.

    $ npm install
    $ bower install

### Build

Finally, run `gulp` to build the theme.

    $ gulp

### Build other themes

Do the same for theme_site.

Running
-------
Assuming you're serving from an Ubuntu box...

First, make sure mongod is running.

    $ service mongod status

For testing purposes, you can start the MongoDB daemon manually with a specific database path with the following command.

    $ mongod --dbpath=./db

Next, as described above, make sure that the `GOPATH` environment variable is set. For convenience, add the workspace **bin** directory to the `PATH` environment variable.

    $ export GOPATH=$HOME/go
    $ export PATH=$PATH:$GOPATH

Finally, run Compose.

    $ compose

This will create a config file if one does not already exist. Now, run again.

    $ compose

This will start an HTTP server listening at [http://127.0.0.1:8000](http://127.0.0.1:8000).

Navigate to [http://127.0.0.1:8000/setup](http://127.0.0.1:8000/setup) to initialize the database.

Now, you can login and write content at [http://127.0.0.1:8000/login](http://127.0.0.1:8000/login).

Deployment
----------
You are free to run Compose independently. For security and efficiency however, I recommend setting up an NGINX reverse proxy with HTTPS and caching enabled.

Useful Tips
-----------
### Save/Restore a MongoDB Database
You can save and restore MongoDB database with relative ease. This is especially
for backing up your database, or deploying it from your development system.

    mongodump -d compose -o compose_dump/
    mongorestore -d compose compose_dump/compose

### Run on Startup
If you're using a version of Ubuntu with Upstart (e.g. 14.04), you can copy the following script to **/etc/init/compose.conf**. This will automatically start Compose after the MongoDB daemon has been started. Assuming your Go workspace is at **/srv/blog/go_workspace**, your config file and themes are at **/srv/blog**, and the user you want to use is **www-data**.

    description "Compose Server"
    author      "Matt Borgerson"

    start on started mongod
    stop on shutdown

    script

        export GOPATH="/srv/blog/go_workspace"
        echo $$ > /var/run/compose.pid
        chdir /srv/blog/
        exec su -s /bin/sh -c 'exec "$0" "$@"' www-data -- /srv/blog/go_workspace/bin/compose

    end script

    pre-start script
        echo "[`date`] Compose Starting" >> /var/log/compose.log
    end script

    pre-stop script
        rm /var/run/compose.pid
        echo "[`date`] Compose Stopping" >> /var/log/compose.log
    end script

License
-------
Compose is licensed under the terms of the GPLv3 license. See LICENSE.txt for the full license text.

Todo List
---------
Just some features I had in mind:

    ☑ Index
    ☑ View Posts
    ☑ Create/Edit/Delete Posts
    ☑ Markdown Support
    ☑ Basic User Authentication
    ☑ File Uploads
    ☑ Basic Post and File Caching Headers
    ☑ Drafts
    ☑ Custom/Changeable Post Slugs
    ☐ User Images and Biography Pages
    ☐ Non-blog Pages (About, Contact, etc)
    ☐ Post Previews
    ☐ API Error Handling
    ☐ API Documentation
    ☐ Better Error Handling/Reporting
    ☐ Multi User Management and Privileges
    ☐ Limit Number of Failed Logins per IP address
    ☐ Live Theme Switching
    ☐ Logging
    ☐ Statistics
    ☐ Dashboard
    ☐ Post Revision Tracking