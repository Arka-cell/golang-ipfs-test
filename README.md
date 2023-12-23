### Dependencies:
Setup your database:
<ol>
<li>Populate your .env file with your desired environment variables.
<li>Install IPFS (In my case I have used IPFS Desktop): https://docs.ipfs.tech/install/ipfs-desktop/
<li>In order to start MySQL as fast as possible, make sure to install Docker and Docker Compose.

<li>You can create your own database with its user and migrate a table into it called messages. To do that, connect with <a href="https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-install-linux-quick.html">mysql</a> client in your machine and run source `init.sql`. This will create a table and seed data into your table.
<li>In order to seed your data run the following commands:
<ul>
<li>Connect to the database with MySQL client:

    $ mysql -h <YOUR_DB_HOST> -P 3306 -u <YOUR_USERNAME> -p <YOUR_DB_NAME> --password=<YOUR_USER_PASSWORD>

    mysql> source init.sql

</ul>

</ol>
