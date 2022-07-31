# Comparison of Google Cloud SQL vs PersistentVolumeClaims with own DBMS images

Advantages of Google Cloud SQL

- > Google handles replication, patch management and database management to ensure availability and performance.

  This includes automatic backups and maintenance.

- Automatic scaling
- Easy setup using google cloud console
- Cheaper to maintain, as maintenance is taken care of by CloudSQL

Advantages of own DBMS images with PersistentVolumeClaims

- Possibility of using any database and any version of that database.
- Ability to spin up a database for each service in a microservice architecture. This means that each service can potentially use a different type of database, one that best fits its needs.
- Everything you need (apps + databases) run on one platform (Kubernetes), allowing easy porting to different cloud providers (as a result of being decoupled from cloud provider).
- Can be cheaper to run if setup correctly.


## Bibliography

[Google Cloud SQL](https://cloud.google.com/sql)  
[PostgreSQL in the Cloud: DBaaS vs Kubernetes - Michal Nosek | Percona Live 2022](https://www.youtube.com/watch?v=CRCkh8mbrpE)


## Exercise 3.07

For this exercise I have chosen to use Google Cloud SQL because I have have not used it before, whereas Postgres with PersistentVolumeClaims has been covered in the course. As such, I am curious to see which approach is more intuitive.

*Edit*: as it turned out, setting up a connection to Google Cloud SQL proved more challenging, as it requires setting up service accounts and Cloud SQL Auth proxy. This means that this approach requires more knowledge of the underlying cloud platform, rather than just of Kubernetes itself.
