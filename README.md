## JoJo

JoJo makes is easy to transform your scripts into an api.

I wanted something stupid easy to install that could be 100% repeated with code from configuration management. 
With this, I can create machines that have a single responsibility.  I can also quit logging in from one 
machine into another.

## Note

This is a rather dirty proof of concept and one of my first go projects...it could be better.  If you'd like to
contribute or make the go code better, I would be happy to accept pull requests.

## Options

    --config 'config file'
    --host 'host'
    --port 'port'
    --cert '/path/to/cert'
    --key '/path/to/key'
    --user 'http basic user'
    --password 'http basic password'

## Get started

1. compile jojo (I'll eventually make a little repository, but it will probably be a while)
2. create a config file - this is where you tell jojo about your scripts and what url to serve them from
3. create an init file - jojo just logs to stderr and stdout - capture those in your logs however you like

## Sample Usage and Responses

*NOTE:* These are from running the examples

##### test.sh

    curl -X POST "https://localhost:3000/test" | python -m json.tool
    {
        "arguments": [], 
        "duration": "4.574145702s", 
        "exit-status": 0, 
        "script": "/tmp/test.sh", 
        "stderr": [], 
        "stdout": [
            "libmongodb.x86_64                       2.2.3-4.fc17                 @updates   ", 
            "mongodb.x86_64                          2.2.3-4.fc17                 @updates   ", 
            "mongodb-server.x86_64                   2.2.3-4.fc17                 @updates   ", 
            "libmongo-client.i686                    0.1.6.1-1.fc17               updates    ", 
            "libmongo-client.x86_64                  0.1.6.1-1.fc17               updates    ", 
            "libmongo-client-devel.i686              0.1.6.1-1.fc17               updates    ", 
            "libmongo-client-devel.x86_64            0.1.6.1-1.fc17               updates    ", 
            "libmongo-client-doc.noarch              0.1.6.1-1.fc17               updates    ", 
            "mongo-10gen.x86_64                      2.4.5-mongodb_1              10gen      ", 
            "mongo-10gen-server.x86_64               2.4.5-mongodb_1              10gen      ", 
            "mongo-10gen-unstable.x86_64             2.5.0-mongodb_1              10gen      ", 
            "mongo-10gen-unstable-server.x86_64      2.5.0-mongodb_1              10gen      ", 
            "mongo-java-driver.noarch                2.7.3-1.fc17                 updates    ", 
            "mongo-java-driver-bson.noarch           2.7.3-1.fc17                 updates    ", 
            "mongo-java-driver-bson-javadoc.noarch   2.7.3-1.fc17                 updates    ", 
            "mongo-java-driver-javadoc.noarch        2.7.3-1.fc17                 updates    ", 
            "mongo18-10gen.x86_64                    1.8.5-mongodb_1              10gen      ", 
            "mongo18-10gen-server.x86_64             1.8.5-mongodb_1              10gen      ", 
            "mongo20-10gen.x86_64                    2.0.8-mongodb_1              10gen      ", 
            "mongo20-10gen-server.x86_64             2.0.8-mongodb_1              10gen      ", 
            "mongodb-devel.i686                      2.2.3-4.fc17                 updates    ", 
            "mongodb-devel.x86_64                    2.2.3-4.fc17                 updates    ", 
            "mongoose.x86_64                         3.1-1.fc17                   updates    ", 
            "mongoose-devel.i686                     3.1-1.fc17                   updates    ", 
            "mongoose-devel.x86_64                   3.1-1.fc17                   updates    ", 
            "mongoose-lib.i686                       3.1-1.fc17                   updates    ", 
            "mongoose-lib.x86_64                     3.1-1.fc17                   updates    ", 
            "pdns-backend-mongodb.x86_64             3.1-4.fc17                   updates    ", 
            "php-Monolog-mongo.noarch                1.2.1-1.fc17                 updates    ", 
            "php-pecl-mongo.x86_64                   1.2.12-1.fc17                updates    ", 
            "pymongo.x86_64                          2.1.1-1.fc17                 fedora     ", 
            "pymongo-gridfs.x86_64                   2.1.1-1.fc17                 fedora     ", 
            "python-asyncmongo.noarch                0.1.3-2.fc17                 fedora     ", 
            "python-mongoengine.noarch               0.7.9-4.fc17                 updates    ", 
            "rubygem-openshift-origin-auth-mongo.noarch", 
            "rubygem-openshift-origin-auth-mongo-doc.noarch"
        ]
    }
    

##### test2.sh

    curl -X POST "https://localhost:3000/test2?--foo=bar&--number=2" | python -m json.tool
    {
        "arguments": [
            {
                "--foo": "bar"
            }, 
            {
                "--number": "2"
            }
        ], 
        "duration": "2.580674ms", 
        "exit-status": 0, 
        "script": "/tmp/test2.sh", 
        "stderr": [], 
        "stdout": [
            "test", 
            "arg: --foo", 
            "arg: bar", 
            "arg: --number", 
            "arg: 2", 
            "----------------\"\"-----", 
            "4", 
            "--foo bar --number 2"
        ]
    }



## Warning

This tool is a really sharp knife.  You can do really stupid things with it.  Use at your own risk and under 
close supervision from someone comfortable with doing dumb things in a safe way.

## License
   Copyright 2013 Tim Ray

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
