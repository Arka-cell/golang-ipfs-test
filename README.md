### Dependencies:
Setup your database with Docker:
<ol>
    <li>Populate your .env file with your environment variables.
    <li>Run `docker-compose up`. This will create a database for you with MySQL and will automatically create a table called messages to serialize into JSON as well as seed it with dummy data.
</ol>
Install IPFS:
<ol>
    <li>Install IPFS (In my case I have used IPFS Desktop): https://docs.ipfs.tech/install/ipfs-desktop/
    <li>Go to settings and get IPNS Publishing Key, add it to the environment variable `IPNS_PUBLISH_KEY` in your .env file.
</ol>


<p>You can use your own database with its user and migrate a table into it called messages. To do that, connect with <a href="https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-install-linux-quick.html">mysql</a> client in your machine and run `source init.sql`. This will create a table and seed data into your table. E.g;

<li>Connect to the database with MySQL client:
    
    $ pwd

    /ummatest/

    $ mysql -h <YOUR_DB_HOST> -P 3306 -u <YOUR_USERNAME> -p <YOUR_DB_NAME> --password=<YOUR_USER_PASSWORD>

    mysql> source init.sql

### Run Go App
Currently, the code is running with go1.18.1
Install your dependencies with:

    go install .

Run your Go app with the following command:

    go run . --user=$DATABASE_USER --password=$DATABASE_PASSWORD --database=$DATABSE_NAME

Once finished, check json_result.json which will show all data from the table messages as json, but it will also show their IPFS and resolved IPNS. 
