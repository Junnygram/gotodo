# gotodo

//configure cli
which aws
aws configure

# Change the permissions of the downloaded godb.pem file to read-only for the owner.

chmod 400 ~/Downloads/godb.pem

# SSH into the EC2 instance with the specified private key (-i) and username.

ssh -i ~/Downloads/godb.pem <Public Dns>

ssh -i ~/Downloads/godb.pem ubuntu@ec2-51-20-10-5.eu-north-1.compute.amazonaws.com

# Update package information and install MySQL client on the EC2 instance.

sudo apt update -y
sudo apt install mysql-client -y

# Connect to the MySQL database server using the provided host (-h), username (-u), and password (-p).

mysql -h mysqlgo-db.c7gi48seic7n.eu-north-1.rds.amazonaws.com -u admin -p

show databases;
use mysqlgodb
CREATE TABLE Task (
-> ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
-> Completed BOOLEAN NOT NULL,
-> Body TEXT NOT NULL
-> );
SHOW TABLES;
DESCRIBE Task;
