Compose
=======
**Compose** is a simple blogging platform that I initially built over the course of a couple weekends to write and serve content on my own [blog](https://mborgerson.com).

I decided to use the [Go language](https://golang.org/) over other more popular languages largely because it's a language that I had been wanting to learn. As someone who loves C and Python, I think Go is pretty awesome. I highly recommend spending some time with it or using it for your next project. 

I had fun building and using Compose, and I think others might too. It is still in its infancy, but it is usable. If you are interested in contributing, feel free to send me a pull request.

Installing
----------

### Install Required Dependencies

Install Go and Git

On Ubuntu

    $ sudo apt-get install golang git

On Mac OS X

    $ brew install go git

#### Install MongoDB

On Ubuntu

    $ sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
    $ echo "deb http://repo.mongodb.org/apt/ubuntu "$(lsb_release -sc)"/mongodb-org/3.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.0.list
    $ sudo apt-get update
    $ sudo apt-get install -y mongodb-org

On Mac OS X

    $ brew install mongodb

### Setup Go Workspace

    For those unfamiliar with Go, a workspace is needed when developing or
    installing Go packages. To create a workspace, simply assign the `GOPATH`
    environment variable.

    $ mkdir $HOME/go
    $ export GOPATH=$HOME/go
    $ export PATH=$PATH:$GOPATH/bin

### Install Compose

    $ go install github.com/mborgerson/compose

Build the Themes
----------------
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

Install Node.js packages.

    $ npm install

Install client-side dependencies.

    $ bower install

### Build

Finally, just run `gulp` to build the theme.

    $ gulp

### Build other themes

Do the same for theme_site.

Running
-------
First, make sure mongod is running.

    $ service mongod status

For testing purposes, you can start the MongoDB daemon manually with a specific database path with the following command.

    $ mongod --dbpath=./db

Next, as described above, make sure that the `GOPATH` environment variable is set.

    $ export GOPATH=$HOME/go

Next, for convenience, add the workspace **bin** directory to the `PATH` environment variable.

    $ export PATH=$PATH:$GOPATH

Finally, run Compose.

    $ compose

This will create a config file if one does not already exist. Now, run again.

    $ compose

This will start an HTTP server listening at [http://127.0.0.1:8080](http://127.0.0.1:8080).

Navigate to [http://127.0.0.1:8080/setup](http://127.0.0.1:8080/setup) to initialize the database.

Now, you can login and write content at [http://127.0.0.1:8080/login](http://127.0.0.1:8080/login).

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
        exec su -s /bin/sh -c 'exec "$0" "$@"' www-data -- compose

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